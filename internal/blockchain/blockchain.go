package blockchain

import (
	"GoBlockchain/internal/wallet"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks_%s"
	genesisData = "First Transaction from Genesis"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain(address string, nodeId string) *BlockChain {
	path := fmt.Sprintf(dbPath, nodeId)
	// Create a new blockchain
	var lastHash []byte

	if DBexists(path) {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}
	opts := badger.DefaultOptions(path)
	opts.Dir = path
	opts.ValueDir = path
	opts.EventLogging = false
	opts.Logger = nil
	db, err := openDB(path, opts)
	Handle(err)

	// Our db will hold 2 types of kv pairs - an "lh" / hash pair to store our last hash,
	// And hash / block pairs to store and retrieve each block
	err = db.Update(func(txn *badger.Txn) error {
		// check for lh (last hash) in the database. If it doesn't exist, create the genesis block
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis proof complete")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func ContinueBlockChain(nodeID string) *BlockChain {
	path := fmt.Sprintf(dbPath, nodeID)
	if !DBexists(path) {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}
	var lastHash []byte

	opts := badger.DefaultOptions("")
	opts.Dir = path
	opts.ValueDir = path
	opts.EventLogging = false
	opts.Logger = nil
	db, err := openDB(path, opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(lastHash)

		return err
	})
	Handle(err)
	chain := BlockChain{lastHash, db}
	return &chain
}

func DBexists(path string) bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}

	return true
}

func (chain *BlockChain) MineBlock(transactions []*Transaction) *Block {
	// Create a new block with the given data and the previous hash
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		if !chain.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	// View is a read only badger db transaction
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		// Get the prev hash from the db
		lastHash, err = item.ValueCopy(lastHash)
		Handle(err)

		lastBlockData, _ := item.ValueCopy(nil)
		lastBlock := Deserialize(lastBlockData)
		lastHeight = lastBlock.Height

		return err
	})

	Handle(err)

	newBlock := CreateBlock(transactions, lastHash, lastHeight+1)

	// Update the db with a new kv entry for the new hash + block
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
	return newBlock
}

func (chain *BlockChain) GetBestHeight() int {
	var lastBlock Block

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, _ := item.ValueCopy(nil)

		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData, _ := item.ValueCopy(nil)

		lastBlock = *Deserialize(lastBlockData)

		return nil
	})
	Handle(err)

	return lastBlock.Height
}

func (chain *BlockChain) GetBlockHashes() [][]byte {
	var blocks [][]byte

	iter := chain.Iterator()

	for {
		block := iter.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return blocks
}

func (chain *BlockChain) AddBlock(block *Block) {
	err := chain.Database.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(block.Hash); err == nil {
			return nil
		}

		err := txn.Set(block.Hash, block.Serialize())
		Handle(err)

		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, _ := item.ValueCopy(nil)

		item, err = txn.Get(lastHash)
		Handle(err)
		lastBlockData, _ := item.ValueCopy(nil)
		lastBlock := Deserialize(lastBlockData)

		if block.Height > lastBlock.Height {
			err = txn.Set([]byte("lh"), block.Hash)
			Handle(err)
			chain.LastHash = block.Hash
		}

		return nil
	})
	Handle(err)
}

func (chain *BlockChain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := chain.Database.View(func(txn *badger.Txn) error {
		if item, err := txn.Get(blockHash); err != nil {
			return errors.New("Block is not found")
		} else {
			encodedBlock, _ := item.ValueCopy(nil)
			block = *Deserialize(encodedBlock)
		}
		return nil
	})

	if err != nil {
		return block, err
	}

	return block, nil

}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Find UTXO (Unspent transaction output)
func (chain *BlockChain) FindUTXO() map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs)
	spentTxos := make(map[string][]int)

	iter := chain.Iterator()
	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTxos[txID] != nil {
					for _, spentOut := range spentTxos[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					inTxID := hex.EncodeToString(in.ID)
					spentTxos[inTxID] = append(spentTxos[inTxID], in.Out)
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return UTXO
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey wallet.SerializablePrivateKey) {
	prevTxs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTx, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	tx.Sign(privKey, prevTxs)
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Iterator()

	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction does not exist")
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {

	if tx.IsCoinbase() {
		return true
	}

	prevTxs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTx, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	return tx.Verify(prevTxs)
}

func retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf("removing lock file: %v", err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	db, err := badger.Open(retryOpts)
	Handle(err)
	return db, nil
}

func openDB(dir string, opts badger.Options) (*badger.DB, error) {
	if db, err := badger.Open(opts); err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			if db, err := retry(dir, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				return db, nil
			}
			log.Println("Could not unlock db:", err)
		}
		return nil, err
	} else {
		return db, nil
	}
}
