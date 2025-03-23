package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

type ProofOfWork struct {
	Target     *big.Int
	Difficulty uint64
}

func (pow *ProofOfWork) AdjustDifficulty(value uint64) {
	pow.Difficulty = value
	target := big.NewInt(1)
	target.Lsh(target, uint(MaxBytes-pow.Difficulty))
	pow.Target = target
}

func CreateProofOfWork() *ProofOfWork {
	pow := ProofOfWork{Difficulty: InitialDifficulty}

	target := big.NewInt(1)
	target.Lsh(target, uint(MaxBytes-InitialDifficulty))
	pow.Target = target

	return &pow
}

func (pow *ProofOfWork) InitData(block *Block, nonce uint64) []byte {
	data := bytes.Join([][]byte{
		block.Data,
		block.PrevHash,
		IntToBytes(InitialDifficulty),
		IntToBytes(nonce),
	}, []byte{})

	return data
}

func (pow *ProofOfWork) Run(block *Block) (uint64, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := uint64(0)
	for nonce < math.MaxUint64 {
		data := pow.InitData(block, nonce)

		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])

		fmt.Printf("%d %x \n", nonce, hash)

		if intHash.Cmp(pow.Target) == -1 {
			return nonce, hash[:]
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate(block *Block) bool {
	var intHash big.Int

	data := pow.InitData(block, block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1

}
