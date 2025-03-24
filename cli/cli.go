package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"vlady-kotsev/blocky/blockchain"
	"vlady-kotsev/blocky/blockchain/transaction"
)

type CommandLine struct {
	blockchain *blockchain.Blockchain
}

func NewCLI(blockchain *blockchain.Blockchain) *CommandLine {
	return &CommandLine{blockchain: blockchain}
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("	balance -address <ADDRESS> - Get the balance of the address")
	fmt.Println("	print-chain - Print the blocks in the chain")
	fmt.Println("	send -from <FROM> -to <TO> -amount <AMOUNT> - Sends an amount from address to another")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) balance(address string) {
	utxos := cli.blockchain.FindUTXO(address)

	amount := 0
	for _, utxo := range utxos {
		amount += utxo.Value
	}

	fmt.Printf("Address: %s\nBalance: %d\n", address, amount)
}

func (cli *CommandLine) send(from, to string, amount int) error {
	tx, err := cli.blockchain.NewTransaction(from, to, amount)
	if err != nil {
		return err
	}

	err = cli.blockchain.AddBlock([]*transaction.Transaction{tx})
	if err != nil {
		return err
	}

	return nil
}

// func (cli *CommandLine) addBlock(data string) {
// 	err := cli.blockchain.AddBlock([]*transaction.Transaction{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Block added")
// }

func (cli *CommandLine) printBlockchain() {
	err := cli.blockchain.PrintBlocks()
	if err != nil {
		panic(err)
	}
}

func (cli *CommandLine) Run() error {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	getBalanceData := getBalanceCmd.String("address", "", "Address of the user")

	printChainCmd := flag.NewFlagSet("print-chain", flag.ExitOnError)

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "The sender address")
	sendTo := sendCmd.String("to", "", "The receiver address")
	sendAmount := sendCmd.Int("amount", 0, "The amount to send")

	switch os.Args[1] {
	case "balance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}
	case "print-chain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			cli.printUsage()
			runtime.Goexit()
		}
		cli.balance(*getBalanceData)
	} else if printChainCmd.Parsed() {
		cli.printBlockchain()
	} else if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == 0 {
			cli.printUsage()
			runtime.Goexit()
		}
		err := cli.send(*sendFrom, *sendTo, *sendAmount)
		if err != nil {
			return err
		}
	}
	return nil
}
