package transference_test

import (
	"context"
	"core/transference"

	"models/ent"
	entt "models/ent/transference"
	"models/ent/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createTransference(client *ent.Client, from_usr *ent.User, to_usr *ent.User, amount float64) (*ent.Transference, error) {
	return client.Transference.
		Create().
		SetFromUser(from_usr).
		SetToUser(to_usr).
		SetAmount(amount).
		SetStatus(entt.StatusPending).
		Save(context.Background())
}

func TestProcessDepositNotFound(t *testing.T) {
	s, client, l := CreateTransferService(t)
	defer client.Close()
	defer l.Sync()

	err := s.ProcessTransfer(context.Background(), transference.TransferenceToProcessPayload{
		TransferID: uuid.New(),
	}, 0.0)
	assert.Error(t, err)
}

func TestProcessDepositAlreadyProcessed(t *testing.T) {
	s, client, l := CreateTransferService(t)
	defer client.Close()
	defer l.Sync()

	from_usr, err := generateAnyValidUser(client)
	assert.NotNil(t, from_usr)
	assert.Nil(t, err)

	to_usr, err := generateAnyValidUser(client)
	assert.NotNil(t, to_usr)
	assert.Nil(t, err)

	tr, err := client.Transference.
		Create().
		SetFromUser(from_usr).
		SetToUser(to_usr).
		SetAmount(10).
		SetStatus(entt.StatusCompleted).
		Save(context.Background())
	assert.NotNil(t, tr)
	assert.Nil(t, err)

	err = s.ProcessTransfer(context.Background(), transference.TransferenceToProcessPayload{
		TransferID: tr.ID}, 0.0)
	assert.NoError(t, err)
}

func TestProcessDepositFraudDetected(t *testing.T) {
	s, client, l := CreateTransferService(t)
	defer client.Close()
	defer l.Sync()

	from_usr, err := generateAnyValidUser(client)
	assert.NotNil(t, from_usr)
	assert.Nil(t, err)

	to_usr, err := generateAnyValidUser(client)
	assert.NotNil(t, to_usr)
	assert.Nil(t, err)

	tr, err := createTransference(client, from_usr, to_usr, 10)
	assert.NotNil(t, tr)
	assert.Nil(t, err)

	err = s.ProcessTransfer(context.Background(), transference.TransferenceToProcessPayload{
		TransferID: tr.ID}, 100.0)
	assert.NoError(t, err)
}

func TestProcessDepositWithSuccess(t *testing.T) {
	s, client, l := CreateTransferService(t)
	defer client.Close()
	defer l.Sync()

	from_usr, err := generateAnyValidUser(client)
	assert.NotNil(t, from_usr)
	assert.Nil(t, err)

	to_usr, err := generateAnyValidUser(client)
	assert.NotNil(t, to_usr)
	assert.Nil(t, err)

	tr, err := createTransference(client, from_usr, to_usr, 10)
	assert.NotNil(t, tr)
	assert.Nil(t, err)

	err = s.ProcessTransfer(context.Background(), transference.TransferenceToProcessPayload{
		TransferID: tr.ID}, 100.0)
	assert.NoError(t, err)

	updatedToUsr, err := client.User.Query().Where(user.ID(to_usr.ID)).Only(context.Background())
	assert.NotNil(t, updatedToUsr)
	assert.Nil(t, err)

	assert.Equal(t, updatedToUsr.Balance, 20)
}
