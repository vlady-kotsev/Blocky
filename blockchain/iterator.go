package blockchain

import "github.com/dgraph-io/badger/v4"

type BlockchainIterator struct {
	CurrentHash []byte
	DB          *badger.DB
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	err := iter.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
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
	iter.CurrentHash = block.PrevHash
	return block
}
