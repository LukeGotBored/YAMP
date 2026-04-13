package models

import "time"

type Track struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	AlbumID     *int64      `json:"album_id"`
	ArtistID    *int64      `json:"artist_id"`
	TrackNumber *int64      `json:"track_number"`
	Duration    *int64      `json:"duration"`
	FilePath    string      `json:"file_path"`
	CreatedAt   *time.Time  `json:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at"`

	// DTO associations
	ArtistName  string      `json:"artist_name"`
	AlbumTitle  string      `json:"album_title"`
}
