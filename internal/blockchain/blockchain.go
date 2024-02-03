package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"GoBlockchain/internal/transactions"
)

type Blockchain struct {
	Chain        []*Block
	Transactions []*transactions.Transaction
	mutex        sync.Mutex
}

type Block struct {
	Index        int64
	Timestamp    int64
	Proof        int64
	Hash         []byte
	PrevHash     []byte
	Transactions []*transactions.Transaction
	Difficulty   int
	Nonce        string
}

func NewBlockchain() *Blockchain {
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		PrevHash:     []byte{},
		Transactions: []*transactions.Transaction{},
		Difficulty:   4,
		Nonce:        "",
	}

	minedBlock := ProofOfWork(genesisBlock)

	return &Blockchain{
		Chain:        []*Block{minedBlock},
		Transactions: []*transactions.Transaction{},
	}
}

func HashBlock(block *Block) []byte {
	record := fmt.Sprintf("%d%d%x%v%d%s",
		block.Index,
		block.Timestamp,
		block.PrevHash,
		block.Transactions,
		block.Difficulty,
		block.Nonce)

	encodedBlock := sha256.Sum256([]byte(fmt.Sprintf("%v", record)))
	return []byte(hex.EncodeToString(encodedBlock[:]))
}

func isHashValid(hash []byte, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(string(hash), prefix)
}

func IsChainValid(chain []*Block) bool {
	for i, block := range chain {
		if i > 0 {
			prevBlock := chain[i-1]
			// False if the previous block hash does not equal the current block hash
			if !bytes.Equal(block.PrevHash, HashBlock(prevBlock)) {
				return false
			}

			// False if the proof does not start with difficulty number of 0s
			proofHash := HashBlock(prevBlock)
			if !isHashValid(proofHash, prevBlock.Difficulty) {
				return false
			}
		}
	}
	return true
}

func ProofOfWork(newBlock *Block) *Block {
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex

		hash := HashBlock(newBlock)
		if isHashValid(hash, newBlock.Difficulty) {
			fmt.Println(string(hash), " work done!")
			newBlock.Hash = hash
			break
		}
	}
	return newBlock
}

func MineBlock(bc *Blockchain, nodeID string) *Block {
	prevBlock := bc.Chain[len(bc.Chain)-1]

	AddTransaction(bc, "The Network", nodeID, 1)
	newBlock := &Block{
		Index:        int64(len(bc.Chain)),
		Timestamp:    time.Now().Unix(),
		PrevHash:     prevBlock.Hash,
		Difficulty:   prevBlock.Difficulty,
		Transactions: bc.Transactions,
		Nonce:        "",
	}
	minedBlock := ProofOfWork(newBlock)

	bc.mutex.Lock()
	bc.Chain = append(bc.Chain, minedBlock)
	bc.Transactions = []*transactions.Transaction{}
	bc.mutex.Unlock()
	return newBlock
}

func AddTransaction(bc *Blockchain, sender string, recipient string, amount float32) {
	newTx := &transactions.Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	bc.mutex.Lock()
	bc.Transactions = append(bc.Transactions, newTx)
	bc.mutex.Unlock()
	fmt.Println("New Transaction: ", newTx)
}

// func IsBlockValid(block *Block, prevBlock *Block) bool {
// 	// False if the previous block hash does not equal the current block hash
// 	if !bytes.Equal(block.PrevHash, HashBlock(prevBlock)) {
// 		return false
// 	}

// 	// False if the proof does not start with 0000
// 	proofHash := HashBlock(block)
// 	return isHashValid(proofHash, block.Difficulty)
// }

// func ReplaceChain(blockchain *Blockchain) (bool, error) {
// 	chainReplaced := false
// 	var longestChain []*Block
// 	maxlength := len(blockchain.Chain)
// 	for _, node := range node.Nodes {
// 		nodeChain, err := GetChainFromNode(node)

// 		if err != nil {
// 			return chainReplaced, err
// 		}
// 		// Access the received blockchain's chain
// 		if len(nodeChain) > maxlength && IsChainValid(nodeChain) {
// 			maxlength = len(nodeChain)
// 			longestChain = nodeChain
// 		}
// 	}
// 	if longestChain != nil {
// 		blockchain.Chain = longestChain
// 		chainReplaced = true
// 	}

// 	return chainReplaced, nil
// }

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
