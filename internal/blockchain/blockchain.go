package blockchain

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID        []byte
	sender    string
	recipient string
	amount    int
}

type Block struct {
	Index        int64
	Timestamp    int64
	proof        int64
	PrevHash     []byte
	transactions []*Transaction
}

type Blockchain struct {
	Chain        []*Block
	transactions []*Transaction
	nodes        map[string]bool
}

func NewBlockchain() *Blockchain {
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		proof:        0,
		PrevHash:     []byte{},
		transactions: []*Transaction{},
	}

	return &Blockchain{
		Chain:        []*Block{genesisBlock},
		transactions: []*Transaction{},
		nodes:        make(map[string]bool),
	}
}

func (b *Blockchain) Method2() {
	fmt.Println("Method2 called")

}
