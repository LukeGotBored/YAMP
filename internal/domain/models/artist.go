package models

import "time"

type Artist struct {
	ID        int64
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
