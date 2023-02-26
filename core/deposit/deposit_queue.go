package deposit

import "context"

type DepositQueue interface {
	AddMessageToQueue(context.Context, []byte) error
}
