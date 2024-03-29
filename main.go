package main

import (
	"GoBlockchain/cli"
	"os"
)

func main() {
	// VERY important to shutdown the app properly. Badger db must be properly
	// Garbage collected before terminating or else the data could corrupt!

	// By deferring these calls, we can run runtime.Goexit() to properly terminate the app

	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
