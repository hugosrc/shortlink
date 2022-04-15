package repository

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/hugosrc/shortlink/internal/core/domain"
	"github.com/hugosrc/shortlink/internal/core/port"
)

type LinkRepository struct {
	conn *gocql.Session
}

func NewLinkRepository(conn *gocql.Session) port.LinkRepository {
	return &LinkRepository{
		conn: conn,
	}
}

func (r *LinkRepository) Create(ctx context.Context, link *domain.Link) error {
	if err := r.conn.Query(
		"INSERT INTO shorturl.url_mapping (hash, original_url, user_id, creation_time) VALUES (?, ?, ?, ?);",
		link.Hash,
		link.OriginalURL,
		link.UserID,
		link.CreationTime,
	).WithContext(ctx).Exec(); err != nil {
		return err
	}

	return nil
}

func (r *LinkRepository) Delete(ctx context.Context, hash string) error {
	if err := r.conn.Query(
		"DELETE FROM shorturl.url_mapping WHERE hash = ?;",
		hash,
	).WithContext(ctx).Exec(); err != nil {
		return err
	}

	return nil
}

func (r *LinkRepository) FindByHash(ctx context.Context, hash string) (*domain.Link, error) {
	var link domain.Link
	if err := r.conn.Query(
		"SELECT hash, original_url, user_id, creation_time FROM shorturl.url_mapping WHERE hash = ?;", hash,
	).WithContext(ctx).Consistency(gocql.One).Scan(
		&link.Hash,
		&link.OriginalURL,
		&link.UserID,
		&link.CreationTime,
	); err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *LinkRepository) Update(ctx context.Context, hash string, newURL string) error {
	if err := r.conn.Query(
		"UPDATE shorturl.url_mapping SET original_url = ? WHERE hash = ?;",
		newURL,
		hash,
	).WithContext(ctx).Exec(); err != nil {
		return err
	}

	return nil
}
