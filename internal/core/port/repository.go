package port

import (
	"context"

	"github.com/hugosrc/shortlink/internal/core/domain"
)

// LinkRepository is an abstraction for accessing a data storage system.
type LinkRepository interface {
	Create(ctx context.Context, link *domain.Link) error
	FindByHash(ctx context.Context, hash string) (*domain.Link, error)
}
