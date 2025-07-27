package cli

import (
	"blockchain/internal/algorythms"
	"blockchain/internal/blockchain"
	"flag"
	"fmt"
	"os"
	"time"
)

type CLI struct {
	Bc *blockchain.Blockchain
}

func (c *CLI) Run() {
	switch os.Args[1] {
	case "addblock":
		cmd := flag.NewFlagSet("addblock", flag.ExitOnError)
		data := cmd.String("data", "", "data of block")
		cmd.Parse(os.Args[2:])
		c.AddBlock(*data)
	case "print":
		c.PrintBlockchain()
	case "validateAll":
		c.ValidateBlockchain()
	case "validate":
		cmd := flag.NewFlagSet("validate", flag.ExitOnError)
		hash := cmd.String("hash", "", "hash of block")
		cmd.Parse(os.Args[2:])
		c.ValidateBlock(*hash)
	case "balance":
		cmd := flag.NewFlagSet("balance", flag.ExitOnError)
		address := cmd.String("address", "me", "get balance for this address")
		cmd.Parse(os.Args[2:])
		c.GetBalance(*address)
	default:
		c.PrintHelp()
	}
}

func (c *CLI) PrintBlockchain() {
	for b := range blockchain.ForEach(c.Bc) {
		fmt.Println(b.StringHash())
		fmt.Println(time.UnixMilli(b.Header.Timestamp).String())
		/*todo:*/
		// fmt.Println(b.Data.Name)
		fmt.Println()
	}
}

func (c *CLI) AddBlock(name string) {
	// err := c.Bc.AddBlock(block.BlockData{Name: name})
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func (c *CLI) ValidateBlock(hash string) {
	for b := range blockchain.ForEach(c.Bc) {
		if b.StringHash() != hash {
			continue
		}

		valid := algorythms.Validate(b.PrepareForValidate(), b.Header.TargetBits)
		if valid {
			fmt.Println("Block is valid!")
		} else {
			fmt.Println("Block is invalid!")
		}

		return
	}
}

func (c *CLI) ValidateBlockchain() {
	if c.Bc.ValidateBlocks() {
		fmt.Println("Blockchain is valid!")
	} else {
		fmt.Println("Blockchain is invalid!")
	}
}

func (c *CLI) GetBalance(address string) {
	var balance int64 = 0
	for _, out := range c.Bc.FindUTXO(address) {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (c *CLI) PrintHelp() {
	fmt.Println("Usage:")
	fmt.Println("addblock --data '...' - add block with current data")
	fmt.Println("print - print all blockchain")
	fmt.Println("validateAll - validate all blocks in blockchain")
	fmt.Println("validate --hash '...' - validate block by hash")
}
