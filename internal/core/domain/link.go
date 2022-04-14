package domain

import "time"

// Link is the structural representation of the application domain
type Link struct {
	Hash         string    `json:"hash"`
	OriginalURL  string    `json:"original_url"`
	UserID       string    `json:"user_id"`
	CreationTime time.Time `json:"creation_time"`
}
