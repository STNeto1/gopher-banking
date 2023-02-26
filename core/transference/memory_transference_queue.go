package transference

import "context"

type MemoryTransferenceQueue struct {
}

func (mq MemoryTransferenceQueue) AddTransferenceToQueue(_ context.Context, _ []byte) error {
	return nil
}
