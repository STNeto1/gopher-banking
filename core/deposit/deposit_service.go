package deposit

import (
	"models/ent"

	"go.uber.org/zap"
)

type DepositService struct {
	client *ent.Client
	logger *zap.Logger
	queue  DepositQueue
}

func NewDepositService(client *ent.Client, logger *zap.Logger, queue DepositQueue) *DepositService {
	return &DepositService{client: client, logger: logger, queue: queue}
}
