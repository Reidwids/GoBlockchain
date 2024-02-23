package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 15

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	// By left shifting 1 by (256 - Difficulty), we get a target that is missing the first n bits
	// This will be the target for the proof of work, where we will try to find a hash that is less
	// than the target - meaning it has the first n bits set to 0
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			// Could maybe take out the ToHex function in the future
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Validate() bool {
	// Validate that a proof of work meets the target requirements
	var intHash big.Int
	// Create a new data slice with the provided nonce
	data := pow.InitData(pow.Block.Nonce)

	// Create a new sha256 hash from the data
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	// Check that the hash from the given data is less than the target
	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	// Use a buffer to write the binary data to
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	// Return the bytes from the buffer
	return buff.Bytes()
}

func (pow *ProofOfWork) Run() (int, []byte) {

	var intHash big.Int
	var hash [32]byte
	nonce := 0

	// Increment the nonce until we find a hash that is less than the target
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256((data))
		fmt.Printf("\r%x", hash)

		// Set the hash to the big.Int
		intHash.SetBytes(hash[:])

		// Check if the hash is less than the target
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}
