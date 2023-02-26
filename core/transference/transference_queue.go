package transference

import "context"

type TransferenceQueue interface {
	AddTransferenceToQueue(context.Context, []byte) error
}
