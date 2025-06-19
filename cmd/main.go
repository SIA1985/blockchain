package main

import (
	"blockchain/internal/blockchain"
	"blockchain/internal/cli"
)

func main() {
	bc, err := blockchain.NewBlockchain()
	if err != nil {
		panic(err)
	}

	c := cli.CLI{Bc: bc}

	c.Run()
}
