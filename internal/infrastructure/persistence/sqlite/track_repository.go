package sqlite

import (
        "context"
        "database/sql"
        "errors"

        "changeme/internal/db"
        "changeme/internal/domain/models"
        "changeme/internal/domain/repositories"
)

type TrackRepository struct {
        db      *sql.DB
        queries *db.Queries
}

func NewTrackRepository(dbConn *sql.DB, queries *db.Queries) repositories.TrackRepository {
        return &TrackRepository{db: dbConn, queries: queries}
}

func (r *TrackRepository) Create(ctx context.Context, track *models.Track) (*models.Track, error) {
        t, err := r.queries.CreateTrack(ctx, db.CreateTrackParams{
                Title:       track.Title,
                AlbumID:     track.AlbumID,
                ArtistID:    track.ArtistID,
                TrackNumber: track.TrackNumber,
                Duration:    track.Duration,
                FilePath:    track.FilePath,
        })
        if err != nil {
                return nil, err
        }

        result := &models.Track{
                ID:          t.ID,
                Title:       t.Title,
                AlbumID:     t.AlbumID,
                ArtistID:    t.ArtistID,
                TrackNumber: t.TrackNumber,
                Duration:    t.Duration,
                FilePath:    t.FilePath,
        }
        return result, nil
}

func (r *TrackRepository) CreateBatch(ctx context.Context, tracks []*models.Track) error {
        tx, err := r.db.BeginTx(ctx, nil)
        if err != nil {
                return err
        }
        defer tx.Rollback()

        qtx := r.queries.WithTx(tx)
        for _, track := range tracks {
                _, err := qtx.CreateTrack(ctx, db.CreateTrackParams{
                        Title:       track.Title,
                        AlbumID:     track.AlbumID,
                        ArtistID:    track.ArtistID,
                        TrackNumber: track.TrackNumber,
                        Duration:    track.Duration,
                        FilePath:    track.FilePath,
                })
                if err != nil {
                        return err
                }
        }
        return tx.Commit()
}

func (r *TrackRepository) GetByID(ctx context.Context, id int64) (*models.Track, error) {
        return nil, errors.New("not implemented")
}

func (r *TrackRepository) GetAll(ctx context.Context) ([]*models.Track, error) {
        tracks, err := r.queries.ListTracks(ctx)
        if err != nil {
                return nil, err
        }

        result := make([]*models.Track, 0, len(tracks))
        for _, t := range tracks {
                result = append(result, &models.Track{
                        ID:          t.ID,
                        Title:       t.Title,
                        AlbumID:     t.AlbumID,
                        ArtistID:    t.ArtistID,
                        TrackNumber: t.TrackNumber,
                        Duration:    t.Duration,
                        FilePath:    t.FilePath,
                        ArtistName:  t.ArtistName,
                        AlbumTitle:  t.AlbumTitle,
                })
        }
        return result, nil
}

func (r *TrackRepository) Update(ctx context.Context, track *models.Track) (*models.Track, error) {
        return nil, errors.New("not implemented")
}

func (r *TrackRepository) Delete(ctx context.Context, id int64) error {
        return errors.New("not implemented")
}
