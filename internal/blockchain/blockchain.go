package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"time"

	"GoBlockchain/internal/transactions"
)

type Blockchain struct {
	Chain        []*Block
	Transactions []*transactions.Transaction
	Nodes        []string
}

type Block struct {
	Index        int64
	Timestamp    int64
	Proof        int64
	PrevHash     []byte
	Transactions []*transactions.Transaction
}

func NewBlockchain() *Blockchain {
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Proof:        0,
		PrevHash:     []byte{},
		Transactions: []*transactions.Transaction{},
	}

	return &Blockchain{
		Chain:        []*Block{genesisBlock},
		Transactions: []*transactions.Transaction{},
		Nodes:        []string{"localhost:3001", "localhost:3002"},
	}
}

func CreateBlock(blockchain *Blockchain, proof int64, previousHash []byte) *Block {
	newBlock := &Block{
		Index:        int64(len(blockchain.Chain)),
		Timestamp:    time.Now().Unix(),
		Proof:        proof,
		PrevHash:     previousHash,
		Transactions: blockchain.Transactions,
	}
	blockchain.Chain = append(blockchain.Chain, newBlock)
	blockchain.Transactions = []*transactions.Transaction{}
	return newBlock
}

func GetPreviousBlock(blockchain *Blockchain) *Block {
	return blockchain.Chain[len(blockchain.Chain)-1]
}

func ProofOfWork(previousProof int64) int64 {
	newProof := int64(1)
	checkProof := false

	for !checkProof {
		proofHash := HashProof(newProof, previousProof)

		if proofHash[:4] == "0000" {
			checkProof = true
		} else {
			newProof++
		}
	}
	return newProof
}

func HashBlock(Block *Block) []byte {
	encodedBlock := sha256.Sum256([]byte(fmt.Sprintf("%v", Block)))
	return []byte(hex.EncodeToString(encodedBlock[:]))
}

func HashProof(newProof int64, prevProof int64) string {
	// Take the hash of the difference of squares between the 2 proof vals
	// To create a simple proof of work algorithm
	hashInput := math.Pow(float64(newProof), 2) - math.Pow(float64(prevProof), 2)
	hashBytes := sha256.Sum256([]byte(fmt.Sprintf("%f", hashInput)))
	return hex.EncodeToString(hashBytes[:])
}

func IsChainValid(chain []*Block) bool {
	for i, block := range chain {
		if i > 0 {
			prevBlock := chain[i-1]
			// False if the previous block hash does not equal the current block hash
			if !bytes.Equal(block.PrevHash, HashBlock(prevBlock)) {
				return false
			}

			// False if the proof does not start with 0000
			proofHash := HashProof(block.Proof, prevBlock.Proof)
			if proofHash[:4] != "00000" {
				return false
			}
		}
	}
	return true
}

func AddTransaction(blockchain *Blockchain, sender string, recipient string, amount float64) {
	newTx := &transactions.Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	blockchain.Transactions = append(blockchain.Transactions, newTx)
	fmt.Println("Transaction successfully added!")
}

func ReplaceChain(blockchain *Blockchain) (bool, error) {
	chainReplaced := false
	var longestChain []*Block
	maxlength := len(blockchain.Chain)
	for _, node := range blockchain.Nodes {
		nodeChain, err := GetChainFromNode(node)

		if err != nil {
			return chainReplaced, err
		}
		// Access the received blockchain's chain
		if len(nodeChain) > maxlength && IsChainValid(nodeChain) {
			maxlength = len(nodeChain)
			longestChain = nodeChain
		}
	}
	if longestChain != nil {
		blockchain.Chain = longestChain
		chainReplaced = true
	}

	return chainReplaced, nil
}

func GetChainFromNode(node string) ([]*Block, error) {
	url := fmt.Sprintf("http://%s/getChain", node)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error Connecting with node:", node, err)
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	return nil, err
	// }
	// Decode the protobuf-encoded data
	var chainBuf []*Block
	// err = proto.Unmarshal(body, &chainBuf)
	// if err != nil {
	// 	fmt.Println("Error decoding protobuf data:", err)
	// 	return nil, err
	// }

	// receivedChain := chainBuf

	// return receivedChain, nil
	return chainBuf, nil
}

// Improvements:
// 1. add fork resolution & detection
// 2. add consensus algorithm
// 3. add transaction propagation
// 5. add mining rewards
// 4. add wallet functionality
// 6. add transaction fees
// 7. add wallets and wallet balances
