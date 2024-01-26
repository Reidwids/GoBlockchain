package node

import (
	"GoBlockchain/internal/blockchain"
	"GoBlockchain/internal/proto/pb"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedNodeServiceServer
	Blockchain *blockchain.Blockchain
}

func (s *Server) GetNodesRes(ctx context.Context, req *pb.GetNodesReq) (*pb.NodeList, error) {
	return &pb.NodeList{Nodes: s.Blockchain.Nodes}, nil
}

func AddNode(blockchain *blockchain.Blockchain, address string) {
	blockchain.Nodes = append(blockchain.Nodes, address)
	fmt.Println("Node successfully added!")
}

func GetNodesReq() {
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

// broadcast block
// broadcast transaction
// broadcast node
// broadcast chain
// consensus
// Peer discovery
// fork resolution

// Func to query hardcoded nodes for their nodelist
// Func to add new node to nodelist

// func validateNodes() bool {}

// func broadcastNodes(blockchain *pb.Blockchain) {
// 	for _, node := range blockchain.Nodes {
// 		url := fmt.Sprintf("http://%s/getNodes", node)
// 		response, err :=  http.Get(url)
// 		struct NodeRes{
// 			Nodes string `json:"nodes"`
// 		}
// 		json.NewEncoder(w).Encode(response)
// 		res.Body.nodes
// 	}
// 	// for each node in node list, get nodes and calculate differences.
// 	// If there are new nodes in our list for them, add the new node and re broadcast
// 	// If there are new nodes for us, add the new nodes and re broadcast
// }
