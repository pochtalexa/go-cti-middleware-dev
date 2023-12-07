package models

type StAgent struct {
	ID       int64
	Login    string
	PassHash []byte
}

func NewAgent() *StAgent {
	return &StAgent{}
}
