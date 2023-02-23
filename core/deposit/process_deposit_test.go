package deposit_test

import (
	"context"
	"core/deposit"

	"models/ent"
	entd "models/ent/deposit"
	"models/ent/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createDeposit(client *ent.Client, user *ent.User, amount float64) (*ent.Deposit, error) {
	return client.Deposit.
		Create().
		SetUser(user).
		SetAmount(amount).
		SetStatus(entd.StatusPending).
		Save(context.Background())

}

func TestProcessDepositNotFound(t *testing.T) {
	s, client, l := CreateDepositService(t)
	defer client.Close()
	defer l.Sync()

	err := s.ProcessDeposit(context.Background(), deposit.AddDepositMessagePayload{
		DepositID: uuid.New(),
	}, 0.0)
	assert.Error(t, err)
}

func TestProcessDepositFraudDetected(t *testing.T) {
	s, client, l := CreateDepositService(t)
	defer client.Close()
	defer l.Sync()

	usr, err := generateAnyValidUser(client)
	assert.NotNil(t, usr)
	assert.Nil(t, err)

	dep, err := createDeposit(client, usr, 10)
	assert.NotNil(t, dep)
	assert.Nil(t, err)

	err = s.ProcessDeposit(context.Background(), deposit.AddDepositMessagePayload{
		DepositID: dep.ID,
	}, 100.0)
	assert.NoError(t, err)
}

func TestProcessDepositWithSuccess(t *testing.T) {
	s, client, l := CreateDepositService(t)
	defer client.Close()
	defer l.Sync()

	usr, err := generateAnyValidUser(client)
	assert.NotNil(t, usr)
	assert.Nil(t, err)
	assert.Equal(t, usr.Balance, float64(0))

	dep, err := createDeposit(client, usr, 10)
	assert.NotNil(t, dep)
	assert.Nil(t, err)

	err = s.ProcessDeposit(context.Background(), deposit.AddDepositMessagePayload{
		DepositID: dep.ID,
	}, 0.0)
	assert.NoError(t, err)

	updatedUsr, err := client.User.Query().Where(user.ID(usr.ID)).Only(context.Background())
	assert.NotNil(t, updatedUsr)
	assert.Nil(t, err)

	assert.Equal(t, updatedUsr.Balance, dep.Amount)
}
