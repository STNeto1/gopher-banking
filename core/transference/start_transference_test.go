package transference_test

import (
	"context"
	"core/transference"
	"models/ent/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStartTransferToUserNotFound(t *testing.T) {
	s, client, l := CreateTransferService(t)
	defer client.Close()
	defer l.Sync()

	usr, err := generateAnyValidUser(client)
	assert.NotNil(t, usr)
	assert.Nil(t, err)

	err = s.StartTransfer(context.Background(), usr, transference.StartTransferencePayload{
		ToUser:  uuid.New(),
		Message: nil,
		Amount:  10,
	})
	assert.Error(t, err)
}

func TestStartTransferInsufficientBalance(t *testing.T) {
	s, client, l := CreateTransferService(t)
	defer client.Close()
	defer l.Sync()

	usr1, err := generateAnyValidUser(client)
	assert.NotNil(t, usr1)
	assert.Nil(t, err)

	usr2, err := generateAnyValidUser(client)
	assert.NotNil(t, usr2)
	assert.Nil(t, err)

	err = s.StartTransfer(context.Background(), usr1, transference.StartTransferencePayload{
		ToUser:  usr2.ID,
		Message: nil,
		Amount:  20,
	})
	assert.Error(t, err)
}

func TestStartTransferWithSuccess(t *testing.T) {
	s, client, l := CreateTransferService(t)
	defer client.Close()
	defer l.Sync()

	usr1, err := generateAnyValidUser(client)
	assert.NotNil(t, usr1)
	assert.Nil(t, err)

	usr2, err := generateAnyValidUser(client)
	assert.NotNil(t, usr2)
	assert.Nil(t, err)

	err = s.StartTransfer(context.Background(), usr1, transference.StartTransferencePayload{
		ToUser:  usr2.ID,
		Message: nil,
		Amount:  5,
	})
	assert.NoError(t, err)

	updatedUser1, err := client.User.Query().Where(user.ID(usr1.ID)).Only(context.Background())
	assert.NotNil(t, updatedUser1)
	assert.NoError(t, err)
	assert.Equal(t, float64(5), updatedUser1.Balance)

	updatedUser2, err := client.User.Query().Where(user.ID(usr2.ID)).Only(context.Background())
	assert.NotNil(t, updatedUser2)
	assert.NoError(t, err)
	assert.Equal(t, float64(10), updatedUser2.Balance) // should not update
}
