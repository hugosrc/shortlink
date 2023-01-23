package service

import (
	"context"
	"strconv"
	"time"

	"github.com/hugosrc/shortlink/internal/core/domain"
	"github.com/hugosrc/shortlink/internal/core/port"
	"github.com/hugosrc/shortlink/internal/util"
)

type LinkService struct {
	counter port.Counter
	encoder port.Encoder
	caching port.LinkCaching
	repo    port.LinkRepository
}

func NewLinkService(counter port.Counter, encoder port.Encoder, caching port.LinkCaching, repo port.LinkRepository) port.LinkService {
	return &LinkService{
		counter: counter,
		encoder: encoder,
		caching: caching,
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
		Hash:         hash[0:7],
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
	url, _ := s.caching.Get(ctx, hash)

	if len(url) > 0 {
		return url, nil
	}

	link, err := s.repo.FindByHash(ctx, hash)
	if err != nil {
		return "", err
	}

	_ = s.caching.Set(ctx, hash, link.OriginalURL)

	return link.OriginalURL, nil
}

func (s *LinkService) Delete(ctx context.Context, hash string, userID string) error {
	link, err := s.repo.FindByHash(ctx, hash)
	if err != nil {
		return err
	}

	if link.UserID != userID {
		return util.NewErrorf(util.ErrCodeUnauthorized, "user does not have permission")
	}

	if err := s.repo.Delete(ctx, hash); err != nil {
		return err
	}

	_ = s.caching.Del(ctx, hash)

	return nil
}

func (s *LinkService) Update(ctx context.Context, hash string, newURL string, userID string) (*domain.Link, error) {
	link, err := s.repo.FindByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	if link.UserID != userID {
		return nil, util.NewErrorf(util.ErrCodeUnauthorized, "user does not have permission")
	}

	if err := s.repo.Update(ctx, hash, newURL); err != nil {
		return nil, err
	}

	_ = s.caching.Set(ctx, hash, newURL)

	return &domain.Link{
		Hash:         hash,
		OriginalURL:  newURL,
		UserID:       userID,
		CreationTime: link.CreationTime,
	}, nil
}
