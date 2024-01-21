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
	blockchainKey  ContextKey = "blockchain"
	nodeAddressKey ContextKey = "nodeAddress"
)

func Start() {
	fmt.Println("Starting Go Blockchain...")
	bc := blockchain.NewBlockchain()
	nodeId := uuid.New()

	// Use a wait group to launch both servers
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		startgRPC(bc, nodeId)
		waitGroup.Done()
	}()

	go func() {
		starthttp(bc, nodeId)
		waitGroup.Done()
	}()

	waitGroup.Wait()
}

func starthttp(bc *blockchain.Blockchain, nodeId uuid.UUID) {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/mineBlock", withBlockchain(handleMineBlock, bc, nodeId))
	http.HandleFunc("/getChain", withBlockchain(handleGetChain, bc, nodeId))
	http.HandleFunc("/isChainValid", withBlockchain(handleIsChainValid, bc, nodeId))
	http.HandleFunc("/addTransaction", withBlockchain(handleAddTransaction, bc, nodeId))
	http.HandleFunc("/connectNode", withBlockchain(handleConnectNode, bc, nodeId))
	http.HandleFunc("/replaceChain", withBlockchain(handleReplaceChain, bc, nodeId))
	fmt.Println("Go Blockchain http running on port 3003...")
	http.ListenAndServe(":3003", nil)
}

func startgRPC(bc *blockchain.Blockchain, nodeId uuid.UUID) {
	lis, err := net.Listen("tcp", ":3005")
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

func withBlockchain(next http.HandlerFunc, bc *blockchain.Blockchain, nodeId uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add 'bc' to the request context
		ctx := context.WithValue(r.Context(), blockchainKey, bc)
		ctx = context.WithValue(ctx, nodeAddressKey, nodeId)
		next(w, r.WithContext(ctx))
	}
}

func testgRPC() {
	conn, err := grpc.Dial("localhost:3004", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)
	res, err := client.GetNodes(context.Background(), &pb.GetNodesReq{})
	if err != nil {
		log.Fatalf("Failed to get nodes: %v", err)
	}
	fmt.Println(res.Nodes)
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
	nodeId := r.Context().Value(nodeAddressKey).(uuid.UUID)

	prevBlock := blockchain.GetPreviousBlock(bc)
	prevProof := prevBlock.Proof
	previousHash := blockchain.HashBlock(prevBlock)
	newProof := blockchain.ProofOfWork(prevProof)
	blockchain.AddTransaction(bc, "The Network", nodeId.String(), 1)
	block := blockchain.CreateBlock(bc, newProof, previousHash)
	response := struct {
		Message   string `json:"message"`
		Index     int64  `json:"index"`
		Timestamp int64  `json:"timestamp"`
		PrevHash  []byte `json:"prevHash"`
	}{
		Message:   "Block mined",
		Index:     block.Index,
		Timestamp: block.Timestamp,
		PrevHash:  block.PrevHash,
	}
	json.NewEncoder(w).Encode(response)
}

func handleGetChain(w http.ResponseWriter, r *http.Request) {
	// bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)
	response := struct {
		Chain []*pb.Block `json:"chain"`
	}{
		// Chain: blockchain.Blockchain.Chain,
	}
	json.NewEncoder(w).Encode(response)
}

func handleIsChainValid(w http.ResponseWriter, r *http.Request) {
	bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)
	isValid := blockchain.IsChainValid(bc.Chain)
	if isValid {
		response := struct {
			IsValid     bool `json:"isValid"`
			ChainLength int  `json:"chainLength"`
		}{
			IsValid:     isValid,
			ChainLength: len(bc.Chain),
		}
		json.NewEncoder(w).Encode(response)
	}
}

func handleAddTransaction(w http.ResponseWriter, r *http.Request) {
	// bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)
	var tx pb.Transaction
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tx); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	// blockchain.AddTransaction(bc, tx.Sender, tx.Recipient, tx.Amount)

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Transaction added to pool",
	}
	json.NewEncoder(w).Encode(response)
}

func handleConnectNode(w http.ResponseWriter, r *http.Request) {
	bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)
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
	node.AddNode(bc, body.NodeAddress)

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Succesfully added node",
	}
	json.NewEncoder(w).Encode(response)
}

func handleReplaceChain(w http.ResponseWriter, r *http.Request) {
	bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)

	chainReplaced, err := blockchain.ReplaceChain(bc)
	if err != nil {
		http.Error(w, "Error replacing chain", http.StatusBadRequest)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: func() string {
			if chainReplaced {
				return "Chain was replaced"
			}
			return "Chain was not replaced"
		}(),
	}
	json.NewEncoder(w).Encode(response)
}

func getNodes(w http.ResponseWriter, r *http.Request) {
	bc := r.Context().Value(blockchainKey).(*blockchain.Blockchain)

	response := struct {
		Nodes []string `json:"nodes"`
	}{
		Nodes: bc.Nodes,
	}
	json.NewEncoder(w).Encode(response)
}
