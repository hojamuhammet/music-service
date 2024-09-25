package repository

import (
	"context"
	"database/sql"
	"music-service/internal/domain"
	"music-service/pkg/logger"
	"strconv"
	"strings"

	"log/slog"
)

type SongRepository interface {
	GetSongs(ctx context.Context, filter SongFilter, limit, offset int) ([]domain.Song, error)
	GetSongLyricsPaginated(ctx context.Context, songID int, limit, offset int) ([]string, error)
	DeleteSong(ctx context.Context, songID int) error
	UpdateSong(ctx context.Context, song domain.Song) error
	AddSong(ctx context.Context, song domain.Song) error
}

type SongFilter struct {
	Group       string
	Song        string
	ReleaseDate string
}

type songRepository struct {
	db     *sql.DB
	logger *logger.Loggers
}

func NewSongRepository(db *sql.DB, logger *logger.Loggers) SongRepository {
	return &songRepository{db: db, logger: logger}
}

func (r *songRepository) GetSongs(ctx context.Context, filter SongFilter, limit, offset int) ([]domain.Song, error) {
	r.logger.DebugLogger.Debug("Entering GetSongs", slog.Any("filter", filter))

	query := "SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE 1=1"
	var args []interface{}
	argIndex := 1

	if filter.Group != "" {
		query += " AND group_name ILIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+filter.Group+"%")
		argIndex++
	}

	if filter.Song != "" {
		query += " AND song_name ILIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+filter.Song+"%")
		argIndex++
	}

	if filter.ReleaseDate != "" {
		query += " AND release_date = $" + strconv.Itoa(argIndex)
		args = append(args, filter.ReleaseDate)
		argIndex++
	}

	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	r.logger.DebugLogger.Debug("Executing query", slog.String("query", query), slog.Any("args", args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.ErrorLogger.Error("Error executing GetSongs query", slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()

	var songs []domain.Song
	for rows.Next() {
		var song domain.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			r.logger.ErrorLogger.Error("Error scanning song row", slog.Any("error", err))
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		r.logger.ErrorLogger.Error("Error iterating over song rows", slog.Any("error", err))
		return nil, err
	}

	r.logger.InfoLogger.Info("Successfully fetched songs", slog.Int("count", len(songs)))
	return songs, nil
}

func (r *songRepository) GetSongLyricsPaginated(ctx context.Context, songID int, limit, offset int) ([]string, error) {
	r.logger.DebugLogger.Debug("Entering GetSongLyricsPaginated", slog.Int("songID", songID))

	query := "SELECT text FROM songs WHERE id = $1"
	r.logger.DebugLogger.Debug("Executing query", slog.String("query", query), slog.Int("songID", songID))

	row := r.db.QueryRowContext(ctx, query, songID)

	var lyrics string
	if err := row.Scan(&lyrics); err != nil {
		r.logger.ErrorLogger.Error("Error fetching song lyrics", slog.Any("error", err))
		return nil, err
	}

	verses := strings.Split(lyrics, "\n")
	start := offset
	end := offset + limit

	if start >= len(verses) {
		return nil, nil // no more verses to return
	}

	if end > len(verses) {
		end = len(verses) // don’t go out of bounds
	}

	r.logger.InfoLogger.Info("Successfully fetched lyrics for song", slog.Int("songID", songID), slog.Any("verses", verses[start:end]))
	return verses[start:end], nil
}

func (r *songRepository) DeleteSong(ctx context.Context, songID int) error {
	r.logger.DebugLogger.Debug("Entering DeleteSong", slog.Int("songID", songID))

	query := "DELETE FROM songs WHERE id = $1"
	r.logger.DebugLogger.Debug("Executing query", slog.String("query", query), slog.Int("songID", songID))

	_, err := r.db.ExecContext(ctx, query, songID)
	if err != nil {
		r.logger.ErrorLogger.Error("Error deleting song", slog.Int("songID", songID), slog.Any("error", err))
		return err
	}

	r.logger.InfoLogger.Info("Successfully deleted song", slog.Int("songID", songID))
	return nil
}

func (r *songRepository) UpdateSong(ctx context.Context, song domain.Song) error {
	r.logger.DebugLogger.Debug("Entering UpdateSong", slog.Any("song", song))

	query := `
		UPDATE songs
		SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5
		WHERE id = $6
	`
	r.logger.DebugLogger.Debug("Executing query", slog.String("query", query), slog.Any("song", song))

	_, err := r.db.ExecContext(ctx, query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID)
	if err != nil {
		r.logger.ErrorLogger.Error("Error updating song", slog.Int("songID", song.ID), slog.Any("error", err))
		return err
	}

	r.logger.InfoLogger.Info("Successfully updated song", slog.Int("songID", song.ID))
	return nil
}

func (r *songRepository) AddSong(ctx context.Context, song domain.Song) error {
	r.logger.DebugLogger.Debug("Entering AddSong", slog.Any("song", song))

	query := `
		INSERT INTO songs (group_name, song_name, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
	`
	r.logger.DebugLogger.Debug("Executing query", slog.String("query", query), slog.Any("song", song))

	_, err := r.db.ExecContext(ctx, query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		r.logger.ErrorLogger.Error("Error adding song", slog.Any("error", err))
		return err
	}

	r.logger.InfoLogger.Info("Successfully added song", slog.Any("song", song))
	return nil
}
