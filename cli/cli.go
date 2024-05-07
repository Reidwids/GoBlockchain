package cli

import (
	"GoBlockchain/internal/blockchain"
	"GoBlockchain/internal/utils"
	"GoBlockchain/internal/wallet"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address of the wallet")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address of the wallet to mine the genesis block")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || sendAmount == nil {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - Get the balance for the given address")
	fmt.Println(" createblockchain -address ADDRESS - Creates a blockchain, mines genesis block with given address")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send tokens from one address to another")
	fmt.Println(" createwallet - Creates a new Wallet")
	fmt.Println(" listaddresses - Lists the addresses in our wallet file")
	fmt.Println(" reindexutxo - Rebuilds to UTXO set")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()

		runtime.Goexit()
	}
}

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool((pow.Validate())))
		fmt.Println()

		for _, tx := range block.Transactions {
			fmt.Println(tx.ToString())
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Blockchain Created!")
}

func (cli *CommandLine) getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not valid")
	}

	chain := blockchain.ContinueBlockChain(address)
	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	UTXOSet.Reindex()
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUnspentTransactions(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	if !wallet.ValidateAddress(to) {
		log.Panic("To address is not valid")
	}

	if !wallet.ValidateAddress(from) {
		log.Panic("From address is not valid")
	}

	chain := blockchain.ContinueBlockChain(from)
	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, &UTXOSet)
	block := chain.AddBlock([]*blockchain.Transaction{tx})
	UTXOSet.Update(block)
	fmt.Println("Send was successful!")
}

func (cli *CommandLine) listAddresses() {
	wallets, _ := wallet.CreateWalletsMap()
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CommandLine) reindexUTXO() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}

func (cli *CommandLine) createWallet() {
	wallets, _ := wallet.CreateWalletsMap()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New wallet created successfully. Address is: %s\n", address)
}
