package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"
)

type Transaction struct {
	ID        []byte
	sender    string
	recipient string
	amount    float64
}

type Block struct {
	index        int64
	timestamp    int64
	proof        int64
	prevHash     []byte
	transactions []*Transaction
}

type Blockchain struct {
	chain        []*Block
	transactions []*Transaction
	nodes        []string
}

func NewBlockchain() *Blockchain {
	genesisBlock := &Block{
		index:        0,
		timestamp:    time.Now().Unix(),
		proof:        0,
		prevHash:     []byte{},
		transactions: []*Transaction{},
	}

	return &Blockchain{
		chain:        []*Block{genesisBlock},
		transactions: []*Transaction{},
		nodes:        []string{},
	}
}

func createBlock(Blockchain *Blockchain, proof int64, previousHash []byte) *Block {
	newBlock := &Block{
		index:        int64(len(Blockchain.chain)),
		timestamp:    time.Now().Unix(),
		proof:        proof,
		prevHash:     previousHash,
		transactions: Blockchain.transactions,
	}
	Blockchain.chain = append(Blockchain.chain, newBlock)
	Blockchain.transactions = []*Transaction{}
	return newBlock
}

func getPreviousBlock(Blockchain *Blockchain) *Block {
	return Blockchain.chain[len(Blockchain.chain)-1]
}

func proofOfWork(Blockchain *Blockchain, previousProof int64) int64 {
	newProof := int64(1)
	checkProof := false

	for !checkProof {
		proofHash := hashProof(newProof, previousProof)

		if proofHash[:4] == "0000" {
			checkProof = true
		} else {
			newProof++
		}
	}
	return newProof
}

func hashBlock(Blockchain *Blockchain, Block *Block) []byte {
	encodedBlock := sha256.Sum256([]byte(fmt.Sprintf("%v", Block)))
	return []byte(hex.EncodeToString(encodedBlock[:]))
}

func hashProof(newProof int64, prevProof int64) string {
	// Take the hash of the difference of squares between the 2 proof vals
	// To create a simple proof of work algorithm
	hashInput := math.Pow(float64(newProof), 2) - math.Pow(float64(prevProof), 2)
	hashBytes := sha256.Sum256([]byte(fmt.Sprintf("%f", hashInput)))
	return hex.EncodeToString(hashBytes[:])
}

func isChainValid(Blockchain *Blockchain, chain []*Block) bool {
	for i, block := range chain {
		if i > 0 {
			prevBlock := chain[i-1]
			// False if the previous block hash does not equal the current block hash
			if !bytes.Equal(block.prevHash, hashBlock(Blockchain, prevBlock)) {
				return false
			}

			// False if the proof does not start with 0000
			proofHash := hashProof(block.proof, prevBlock.proof)
			if proofHash[:4] != "0000" {
				return false
			}
		}
	}
	return true
}

func addTransaction(Blockchain *Blockchain, sender string, recipient string, amount float64) {
	newTx := &Transaction{
		sender:    sender,
		recipient: recipient,
		amount:    amount,
	}
	Blockchain.transactions = append(Blockchain.transactions, newTx)
	print("Transaction successfully added!")
}

func addNode(Blockchain *Blockchain, address string) {
	Blockchain.nodes = append(Blockchain.nodes, address)
	print("Node successfully added!")
}

func replaceChain(Blockchain *Blockchain) bool {
	longestChain := nil
	maxlength := len(Blockchain.chain)
	for node := range Blockchain.nodes {
		url := fmt.Sprintf("http://%s/getChain", node)
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("Error Connecting with node %s:", node, err)
			return false
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		// res :=
	}
}
