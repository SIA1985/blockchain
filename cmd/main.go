package main

import (
	"blockchain/internal/blockchain"
	"blockchain/internal/cli"
)

func main() {
	/*todo: address*/
	bc, err := blockchain.NewBlockchain([]byte("meeeeeeeeeeeeeeeeeeeeee"))
	if err != nil {
		panic(err)
	}

	c := cli.CLI{Bc: bc}

	c.Run()
}
