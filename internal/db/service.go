package db

import (
        "context"
        "database/sql"
        "embed"
        "fmt"
        "log"
        "time"

        "github.com/golang-migrate/migrate/v4"
        "github.com/golang-migrate/migrate/v4/database/sqlite"
        "github.com/golang-migrate/migrate/v4/source/iofs"
        _ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var fs embed.FS

type Service struct {
        DB      *sql.DB
        Queries *Queries
}

func NewService(dbPath string) (*Service, error) {
        db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
        if err != nil {
                return nil, fmt.Errorf("failed to open database: %w", err)
        }

        db.SetMaxOpenConns(4)
        db.SetMaxIdleConns(2)
        db.SetConnMaxLifetime(time.Hour)

        pingCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()

        if err := db.PingContext(pingCtx); err != nil {
                return nil, fmt.Errorf("failed to ping database: %w", err)
        }

        d, err := iofs.New(fs, "migrations")
        if err != nil {
                return nil, fmt.Errorf("failed to load embedded migrations: %w", err)
        }

        driver, err := sqlite.WithInstance(db, &sqlite.Config{})
        if err != nil {
                return nil, fmt.Errorf("failed to wrap db for migrations: %w", err)
        }

        m, err := migrate.NewWithInstance("iofs", d, "sqlite", driver)
        if err != nil {
                return nil, fmt.Errorf("failed to instantiate migrator: %w", err)
        }

        if err := m.Up(); err != nil && err != migrate.ErrNoChange {
                return nil, fmt.Errorf("failed to run migrations: %w", err)
        }

        log.Println("Database initialized and migrated successfully")
        return &Service{
                DB:      db,
                Queries: New(db),
        }, nil
}

func (s *Service) Close() error {
        return s.DB.Close()
}

// ResetDatabase deletes all data from all tracking tables and resets auto-incrementing ID sequences
func (s *Service) ResetDatabase() error {
        tx, err := s.DB.Begin()
        if err != nil {
                return err
        }
        defer tx.Rollback()

        statements := []string{
                "DELETE FROM playlist_tracks;",
                "DELETE FROM playlists;",
                "DELETE FROM tracks;",
                "DELETE FROM albums;",
                "DELETE FROM artists;",
                "DELETE FROM sqlite_sequence;", // Resets AUTOINCREMENT counters back to 0
        }

        for _, stmt := range statements {
                if _, err := tx.Exec(stmt); err != nil {
                        return fmt.Errorf("failed to execute %q: %w", stmt, err)
                }
        }

        log.Println("Database completely reset")
        return tx.Commit()
}
