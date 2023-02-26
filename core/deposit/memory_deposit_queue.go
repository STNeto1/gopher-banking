package deposit

import "context"

type MemoryDepositQueue struct {
}

func (mq MemoryDepositQueue) AddMessageToQueue(_ context.Context, _ []byte) error {
	return nil
}
