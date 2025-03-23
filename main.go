package main

import (
	"vlady-kotsev/blocky/blockchain"

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
	chain, err := blockchain.InitBlockchain(db)
	if err != nil {
		panic(err)
	}
	cli := CommandLine{
		blockchain: chain,
	}

	cli.run()
}
