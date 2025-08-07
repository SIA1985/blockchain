package cli

import (
	"blockchain/internal/algorythms"
	"blockchain/internal/block"
	"blockchain/internal/blockchain"
	"blockchain/internal/transaction"
	"encoding/hex"
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
	case "send":
		cmd := flag.NewFlagSet("send", flag.ExitOnError)
		from := cmd.String("from", "me", "address from")
		to := cmd.String("to", "me", "address to")
		amount := cmd.Int64("amount", 0, "amount to send")
		cmd.Parse(os.Args[2:])
		c.AddBlock(*from, *to, *amount)
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

func (c *CLI) AddBlock(from, to string, amount int64) {
	var err error

	srcAddress, err := hex.DecodeString(from)
	if err != nil {
		fmt.Printf("Can't decode address: %w\n", err)
		return
	}
	destAddress, err := hex.DecodeString(to)
	if err != nil {
		fmt.Printf("Can't decode address: %w\n", err)
		return
	}

	unspentOuts, accumulated := c.Bc.FindOutputsToSpend(srcAddress, amount)

	tx := transaction.NewUTXOTransaction(srcAddress, destAddress, amount, unspentOuts, accumulated)
	err = c.Bc.AddBlock(block.BlockData{Transactions: []*transaction.Transaction{tx}})

	if err != nil {
		fmt.Println("Can't send!")
		return
	}

	fmt.Printf("%s -> %s sended: %d\n", from, to, amount)
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
	srcAddress, err := hex.DecodeString(address)
	if err != nil {
		fmt.Printf("Can't decode address: %w\n", err)
		return
	}

	var balance int64 = 0
	for _, out := range c.Bc.FindUTXO(srcAddress) {
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
