package services

import (
	"context"
	"fmt"
	"math/rand"

	"changeme/internal/domain/models"
	"changeme/internal/domain/repositories"
)

// LibraryService orchestrates business logic involving tracks, albums, and artists.
// Notice how it doesn't know about Wails or SQLite directly.
type LibraryService struct {
	trackRepo  repositories.TrackRepository
	albumRepo  repositories.AlbumRepository
	artistRepo repositories.ArtistRepository
}

func NewLibraryService(
	trackRepo repositories.TrackRepository,
	albumRepo repositories.AlbumRepository,
	artistRepo repositories.ArtistRepository,
) *LibraryService {
	return &LibraryService{
		trackRepo:  trackRepo,
		albumRepo:  albumRepo,
		artistRepo: artistRepo,
	}
}

func (s *LibraryService) ListAllTracks(ctx context.Context) ([]*models.Track, error) {
	return s.trackRepo.GetAll(ctx)
}

func (s *LibraryService) ListAllAlbums(ctx context.Context) ([]*models.Album, error) {
	return s.albumRepo.GetAll(ctx)
}

func (s *LibraryService) ListAllArtists(ctx context.Context) ([]*models.Artist, error) {
	return s.artistRepo.GetAll(ctx)
}

// System Benchmarking & Dev Tools
// -------------------------------

// GenerateDummyData generates large datasets to benchmark framework scaling
func (s *LibraryService) GenerateDummyData(ctx context.Context, count int) error {
	artist, err := s.artistRepo.Create(ctx, &models.Artist{Name: "The Benchmarkers"})
	if err != nil {
		return err
	}

	year := int64(2026)
	album, err := s.albumRepo.Create(ctx, &models.Album{
		Title:    "Volume Scale",
		ArtistID: &artist.ID,
		Year:     &year,
	})
	if err != nil {
		return err
	}

	tracks := make([]*models.Track, 0, count)
	for i := 0; i < count; i++ {
		duration := int64(rand.Intn(300-120) + 120)
		trackNumber := int64(i + 1)
		tracks = append(tracks, &models.Track{
			Title:       fmt.Sprintf("Benchmark Track #%d", i+1),
			AlbumID:     &album.ID,
			ArtistID:    &artist.ID,
			TrackNumber: &trackNumber,
			Duration:    &duration,
			FilePath:    fmt.Sprintf("file:///benchmark_mock_%f.mp3", rand.Float64()),
		})
	}
	
	if err := s.trackRepo.CreateBatch(ctx, tracks); err != nil {
		return err
	}

	return nil
}
