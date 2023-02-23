package deposit

import (
	"context"
	"core/utils"
	"errors"
	"math/rand"
	"models/ent/deposit"

	"go.uber.org/zap"
)

func (s *DepositService) ProcessDeposit(ctx context.Context, payload AddDepositMessagePayload, offset float64) error {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		s.logger.Error("failed to create transaction", zap.Error(err))
		return errors.New("failed to create transaction")
	}

	dbDeposit, err := tx.Deposit.Query().Where(deposit.ID(payload.DepositID)).WithUser().Only(ctx)
	if err != nil {
		s.logger.Error("failed to fetch deposit", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to fetch deposit"))
	}

	randNumber := rand.Float64() * 100
	if randNumber < offset {
		// "fraud" detected
		_, err = tx.Deposit.UpdateOne(dbDeposit).SetStatus(deposit.StatusDenied).Save(ctx)
		if err != nil {
			s.logger.Error("failed to update deposit", zap.Error(err))
			return utils.Rollback(tx, errors.New("failed to update deposit"))
		}

		err = tx.Commit()
		if err != nil {
			s.logger.Error("failed to commit transaction", zap.Error(err))
			return errors.New("failed to commit transaction")
		}

		return nil
	}

	_, err = tx.Deposit.UpdateOne(dbDeposit).SetStatus(deposit.StatusCompleted).Save(ctx)
	if err != nil {
		s.logger.Error("failed to update deposit", zap.Error(err))
		return utils.Rollback(tx, errors.New("failed to update deposit"))
	}

	_, err = tx.User.UpdateOne(dbDeposit.Edges.User).SetBalance(dbDeposit.Edges.User.Balance + dbDeposit.Amount).Save(ctx)
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
