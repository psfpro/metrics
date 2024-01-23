package storage

import "context"

type DumbAdapter struct {
}

func NewDumbAdapter() *DumbAdapter {
	return &DumbAdapter{}
}

func (obj *DumbAdapter) Flush(ctx context.Context) error {
	return nil
}

func (obj *DumbAdapter) Restore(ctx context.Context) error {
	return nil
}
