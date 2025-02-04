package blockchain

import "github.com/dgraph-io/badger"

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
