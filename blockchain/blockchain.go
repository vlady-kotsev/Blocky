package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

const LashHashKey = "lh"

type Blockchain struct {
	LastHash []byte
	DB       *badger.DB
	Pow      *ProofOfWork
}

type BlockchainIterator struct {
	CurrentHash []byte
	DB          *badger.DB
}

func InitBlockchain(db *badger.DB) (*Blockchain, error) {
	pow := CreateProofOfWork()
	chain := Blockchain{Pow: pow, DB: db}

	err := chain.DB.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(LashHashKey)); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := CreateGenesis()
			genesisSerialized, err := genesis.SerializeBlock()
			if err != nil {
				return err
			}
			err = txn.Set(genesis.Hash, genesisSerialized)
			if err != nil {
				return err
			}
			err = txn.Set([]byte(LashHashKey), genesis.Hash)
			if err != nil {
				return err
			}
			chain.LastHash = genesis.Hash
		} else {
			item, err := txn.Get([]byte(LashHashKey))
			if err != nil {
				return err
			}
			var dbLastHash []byte
			err = item.Value(func(val []byte) error {
				dbLastHash = val
				return nil
			})
			chain.LastHash = dbLastHash
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &chain, nil
}

func (bc *Blockchain) CreateBlock(data string, prevHash []byte) *Block {
	block := Block{
		Data:     []byte(data),
		PrevHash: prevHash,
	}
	nonce, hash := bc.Pow.Run(&block)

	block.Hash = hash
	block.Nonce = nonce

	return &block
}

func (bc *Blockchain) AddBlock(data string) error {
	var dbLashHash []byte
	err := bc.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LashHashKey))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			dbLashHash = val
			return nil
		})

		return nil
	})
	if err != nil {
		return err
	}

	// create block
	newBlock := bc.CreateBlock(data, dbLashHash)

	// validate block
	if !bc.Pow.Validate(newBlock) {
		return fmt.Errorf("Invalid block")
	}

	// add block
	err = bc.DB.Update(func(txn *badger.Txn) error {
		newBlockSerialized, err := newBlock.SerializeBlock()
		if err != nil {
			return err
		}
		err = txn.Set(newBlock.Hash, newBlockSerialized)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(LashHashKey), newBlock.Hash)
		if err != nil {
			return err
		}
		// update in meory blockchain
		bc.LastHash = newBlock.Hash
		return nil
	})

	return nil
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	iter := BlockchainIterator{
		CurrentHash: bc.LastHash,
		DB:          bc.DB,
	}
	return &iter
}

func (bi *BlockchainIterator) Next() *Block {
	var block *Block
	err := bi.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(bi.CurrentHash)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			block, err = DeserializeBlock(val)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// log error
		return nil
	}
	bi.CurrentHash = block.PrevHash
	return block
}

func (bc *Blockchain) IterateBlockchain() ([]*Block, error) {
	iter := bc.Iterator()

	var blocks []*Block
	for !IsEqualToZeroHash(iter.CurrentHash) {
		currentBlock := iter.Next()

		blocks = append(blocks, currentBlock)
	}
	return blocks, nil
}

func (bc *Blockchain) PrintBlocks() error {
	blocks, err := bc.IterateBlockchain()
	if err != nil {
		return err
	}
	for _, block := range blocks {
		fmt.Printf("Hash: %x, Data: %s, PrevHash: %x\n", block.Hash, block.Data, block.PrevHash)
	}
	return nil
}
