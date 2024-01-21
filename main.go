package main

import (
	"GoBlockchain/internal/api"
)

func main() {
	// Define command-line flags
	api.Start()

	// createBlockchain := flag.Bool("create", false, "Create a new blockchain")
	// addBlock := flag.String("addblock", "", "Add a new block to the blockchain")
	// listBlocks := flag.Bool("list", false, "List all blocks in the blockchain")
	// startApi := flag.Bool("api", false, "Start the API server")
	// flag.Parse()
	// Process command-line flags
	// if *createBlockchain {
	// 	// Create a new blockchain
	// 	fmt.Println("Creating a new blockchain...")
	// } else if *addBlock != "" {
	// 	// Add a new block to the blockchain
	// 	fmt.Println("Adding a new block to the blockchain...")
	// 	// bc.AddBlock(*addBlock)
	// } else if *listBlocks {
	// 	// List all blocks in the blockchain
	// 	fmt.Println("Listing all blocks in the blockchain:")
	// 	// for _, block := range bc.Blocks {
	// 	// 	fmt.Printf("Block %d: %s\n", block.Index, block.Data)
	// 	// }
	// } else if *startApi {
	// 	// Start the API server
	// 	fmt.Println("Starting the API server...")
	// 	api.Start()

	// } else {
	// 	// No valid command-line flags provided
	// 	fmt.Println("Usage:")
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }
}
