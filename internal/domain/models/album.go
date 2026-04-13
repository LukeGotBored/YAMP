package models

import "time"

type Album struct {
	ID        int64
	Title     string
	ArtistID  *int64
	Year      *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
