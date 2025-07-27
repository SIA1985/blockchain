package main

import (
	"blockchain/internal/blockchain"
	"blockchain/internal/cli"
)

func main() {
	/*todo: address*/
	bc, err := blockchain.NewBlockchain("me")
	if err != nil {
		panic(err)
	}

	c := cli.CLI{Bc: bc}

	c.Run()
}
