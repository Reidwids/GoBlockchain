package blockchain

type TxOutput struct {
	Value  int    // value of output tokens in the tx. outputs cannot be split
	PubKey string // the key needed to unlock the tokens
}

type TxInput struct {
	ID  []byte // references the transaction ID the output is inside of
	Out int    // the index the output appears
	Sig string // provides the data used in the outputs pub key
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
