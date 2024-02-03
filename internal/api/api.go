package api

import (
	"GoBlockchain/internal/blockchain"
	"GoBlockchain/internal/node"
	"GoBlockchain/internal/proto/pb"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type ContextKey string

const (
	blockchainKey ContextKey = "blockchain"
	nodeKey       ContextKey = "nodeKey"
)

func Start() {
	fmt.Println("Starting Go Blockchain...")
	bc := blockchain.NewBlockchain()
	n := &node.Node{
		ID:      uuid.New().String(),
		Address: "localhost:3003",
		Peers:   []string{},
	}

	// Use a wait group to launch both servers
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		startgRPC(bc, n)
		waitGroup.Done()
	}()

	go func() {
		starthttp(bc, n)
		waitGroup.Done()
	}()

	waitGroup.Wait()
}

func starthttp(bc *blockchain.Blockchain, n *node.Node) {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/mineBlock", withContext(handleMineBlock, bc, n))
	http.HandleFunc("/getChain", withContext(handleGetChain, bc, n))
	http.HandleFunc("/isChainValid", withContext(handleIsChainValid, bc, n))
	http.HandleFunc("/addTransaction", withContext(handleAddTransaction, bc, n))
	http.HandleFunc("/connectNode", withContext(handleConnectNode, bc, n))
	// http.HandleFunc("/replaceChain", withBlockchain(handleReplaceChain, bc, n))
	fmt.Println("Go Blockchain http running on port 3003...")
	http.ListenAndServe(":3000", nil)
}

func startgRPC(bc *blockchain.Blockchain, n *node.Node) {
	lis, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gRPCserver := grpc.NewServer()
	nodeServer := node.Server{
		Blockchain: bc,
	}

	pb.RegisterNodeServiceServer(gRPCserver, &nodeServer)

	fmt.Println("Go Blockchain gRPC running on port 3004...")
	if err := gRPCserver.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server on port 3004: %v", err)
	}
}

func withContext(next http.HandlerFunc, bc *blockchain.Blockchain, n *node.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add 'bc' to the request context
		ctx := context.WithValue(r.Context(), blockchainKey, bc)
		ctx = context.WithValue(ctx, nodeKey, n)
		next(w, r.WithContext(ctx))
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Message: "Golang Blockchain",
		Version: "v1.3.0",
	}
	json.NewEncoder(w).Encode(response)
}

func handleMineBlock(w http.ResponseWriter, r *http.Request) {
	bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)
	n := r.Context().Value(nodeKey).(*node.Node)

	block := blockchain.MineBlock(bc, n.ID)
	response := struct {
		Message   string `json:"message"`
		Index     int64  `json:"index"`
		Timestamp int64  `json:"timestamp"`
		PrevHash  string `json:"prevHash"`
	}{
		Message:   "Block mined",
		Index:     block.Index,
		Timestamp: block.Timestamp,
		PrevHash:  string(block.PrevHash),
	}
	json.NewEncoder(w).Encode(response)
}

func handleGetChain(w http.ResponseWriter, r *http.Request) {
	bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)
	response := struct {
		Chain []*blockchain.Block `json:"chain"`
	}{
		Chain: bc.Chain,
	}
	json.NewEncoder(w).Encode(response)
}

func handleIsChainValid(w http.ResponseWriter, r *http.Request) {
	bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)
	isValid := blockchain.IsChainValid(bc.Chain)
	response := struct {
		IsValid     bool `json:"isValid"`
		ChainLength int  `json:"chainLength"`
	}{
		IsValid:     isValid,
		ChainLength: len(bc.Chain),
	}
	json.NewEncoder(w).Encode(response)
}

func handleAddTransaction(w http.ResponseWriter, r *http.Request) {
	bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)
	var tx pb.Transaction
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tx); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	blockchain.AddTransaction(bc, tx.Sender, tx.Recipient, tx.Amount)

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Transaction added to pool",
	}
	json.NewEncoder(w).Encode(response)
}

func handleConnectNode(w http.ResponseWriter, r *http.Request) {
	// bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)
	type ConnectNodeRq struct {
		NodeAddress string `json:"nodeAddress"`
	}
	var body ConnectNodeRq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	// Get requesting node's IP address
	// node.AddNode(bc, body.NodeAddress)

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Succesfully added node",
	}
	json.NewEncoder(w).Encode(response)
}

// func handleReplaceChain(w http.ResponseWriter, r *http.Request) {
// 	bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)

// 	chainReplaced, err := blockchain.ReplaceChain(bc)
// 	if err != nil {
// 		http.Error(w, "Error replacing chain", http.StatusBadRequest)
// 		return
// 	}

// 	response := struct {
// 		Message string `json:"message"`
// 	}{
// 		Message: func() string {
// 			if chainReplaced {
// 				return "Chain was replaced"
// 			}
// 			return "Chain was not replaced"
// 		}(),
// 	}
// 	json.NewEncoder(w).Encode(response)
// }

func getNodes(w http.ResponseWriter, r *http.Request) {
	n := r.Context().Value(nodeKey).(*node.Node)

	response := struct {
		Nodes []string `json:"nodes"`
	}{
		Nodes: n.Peers,
	}
	json.NewEncoder(w).Encode(response)
}
