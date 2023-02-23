package deposit_test

import (
	"context"
	"core/deposit"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddDeposit(t *testing.T) {
	s, client, l := CreateDepositService(t)
	defer client.Close()
	defer l.Sync()

	usr, err := generateAnyValidUser(client)
	assert.NotNil(t, usr)
	assert.Nil(t, err)

	err = s.AddDeposit(context.Background(), usr, deposit.AddDepositPayload{
		Amount: 10,
	})
	assert.NoError(t, err)
}
