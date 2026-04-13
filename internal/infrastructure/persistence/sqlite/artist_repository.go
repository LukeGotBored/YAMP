package sqlite

import (
	"context"
	"errors"

	"changeme/internal/db"
	"changeme/internal/domain/models"
	"changeme/internal/domain/repositories"
)

type ArtistRepository struct {
	queries *db.Queries
}

func NewArtistRepository(queries *db.Queries) repositories.ArtistRepository {
	return &ArtistRepository{queries: queries}
}

func (r *ArtistRepository) Create(ctx context.Context, artist *models.Artist) (*models.Artist, error) {
	a, err := r.queries.CreateArtist(ctx, artist.Name)
	if err != nil {
		return nil, err
	}

	result := &models.Artist{
		ID:   a.ID,
		Name: a.Name,
	}
	return result, nil
}

func (r *ArtistRepository) GetByID(ctx context.Context, id int64) (*models.Artist, error) {
	return nil, errors.New("not implemented")
}

func (r *ArtistRepository) GetAll(ctx context.Context) ([]*models.Artist, error) {
	artists, err := r.queries.ListArtists(ctx)
	if err != nil {
		return nil, err
	}

	var result []*models.Artist
	for _, a := range artists {
		result = append(result, &models.Artist{
			ID:   a.ID,
			Name: a.Name,
		})
	}
	return result, nil
}

func (r *ArtistRepository) Update(ctx context.Context, artist *models.Artist) (*models.Artist, error) {
	return nil, errors.New("not implemented")
}

func (r *ArtistRepository) Delete(ctx context.Context, id int64) error {
	return errors.New("not implemented")
}
