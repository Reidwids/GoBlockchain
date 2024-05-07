package blockchain

import (
	"GoBlockchain/internal/utils"
	"GoBlockchain/internal/wallet"
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
)

type TxOutput struct {
	Value      int    // value of output tokens in the tx. outputs cannot be split
	PubKeyHash []byte // the key needed to unlock the tokens
}

type TxOutputs struct {
	Outputs []TxOutput
}

type TxInput struct {
	ID        []byte // references the transaction ID the output is inside of
	Out       int    // the index the output appears
	Signature []byte // provides the data used in the outputs pub key
	PubKey    []byte
}

func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (tx *Transaction) ToString() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("Input # %d:", i))
		lines = append(lines, fmt.Sprintf("  Input TxID: %x", input.ID))
		lines = append(lines, fmt.Sprintf("  Out: %d", input.Out))
		lines = append(lines, fmt.Sprintf("  Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("  PubKey: %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("Output %d:", i))
		lines = append(lines, fmt.Sprintf("  Value: %d", output.Value))
		lines = append(lines, fmt.Sprintf("  PubKeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	Handle(err)
	return buffer.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs
	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	Handle(err)
	return outputs
}
