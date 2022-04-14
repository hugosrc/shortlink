package service

import (
	"context"
	"strconv"
	"time"

	"github.com/hugosrc/shortlink/internal/core/domain"
	"github.com/hugosrc/shortlink/internal/core/port"
)

type LinkService struct {
	counter port.Counter
	encoder port.Encoder
	repo    port.LinkRepository
}

func NewLinkService(counter port.Counter, encoder port.Encoder, repo port.LinkRepository) port.LinkService {
	return &LinkService{
		counter: counter,
		encoder: encoder,
		repo:    repo,
	}
}

func (s *LinkService) Create(ctx context.Context, url string, userID string) (*domain.Link, error) {
	c, err := s.counter.Inc()
	if err != nil {
		return nil, err
	}

	hash := s.encoder.EncodeToString([]byte(strconv.Itoa(c)))
	link := &domain.Link{
		Hash:         hash,
		OriginalURL:  url,
		UserID:       userID,
		CreationTime: time.Now(),
	}

	if err := s.repo.Create(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *LinkService) FindByHash(ctx context.Context, hash string) (string, error) {
	link, err := s.repo.FindByHash(ctx, hash)
	if err != nil {
		return "", err
	}

	return link.OriginalURL, nil
}
