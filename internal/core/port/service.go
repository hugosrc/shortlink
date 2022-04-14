package port

import (
	"context"

	"github.com/hugosrc/shortlink/internal/core/domain"
)

type LinkService interface {
	Create(ctx context.Context, url string, userID string) (*domain.Link, error)
	FindByHash(ctx context.Context, hash string) (string, error)
}
