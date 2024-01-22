package storage

import "context"

type Adapter interface {
	Flush(ctx context.Context) error
	Restore(ctx context.Context) error
}
