package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"vlady-kotsev/blocky/blockchain/transaction"

	"github.com/dgraph-io/badger/v4"
)

const LashHashKey = "lh"

type Blockchain struct {
	LastHash []byte
	DB       *badger.DB
	Pow      *ProofOfWork
}

func InitBlockchain(db *badger.DB, address string) (*Blockchain, error) {
	pow := CreateProofOfWork()
	chain := Blockchain{Pow: pow, DB: db}

	err := chain.DB.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(LashHashKey)); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			coinbaseTx, err := transaction.CoinbaseTx(address, GenesisTxData)
			if err != nil {
				return err
			}
			genesis := CreateGenesis(coinbaseTx)
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

func (bc *Blockchain) CreateBlock(txs []*transaction.Transaction, prevHash []byte) *Block {
	block := Block{
		Transactions: txs,
		PrevHash:     prevHash,
	}
	nonce, hash := bc.Pow.Run(&block)

	block.Hash = hash
	block.Nonce = nonce

	return &block
}

func CreateGenesis(coinbase *transaction.Transaction) *Block {
	zeroHash := [32]byte{}
	data := []byte(GenesisData)
	bytes := bytes.Join([][]byte{
		data,
		zeroHash[:],
	}, []byte{})

	hash := sha256.Sum256(bytes)
	return &Block{
		Transactions: []*transaction.Transaction{
			coinbase,
		},
		PrevHash: zeroHash[:],
		Hash:     hash[:],
		Nonce:    0,
	}
}

func (bc *Blockchain) AddBlock(txs []*transaction.Transaction) error {
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
	newBlock := bc.CreateBlock(txs, dbLashHash)

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
		fmt.Printf("Hash: %x\nData: %x\nPrevHash: %x\n\n", block.Hash, block.HashTransactions(), block.PrevHash)
	}
	return nil
}

func (bc *Blockchain) FindUnspentTransactions(address string) []*transaction.Transaction {
	var unspentTxs []*transaction.Transaction

	spentTXOs := make(map[string][]int)

	iter := bc.Iterator()

	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}
		if IsEqualToZeroHash(block.PrevHash) {
			break
		}
	}

	return unspentTxs
}

func (bc *Blockchain) FindUTXO(address string) []transaction.TxOutput {
	var UTXOs []transaction.TxOutput
	unspentTransactions := bc.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		for _, output := range tx.Outputs {
			if output.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, output)
			}
		}
	}
	return UTXOs
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
	Outputs:
		for outIndex, output := range tx.Outputs {
			if output.CanBeUnlocked(address) && accumulated < amount {
				accumulated += output.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIndex)
			}

			if accumulated >= amount {
				break Outputs
			}
		}

	}
	return accumulated, unspentOuts
}

func (bc *Blockchain) NewTransaction(from, to string, amount int) (*transaction.Transaction, error) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	accumulated, validOutputs := bc.FindSpendableOutputs(from, amount)
	if amount > accumulated {
		return nil, fmt.Errorf("Not enough balance")
	}

	for txID, outputs := range validOutputs {
		id, err := hex.DecodeString(txID)
		if err != nil {
			return nil, err
		}
		for _, outIdx := range outputs {
			input := transaction.TxInput{ID: id, Out: outIdx, Sig: from}
			inputs = append(inputs, input)
		}

	}

	outputs = append(outputs, transaction.TxOutput{
		Value:  amount,
		Pubkey: to,
	})

	if accumulated > amount {
		outputs = append(outputs, transaction.TxOutput{
			Value:  accumulated - amount,
			Pubkey: from,
		})
	}

	tx := transaction.Transaction{Inputs: inputs, Outputs: outputs}
	tx.SetID()

	return &tx, nil
}
