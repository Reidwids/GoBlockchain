package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	// Create an empty block with the prev hash
	block := &Block{[]byte{}, txs, prevHash, 0}

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
	return CreateBlock([]*Transaction{coinbase}, []byte{})
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
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	// Concat all txHashes together and take the hash of the result
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}
