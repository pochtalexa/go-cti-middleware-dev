package flags

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseFlagsDefault(t *testing.T) {

	defaultServAddr := "http://localhost:9595"
	defaultLogin := "agent"
	defaultPassword := "123"
	defaultRegister := false

	ParseFlags()

	assert.Equal(t, ServAddr, defaultServAddr)
	assert.Equal(t, Login, defaultLogin)
	assert.Equal(t, Password, defaultPassword)
	assert.Equal(t, Register, defaultRegister)
}

func TestParseFlags(t *testing.T) {

	testServAddr := "http://localhost:3535"
	testLogin := "agent1"
	testPassword := "1234"
	testRegister := true

	os.Args = []string{"test", "-a", testServAddr, "-l", testLogin, "-p", testPassword, "-r", "true"}

	ParseFlags()

	assert.Equal(t, ServAddr, testServAddr)
	assert.Equal(t, Login, testLogin)
	assert.Equal(t, Password, testPassword)
	assert.Equal(t, Register, testRegister)
}
