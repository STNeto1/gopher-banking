package transference

import (
	"context"
	"core/utils"
	"errors"
	"math/rand"
	"models/ent/transference"

	"go.uber.org/zap"
)

func (s *TransferService) ProcessTransfer(ctx context.Context, payload TransferenceToProcessPayload, offset float64) error {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		s.logger.Error("failed to create transaction", zap.Error(err))
		return errors.New("failed to create transaction")
	}

	dbTransfer, err := tx.Transference.Query().Where(transference.ID(payload.TransferID)).WithToUser().WithFromUser().Only(ctx)
	if err != nil {
		s.logger.Error("failed to fetch transference", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to fetch transference"))
	}

	if dbTransfer.Status != transference.StatusPending {
		s.logger.Error("transference already processed", zap.Error(err))

		err = tx.Commit()
		if err != nil {
			s.logger.Error("failed to commit transaction", zap.Error(err))
			return errors.New("failed to commit transaction")
		}

		return nil
	}

	randNumber := rand.Float64() * 100
	if randNumber < offset {
		// "fraud" detected
		_, err = tx.Transference.UpdateOne(dbTransfer).SetStatus(transference.StatusDenied).Save(ctx)
		if err != nil {
			s.logger.Error("failed to update transference", zap.Error(err))
			return utils.Rollback(tx, errors.New("failed to update transference"))
		}

		_, err = tx.User.UpdateOne(dbTransfer.Edges.FromUser).SetBalance(dbTransfer.Amount + dbTransfer.Edges.FromUser.Balance).Save(ctx)
		if err != nil {
			s.logger.Error("failed to update user balance", zap.Error(err))
			return utils.Rollback(tx, errors.New("failed to update user balance"))
		}

		err = tx.Commit()
		if err != nil {
			s.logger.Error("failed to commit transaction", zap.Error(err))
			return errors.New("failed to commit transaction")
		}

		return nil
	}

	_, err = tx.Transference.UpdateOne(dbTransfer).SetStatus(transference.StatusCompleted).Save(ctx)
	if err != nil {
		s.logger.Error("failed to update transference", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to update transference"))
	}

	_, err = tx.User.UpdateOne(dbTransfer.Edges.ToUser).SetBalance(dbTransfer.Amount + dbTransfer.Edges.ToUser.Balance).Save(ctx)
	if err != nil {
		s.logger.Error("failed to update user balance", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to update user balance"))
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("failed to commit transaction", zap.Error(err))
		return errors.New("failed to commit transaction")
	}

	return nil
}
