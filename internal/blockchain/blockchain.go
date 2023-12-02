package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"GoBlockchain/internal/blockchain/pb"

	"google.golang.org/protobuf/proto"
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

func createBlock(blockchain *Blockchain, proof int64, previousHash []byte) *Block {
	newBlock := &Block{
		index:        int64(len(blockchain.chain)),
		timestamp:    time.Now().Unix(),
		proof:        proof,
		prevHash:     previousHash,
		transactions: blockchain.transactions,
	}
	blockchain.chain = append(blockchain.chain, newBlock)
	blockchain.transactions = []*Transaction{}
	return newBlock
}

func getPreviousBlock(blockchain *Blockchain) *Block {
	return blockchain.chain[len(blockchain.chain)-1]
}

func proofOfWork(previousProof int64) int64 {
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

func hashBlock(Block *Block) []byte {
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

func isChainValid(blockchain *Blockchain, chain []*Block) bool {
	for i, block := range chain {
		if i > 0 {
			prevBlock := chain[i-1]
			// False if the previous block hash does not equal the current block hash
			if !bytes.Equal(block.prevHash, hashBlock(prevBlock)) {
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

func addTransaction(blockchain *Blockchain, sender string, recipient string, amount float64) {
	newTx := &Transaction{
		sender:    sender,
		recipient: recipient,
		amount:    amount,
	}
	blockchain.transactions = append(blockchain.transactions, newTx)
	print("Transaction successfully added!")
}

func addNode(blockchain *Blockchain, address string) {
	blockchain.nodes = append(blockchain.nodes, address)
	print("Node successfully added!")
}

func replaceChain(blockchain *Blockchain) bool {
	// longestChain := nil
	// maxlength := len(Blockchain.chain)
	for node := range blockchain.nodes {
		url := fmt.Sprintf("http://%s/getChain", node)
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("Error Connecting with node %s:", node, err)
			return false
		}
		defer response.Body.Close()

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return false
		}
		// Decode the protobuf-encoded data
		var receivedChain pb.Chain
		err = proto.Unmarshal(body, &receivedChain)
		if err != nil {
			fmt.Println("Error decoding protobuf data:", err)
			return false
		}

		// Access the received blockchain's chain
		for _, block := range receivedChain.Chain {
			fmt.Printf("Block %d\n", block.Index)
			for _, tx := range block.Transactions {
				fmt.Printf("  Transaction ID: %x\n", tx.ID)
				fmt.Printf("  Sender: %s\n", tx.Sender)
				fmt.Printf("  Recipient: %s\n", tx.Recipient)
				fmt.Printf("  Amount: %f\n", tx.Amount)
			}
		}
	}
	return true
}

// Improvements:
// 1. add fork resolution & detection
// 2. add consensus algorithm
// 3. add transaction propagation
// 5. add mining rewards
// 4. add wallet functionality
// 6. add transaction fees
// 7. add wallets and wallet balances
