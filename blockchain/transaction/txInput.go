package transaction

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}
