package sqlite

import (
	"context"
	"errors"

	"changeme/internal/db"
	"changeme/internal/domain/models"
	"changeme/internal/domain/repositories"
)

type AlbumRepository struct {
	queries *db.Queries
}

func NewAlbumRepository(queries *db.Queries) repositories.AlbumRepository {
	return &AlbumRepository{queries: queries}
}

func (r *AlbumRepository) Create(ctx context.Context, album *models.Album) (*models.Album, error) {
	a, err := r.queries.CreateAlbum(ctx, db.CreateAlbumParams{
		Title:    album.Title,
		ArtistID: album.ArtistID,
		Year:     album.Year,
	})
	if err != nil {
		return nil, err
	}

	result := &models.Album{
		ID:       a.ID,
		Title:    a.Title,
		ArtistID: a.ArtistID,
		Year:     a.Year,
	}
	return result, nil
}

func (r *AlbumRepository) GetByID(ctx context.Context, id int64) (*models.Album, error) {
	return nil, errors.New("not implemented")
}

func (r *AlbumRepository) GetAll(ctx context.Context) ([]*models.Album, error) {
	albums, err := r.queries.ListAlbums(ctx)
	if err != nil {
		return nil, err
	}

	var result []*models.Album
	for _, a := range albums {
		result = append(result, &models.Album{
			ID:       a.ID,
			Title:    a.Title,
			ArtistID: a.ArtistID,
			Year:     a.Year,
		})
	}
	return result, nil
}

func (r *AlbumRepository) Update(ctx context.Context, album *models.Album) (*models.Album, error) {
	return nil, errors.New("not implemented")
}

func (r *AlbumRepository) Delete(ctx context.Context, id int64) error {
	return errors.New("not implemented")
}

func (r *AlbumRepository) GetByArtistID(ctx context.Context, artistID int64) ([]*models.Album, error) {
	return nil, errors.New("not implemented")
}
