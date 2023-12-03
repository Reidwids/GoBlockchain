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

func NewBlockchain() *pb.Blockchain {
	genesisBlock := &pb.Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Proof:        0,
		PrevHash:     []byte{},
		Transactions: []*pb.Transaction{},
	}

	return &pb.Blockchain{
		Chain:        []*pb.Block{genesisBlock},
		Transactions: []*pb.Transaction{},
		Nodes:        []string{},
	}
}

func createBlock(blockchain *pb.Blockchain, proof int64, previousHash []byte) *pb.Block {
	newBlock := &pb.Block{
		Index:        int64(len(blockchain.Chain)),
		Timestamp:    time.Now().Unix(),
		Proof:        proof,
		PrevHash:     previousHash,
		Transactions: blockchain.Transactions,
	}
	blockchain.Chain = append(blockchain.Chain, newBlock)
	blockchain.Transactions = []*pb.Transaction{}
	return newBlock
}

func getPreviousBlock(blockchain *pb.Blockchain) *pb.Block {
	return blockchain.Chain[len(blockchain.Chain)-1]
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

func hashBlock(Block *pb.Block) []byte {
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

func isChainValid(chain []*pb.Block) bool {
	for i, block := range chain {
		if i > 0 {
			prevBlock := chain[i-1]
			// False if the previous block hash does not equal the current block hash
			if !bytes.Equal(block.PrevHash, hashBlock(prevBlock)) {
				return false
			}

			// False if the proof does not start with 0000
			proofHash := hashProof(block.Proof, prevBlock.Proof)
			if proofHash[:4] != "0000" {
				return false
			}
		}
	}
	return true
}

func addTransaction(blockchain *pb.Blockchain, sender string, recipient string, amount float32) {
	newTx := &pb.Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	blockchain.Transactions = append(blockchain.Transactions, newTx)
	print("Transaction successfully added!")
}

func addNode(blockchain *pb.Blockchain, address string) {
	blockchain.Nodes = append(blockchain.Nodes, address)
	print("Node successfully added!")
}

func replaceChain(blockchain *pb.Blockchain) error {
	var longestChain []*pb.Block
	maxlength := len(blockchain.Chain)
	for _, node := range blockchain.Nodes {
		nodeChain, err := getChainFromNode(node)

		if err != nil {
			return err
		}
		// Access the received blockchain's chain
		if len(nodeChain) > maxlength && isChainValid(nodeChain) {
			maxlength = len(nodeChain)
			longestChain = nodeChain
		}
	}
	if longestChain != nil {
		blockchain.Chain = longestChain
	}

	return nil
}

func getChainFromNode(node string) ([]*pb.Block, error) {
	url := fmt.Sprintf("http://%s/getChain", node)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error Connecting with node:", node, err)
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	// Decode the protobuf-encoded data
	var chainBuf pb.Chain
	err = proto.Unmarshal(body, &chainBuf)
	if err != nil {
		fmt.Println("Error decoding protobuf data:", err)
		return nil, err
	}

	receivedChain := chainBuf.Chain

	return receivedChain, nil
}

// Improvements:
// 1. add fork resolution & detection
// 2. add consensus algorithm
// 3. add transaction propagation
// 5. add mining rewards
// 4. add wallet functionality
// 6. add transaction fees
// 7. add wallets and wallet balances
