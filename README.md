# Go Blockchain
Welcome to the Golang Blockchain! This was an exercise in learning Golang, through a fun DIY project - a small cryptocurrency implementation. This project is a local proof-of-work crypto currency, including wallets & transactions. Spin up a few terminals and try it out!

## How-to
1. Using an initial terminal (Lets call this terminal t1)- create a node ID, as the local host port to run the node on - `export NODE_ID=3000`
2. Create a wallet with `go run main.go createwallet`. Copy the addresses into a notepad - we'll be reusing them!
3. Copy the wallet address, and create a blockchain using the wallet as the target for rewards - `go run main.go createblockchain -address <t1-wallet-addr>`

To see the full capabilities of the current state, set up 2 more nodes:
1. Open 2 terminals (Lets call themn t2 and t3), and defined a new `NODE_ID` for each - `export NODE_ID=3001` `export NODE_ID=3002`
2. There is no blockchain seeding yet - so we must copy the initial blockchain - copy `/tmp/blocks_3000` twice, as `/tmp/blocks_3001` and `/tmp/blocks_3002` 
3. Create a wallet for the two nodes - `go run main.go createwallet` - Again, copy the addresses

Now, add a block to the original blockchain so we can see these 2 nodes sync when they start up.
1. In t1, commit a new tx and mine it with  - `go run main.go send -from <t1-wallet-addr> -to <t2-wallet-addr> -amount 10 -mine`
2. Start the t1 node - `go run main.go startnode`
3. Start the t2 node - `go run main.go startnode`. Notice that it copies the new blocks on startup!! The interaction logs can be seen in t1.
4. Start the t3 node as a miner - `go run main.go startnode -miner <t3-wallet-addr>`
5. Stop the t2 node, leaving the t2 terminal open. With an active miner in t3, we can commit transactions and watch them be mined by t3! Make a transaction in the t2 terminal, such as `go run main.go send -from <t2-wallet-addr> -to <t3-wallet-addr> -amount 10`
