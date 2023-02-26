package transference

import (
	"models/ent"

	"go.uber.org/zap"
)

type TransferService struct {
	client *ent.Client
	logger *zap.Logger
	queue  TransferenceQueue
}

func NewTransferService(client *ent.Client, logger *zap.Logger, queue TransferenceQueue) *TransferService {
	return &TransferService{client: client, logger: logger, queue: queue}
}
