package deposit_test

import (
	"context"
	"core/deposit"
	"errors"
	"fmt"
	"math/rand"
	"models/ent"
	"models/ent/enttest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func generateAnyValidUser(c *ent.Client) (*ent.User, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error while hashing password")
	}

	someEmail := fmt.Sprintf("some-%d@mail.com", rand.Int())

	return c.User.
		Create().
		SetName("some name").
		SetEmail(someEmail).
		SetPassword(string(pwdHash)).
		SetBalance(0).
		Save(context.Background())
}

func CreateDepositService(t *testing.T) (*deposit.DepositService, *ent.Client, *zap.Logger) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	logger := zap.NewNop()
	mq := deposit.MemoryDepositQueue{}

	return deposit.NewDepositService(client, logger, mq), client, logger
}

func TestCreateDepositService(t *testing.T) {
	s, client, l := CreateDepositService(t)
	defer client.Close()
	defer l.Sync()

	assert.NotNil(t, s)
	assert.NotNil(t, client)
}
