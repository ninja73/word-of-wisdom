package cache

import "context"

type Cache interface {
	ContainsOrAdd(ctx context.Context, val uint64) (bool, error)
}
