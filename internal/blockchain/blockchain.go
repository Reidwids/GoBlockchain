package blockchain

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	// Create a new blockchain

	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	db, err := badger.Open(opts)
	Handle(err)

	// Our db will hold 2 types of kv pairs - an "lh" / hash pair to store our last hash,
	// And hash / block pairs to store and retrieve each block
	err = db.Update(func(txn *badger.Txn) error {
		// check for lh (last hash) in the database. If it doesn't exist, create the genesis block
		if item, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println(("No existing blockchain found"))
			genesis := Genesis()
			fmt.Println("Genesis proof complete")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash
			return err
		} else {
			// If lh is found, pull the last hash from the db
			Handle(err)
			lastHash, err = item.ValueCopy(lastHash)
			return err
		}
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
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

	newBlock := CreateBlock(data, lastHash)

	// Update the db with a new kv entry for the new hash + block
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
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
