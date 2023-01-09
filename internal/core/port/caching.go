package port

import "context"

type LinkCaching interface {
	Get(ctx context.Context, hash string) (string, error)
	Set(ctx context.Context, hash string, originalURL string) error
	Del(ctx context.Context, hash string) error
}
