package transaction

type TxOutput struct {
	Value  int
	Pubkey string
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.Pubkey == data
}
