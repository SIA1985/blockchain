package cli

import (
	"blockchain/internal/algorythms"
	"blockchain/internal/block"
	"blockchain/internal/blockchain"
	"blockchain/internal/transaction"
	"blockchain/internal/wallet"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"time"
)

func getBlockchain() *blockchain.Blockchain {
	/*todo: my address = 35fad8a91040ce20658c12e19b488c6e4d325edff153eadc4a8e79fab4bb403f7bf161f25a1f6382bb0f6dc58edfbf9e1d0389c7519a5c1c8da1becb00d71e00*/
	address, _ := hex.DecodeString("35fad8a91040ce20658c12e19b488c6e4d325edff153eadc4a8e79fab4bb403f7bf161f25a1f6382bb0f6dc58edfbf9e1d0389c7519a5c1c8da1becb00d71e00")

	bc, err := blockchain.NewBlockchain(address)
	if err != nil {
		panic(err)
	}

	return bc
}

type CLI struct {
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
	case "createWallet":
		c.CreateWallet()
	case "addresses":
		c.GetAddresses()
	default:
		c.PrintHelp()
	}
}

func (c *CLI) PrintBlockchain() {
	bc := getBlockchain()

	for b := range blockchain.ForEach(bc) {
		fmt.Println(b.StringHash())
		fmt.Println(time.UnixMilli(b.Header.Timestamp).String())
		/*todo:*/
		// fmt.Println(b.Data.Name)
		fmt.Println()
	}
}

func (c *CLI) AddBlock(from, to string, amount int64) {
	bc := getBlockchain()
	var err error

	srcAddress, err := hex.DecodeString(from)
	if err != nil {
		fmt.Printf("Can't decode 'from' address: %x\n", err)
		return
	}
	destAddress, err := hex.DecodeString(to)
	if err != nil {
		fmt.Printf("Can't decode 'to' address: %x\n", err)
		return
	}

	unspentOuts, accumulated := bc.FindOutputsToSpend(srcAddress, amount)

	tx := transaction.NewUTXOTransaction(srcAddress, destAddress, amount, unspentOuts, accumulated)
	err = bc.AddBlock(block.BlockData{Transactions: []*transaction.Transaction{tx}})

	if err != nil {
		fmt.Println("Can't send!")
		return
	}

	fmt.Printf("%s -> %s sended: %d\n", from, to, amount)
}

func (c *CLI) ValidateBlock(hash string) {
	bc := getBlockchain()

	for b := range blockchain.ForEach(bc) {
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
	bc := getBlockchain()

	if bc.ValidateBlocks() {
		fmt.Println("Blockchain is valid!")
	} else {
		fmt.Println("Blockchain is invalid!")
	}
}

func (c *CLI) GetBalance(address string) {
	bc := getBlockchain()

	srcAddress, err := hex.DecodeString(address)
	if err != nil {
		fmt.Printf("Can't decode address: %w\n", err)
		return
	}

	var balance int64 = 0
	for _, out := range bc.FindUTXO(srcAddress) {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (c *CLI) CreateWallet() {
	w := wallet.NewWallet()

	if err := w.Save(); err != nil {
		panic(err)
	}

	fmt.Printf("address: %x\n", w.GetAddress())
}

func (c *CLI) GetAddresses() {
	addresses, err := wallet.Addresses()
	if err != nil {
		panic(err)
	}

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (c *CLI) PrintHelp() {
	fmt.Println("Usage:")
	fmt.Println("addblock --data '...' - add block with current data")
	fmt.Println("print - print all blockchain")
	fmt.Println("validateAll - validate all blocks in blockchain")
	fmt.Println("validate --hash '...' - validate block by hash")
	fmt.Println("balance --address '...' - get balance by address")
	fmt.Println("createWallet - create new wallet")
	fmt.Println("addresses - all saved addresses")
}
