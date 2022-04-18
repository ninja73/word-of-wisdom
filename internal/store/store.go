package store

import "context"

type Store interface {
	RandomQuote(ctx context.Context) (string, error)
}
