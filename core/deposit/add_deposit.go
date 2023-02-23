package deposit

import (
	"context"
	"core/utils"
	"errors"
	"models/ent"
	"models/ent/deposit"

	"github.com/google/uuid"
	"github.com/near/borsh-go"
	"go.uber.org/zap"
)

type AddDepositPayload struct {
	Amount float64
}

type AddDepositMessagePayload struct {
	DepositID uuid.UUID
}

func (s *DepositService) AddDeposit(ctx context.Context, user *ent.User, payload AddDepositPayload) error {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		s.logger.Error("failed to create transaction", zap.Error(err))
		return errors.New("failed to create transaction")
	}

	newDeposit, err := tx.Deposit.
		Create().
		SetUser(user).
		SetAmount(payload.Amount).
		SetStatus(deposit.StatusPending).
		Save(ctx)
	if err != nil {
		s.logger.Error("failed to create deposit", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to create deposit"))
	}

	msg := AddDepositMessagePayload{
		DepositID: newDeposit.ID,
	}
	data, err := borsh.Serialize(msg)
	if err != nil {
		s.logger.Error("failed to serialize", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to serialize"))
	}
	s.logger.Info(string(data))

	err = s.queue.AddMessageToQueue(ctx, data)
	if err != nil {
		s.logger.Error("failed to add message to queue", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to add message to queue"))
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("failed to commit transaction", zap.Error(err))
		return errors.New("failed to commit transaction")
	}

	return nil
}
