package port

import (
	"context"

	"github.com/hugosrc/shortlink/internal/core/domain"
)

type LinkService interface {
	Create(ctx context.Context, url string, userID string) (*domain.Link, error)
	FindByHash(ctx context.Context, hash string) (string, error)
	Delete(ctx context.Context, hash string, userID string) error
	Update(ctx context.Context, hash string, url string, userID string) (*domain.Link, error)
}
