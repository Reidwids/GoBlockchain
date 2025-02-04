package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
	Height       int
	Timestamp    int64
}

func CreateBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	// Create an empty block with the prev hash
	block := &Block{[]byte{}, txs, prevHash, 0, height, time.Now().Unix()}

	// Create a new pow instance from the new block and run the proof
	pow := NewProof(block)
	nonce, hash := pow.Run()

	// set the hash and nonce
	block.Hash = hash[:]
	fmt.Printf("Hash: %x\n", block.Hash)
	block.Nonce = nonce

	return block
}

func Genesis(coinbase *Transaction) *Block {
	// Init the chain with the genesis block
	return CreateBlock([]*Transaction{coinbase}, []byte{}, 0)
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}

	tree := NewMerkleTree(txHashes)
	return tree.RootNode.Data
}
