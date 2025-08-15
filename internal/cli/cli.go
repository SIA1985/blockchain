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
	/*todo: my address = 3031a099fab71777d092760196c49d4a72fe69a922fd37e31a17*/
	address, _ := hex.DecodeString("3031a099fab71777d092760196c49d4a72fe69a922fd37e31a17")

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

/*todo: пул транзакций*/
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

	unspentOuts, accumulated, err := bc.FindOutputsToSpend(srcAddress, amount)
	if err != nil {
		fmt.Printf("Can't find utxo: %x\n", err)
		return
	}

	var TXs []*transaction.Transaction
	/*todo: взять самые выгодные транзакции из пула*/
	TXs = append(TXs, transaction.NewUTXOTransaction(srcAddress, destAddress, amount, unspentOuts, accumulated))

	err = bc.AddBlock(block.BlockData{Transactions: TXs}, accumulated-amount)
	if err != nil {
		fmt.Printf("Can't send: %x\n", err)
		return
	}

	err = bc.UpdateUTXO(TXs)
	if err != nil {
		fmt.Printf("Can't update utxo: %x\n", err)
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
		fmt.Printf("Can't decode address: %x\n", err)
		return
	}

	UTXO, err := bc.FindUTXO(srcAddress)
	if err != nil {
		return
	}

	var balance int64 = 0
	for _, outs := range UTXO {
		for _, out := range outs {
			balance += out.Value
		}
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (c *CLI) CreateWallet() {
	w := wallet.NewWallet()

	if err := w.Save(); err != nil {
		panic(err)
	}

	fmt.Printf("address: %x\n", w.Address())
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
