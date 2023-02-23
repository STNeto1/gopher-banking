package deposit

import "context"

type DepositQueue interface {
	AddMessageToQueue(context.Context, []byte) error
}

type MemoryDepositQueue struct {
}

func (mq MemoryDepositQueue) AddMessageToQueue(_ context.Context, _ []byte) error {
	return nil
}
