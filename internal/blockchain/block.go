package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	// Create an empty block with the prev hash
	block := &Block{[]byte{}, []byte(data), prevHash, 0}

	// Create a new pow instance from the new block and run the proof
	pow := NewProof(block)
	nonce, hash := pow.Run()

	// set the hash and nonce
	block.Hash = hash[:]
	fmt.Printf("Hash: %x\n", block.Hash)
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	// Init the chain with the genesis block
	return CreateBlock("Genesis", []byte{})
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
