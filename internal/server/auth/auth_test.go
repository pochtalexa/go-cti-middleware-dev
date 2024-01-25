package auth

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/mock"
)

import (
	"github.com/pochtalexa/go-cti-middleware/internal/server/auth/mocks"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
)

func TestLogin(t *testing.T) {
	type args struct {
		login      string
		password   string
		curStorage IntStorage
	}
	tests := []struct {
		name     string
		args     args
		tokenTTL int64
		want     string
		wantErr  bool
	}{
		{
			name: "Login ok",
			args: args{
				login:    "agent",
				password: "123",
			},
			tokenTTL: 1,
			wantErr:  false,
		},
		{
			name: "Login err",
			args: args{
				login:    "agent",
				password: "1234",
			},
			tokenTTL: 1,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Init()
			config.ServerConfig.TokenTTL = time.Duration(tt.tokenTTL) * time.Minute

			testStorage := mocks.NewIntStorage(t)
			testStorage.
				On("GetAgent", tt.args.login).
				Once().
				Return(&storage.StAgent{
					ID:       1,
					Login:    tt.args.login,
					PassHash: []byte("$2a$10$skQdruyEZ54IOFPHNUd/wejcZZSbwzZM13lp4TKN2fQ26j6OLKCAS"),
				}, nil)

			got, err := Login(tt.args.login, tt.args.password, testStorage)
			fmt.Println("got:", got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if got != tt.want {
			//	t.Errorf("Login() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestRegisterNewUser(t *testing.T) {
	type args struct {
		login      string
		password   string
		curStorage IntStorage
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "RegisterNewUser ok",
			args: args{
				login:    "agent",
				password: "123",
			},
			wantErr: false,
		},
		{
			name: "RegisterNewUser err",
			args: args{
				login:    "agent",
				password: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStorage := mocks.NewIntStorage(t)

			if !strings.Contains(tt.name, "err") {
				testStorage.
					On("SaveAgent", tt.args.login, mock.Anything).
					Once().
					Return(int64(1), nil)
			}

			got, err := RegisterNewUser(tt.args.login, tt.args.password, testStorage)
			fmt.Println("got:", got)

			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//if got != tt.want {
			//	t.Errorf("RegisterNewUser() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
