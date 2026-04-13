package repositories

import (
        "context"

        "changeme/internal/domain/models"
)

type TrackRepository interface {
        Create(ctx context.Context, track *models.Track) (*models.Track, error)
        CreateBatch(ctx context.Context, tracks []*models.Track) error
        GetByID(ctx context.Context, id int64) (*models.Track, error)
        GetAll(ctx context.Context) ([]*models.Track, error)
        Update(ctx context.Context, track *models.Track) (*models.Track, error)
        Delete(ctx context.Context, id int64) error
}

type AlbumRepository interface {
        Create(ctx context.Context, album *models.Album) (*models.Album, error)
        GetByID(ctx context.Context, id int64) (*models.Album, error)
        GetAll(ctx context.Context) ([]*models.Album, error)
        Update(ctx context.Context, album *models.Album) (*models.Album, error)
        Delete(ctx context.Context, id int64) error
        GetByArtistID(ctx context.Context, artistID int64) ([]*models.Album, error)
}

type ArtistRepository interface {
        Create(ctx context.Context, artist *models.Artist) (*models.Artist, error)
        GetByID(ctx context.Context, id int64) (*models.Artist, error)
        GetAll(ctx context.Context) ([]*models.Artist, error)
        Update(ctx context.Context, artist *models.Artist) (*models.Artist, error)
        Delete(ctx context.Context, id int64) error
}
