package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"vlady-kotsev/blocky/blockchain"
)

const (
	addCmd   = "add"
	printCmd = "print"
)

type CommandLine struct {
	blockchain *blockchain.Blockchain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Print the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	err := cli.blockchain.AddBlock(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("Block added")
}

func (cli *CommandLine) printBlockchain() {
	err := cli.blockchain.PrintBlocks()
	if err != nil {
		panic(err)
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case addCmd:
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	case printCmd:
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			cli.printUsage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printBlockchain()
	}
}
