package blockchain

import (
	"GoBlockchain/internal/wallet"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
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

func InitBlockChain(address string) *BlockChain {
	// Create a new blockchain
	var lastHash []byte

	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.EventLogging = false
	opts.Logger = nil
	db, err := badger.Open(opts)
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

func ContinueBlockChain(address string) *BlockChain {
	if !DBexists() {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}
	var lastHash []byte

	opts := badger.DefaultOptions("")
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.EventLogging = false
	opts.Logger = nil
	db, err := badger.Open(opts)
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

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) *Block {
	// Create a new block with the given data and the previous hash
	var lastHash []byte

	// View is a read only badger db transaction
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		// Get the prev hash from the db
		lastHash, err = item.ValueCopy(lastHash)
		return err
	})

	Handle(err)

	newBlock := CreateBlock(transactions, lastHash)

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

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	// Create an iterator instance to navigate through the blockchain
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	var encodedBlock []byte
	// Get the previous block from the db and decode it
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err = item.ValueCopy(encodedBlock)
		block = Deserialize(encodedBlock)
		return err
	})
	Handle(err)

	// Set the current hash to current block's prev hash
	iter.CurrentHash = block.PrevHash
	return block
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
	prevTxs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTx, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	return tx.Verify(prevTxs)
}
