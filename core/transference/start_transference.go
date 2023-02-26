package transference

import (
	"context"
	"core/utils"
	"errors"
	"models/ent"
	"models/ent/transference"
	"models/ent/user"

	"github.com/google/uuid"
	"github.com/near/borsh-go"
	"go.uber.org/zap"
)

type StartTransferencePayload struct {
	Amount  float64
	ToUser  uuid.UUID
	Message *string
}

type TransferenceToProcessPayload struct {
	TransferID uuid.UUID
}

func (s *TransferService) StartTransfer(ctx context.Context, usr *ent.User, payload StartTransferencePayload) error {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		s.logger.Error("failed to create transaction", zap.Error(err))
		return errors.New("failed to create transaction")
	}

	to_user, err := s.client.User.Query().Where(user.ID(payload.ToUser)).Only(ctx)
	if err != nil {
		s.logger.Error("from user was not found", zap.Error(err))
		return utils.Rollback(tx, errors.New("from user was not found"))
	}

	if usr.Balance < payload.Amount {
		s.logger.Error("insufficient balance", zap.Error(err))
		return utils.Rollback(tx, errors.New("insufficient balance"))
	}

	newTransfer, err := tx.Transference.
		Create().
		SetFromUser(usr).
		SetToUser(to_user).
		SetAmount(payload.Amount).
		SetStatus(transference.StatusPending).
		SetNillableMessage(payload.Message).
		Save(ctx)
	if err != nil {
		s.logger.Error("failed to create transference", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to create transference"))
	}

	_, err = tx.User.UpdateOne(usr).SetBalance(usr.Balance - payload.Amount).Save(ctx)
	if err != nil {
		s.logger.Error("failed to update user balance", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to update user balance"))
	}

	msg := TransferenceToProcessPayload{
		TransferID: newTransfer.ID,
	}
	data, err := borsh.Serialize(msg)
	if err != nil {
		s.logger.Error("failed to serialize", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to serialize"))
	}

	err = s.queue.AddTransferenceToQueue(ctx, data)
	if err != nil {
		s.logger.Error("failed to add message to transference queue", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to add message to transference queue"))
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("failed to commit transaction", zap.Error(err))
		return errors.New("failed to commit transaction")
	}

	return nil
}
