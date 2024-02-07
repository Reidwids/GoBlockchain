package blockchain

import "fmt"

type BlockChain struct {
	Blocks []*Block
}

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

func (chain *BlockChain) AddBlock(data string) {
	// Create a new block with the given data and the previous hash
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)

	// Append the new block to the chain
	chain.Blocks = append(chain.Blocks, new)
}

func Genesis() *Block {
	// Init the chain with the genesis block
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	// Create a new blockchain with an instance of a genesis block
	return &BlockChain{[]*Block{Genesis()}}
}
