package storage

import "github.com/golang-jwt/jwt/v5"

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
}

func NewAgent() *StAgent {
	return &StAgent{}
}

// StClaims Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type StClaims struct {
	ID    int64  `json:"uid"`
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func NewClaims() *StClaims {
	return &StClaims{}
}
