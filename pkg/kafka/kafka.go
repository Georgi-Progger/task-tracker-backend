package kafka

import "context"

type Producer interface {
	Send(ctx context.Context, value []byte) error
}
