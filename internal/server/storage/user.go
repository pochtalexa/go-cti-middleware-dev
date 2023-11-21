package storage

type StRegister struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type StRegisterSuccess struct {
	Id int64 `json:"id"`
}

func NewRegister() *StRegister {
	return &StRegister{}
}

func NewRegisterSuccess() *StRegisterSuccess {
	return &StRegisterSuccess{}
}
