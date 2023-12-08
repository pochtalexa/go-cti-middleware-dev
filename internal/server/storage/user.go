package storage

type StCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewCredentials() *StCredentials {
	return &StCredentials{}
}

type StRegisterSuccess struct {
	Id int64 `json:"id"`
}

func NewRegisterSuccess() *StRegisterSuccess {
	return &StRegisterSuccess{}
}

type StLoginSuccess struct {
	Token string `json:"token"`
}

func NewLoginSuccess() *StLoginSuccess {
	return &StLoginSuccess{}
}

type StAgent struct {
	ID       int64
	Login    string
	PassHash []byte
	TokenTTL int64
}

func NewAgent() *StAgent {
	return &StAgent{}
}
