package service

import (
	"context"
	"errors"
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

func (s *LinkService) Delete(ctx context.Context, hash string, userID string) error {
	link, err := s.repo.FindByHash(ctx, hash)
	if err != nil {
		return err
	}

	if link.UserID != userID {
		return errors.New("user does not have permission")
	}

	if err := s.repo.Delete(ctx, hash); err != nil {
		return err
	}

	return nil
}

func (s *LinkService) Update(ctx context.Context, hash string, newURL string, userID string) (*domain.Link, error) {
	link, err := s.repo.FindByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	if link.UserID != userID {
		return nil, errors.New("user does not have permission")
	}

	if err := s.repo.Update(ctx, hash, newURL); err != nil {
		return nil, err
	}

	return &domain.Link{
		Hash:         hash,
		OriginalURL:  newURL,
		UserID:       userID,
		CreationTime: link.CreationTime,
	}, nil
}
