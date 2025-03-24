package main

import (
	"vlady-kotsev/blocky/blockchain"
	"vlady-kotsev/blocky/cli"

	"github.com/dgraph-io/badger/v4"
)

func main() {
	opt := badger.DefaultOptions(blockchain.DbPath)
	opt.Logger = nil
	db, err := badger.Open(opt)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	address := "123"
	chain, err := blockchain.InitBlockchain(db, address)
	if err != nil {
		panic(err)
	}
	cli := cli.NewCLI(chain)

	err = cli.Run()
	if err != nil {
		panic(err)
	}
}
