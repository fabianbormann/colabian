package main

import (
	"flag"
	"fmt"
	"os"

	"colabian/core"
	"colabian/persistence"
)

func main() {
	if len(os.Args) < 2 {
        fmt.Println("Usage:", os.Args[0], "[command]")
		fmt.Println("Available commands:")
		fmt.Println("\tmine:\t\tMine a new block")
		fmt.Println("\tsummarize:\tSummarize the blockchain")
        return
    }

	mineFlagSet := flag.NewFlagSet("mine", flag.ExitOnError)
    dataFlag := mineFlagSet.String("data", "", "The block data")

	command := os.Args[1]
	blockchain := persistence.Load("blockchain.dat")

	if (command == "mine") {
		mineFlagSet.Parse(os.Args[2:])

		if *dataFlag == "" {
			fmt.Println("Block data cannot be empty")
			fmt.Println("Usage:", os.Args[0], "mine -data \"some data\"")
			return
		}

		blockchain = core.Mine(*dataFlag, blockchain)

		persistence.Save(blockchain, "", "1.0")
	} else if (command == "summarize") {
		core.Summarize(blockchain)
	} else {
		fmt.Println("Unknown command", "\""+command+"\"", "use one of: mine, summarize")
	}
}
