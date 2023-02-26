package transference_test

import (
	"context"
	"core/transference"
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
		SetBalance(10).
		Save(context.Background())
}

func CreateTransferService(t *testing.T) (*transference.TransferService, *ent.Client, *zap.Logger) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	logger := zap.NewNop()
	mq := transference.MemoryTransferenceQueue{}

	return transference.NewTransferService(client, logger, mq), client, logger
}

func TestCreateTransferService(t *testing.T) {
	s, client, l := CreateTransferService(t)
	defer client.Close()
	defer l.Sync()

	assert.NotNil(t, s)
	assert.NotNil(t, client)
}
