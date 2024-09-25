package service

import (
	"context"
	"music-service/internal/domain"
	"music-service/internal/repository"
	"music-service/pkg/logger"

	"log/slog"
)

type SongService interface {
	GetSongs(ctx context.Context, filter repository.SongFilter, limit, offset int) ([]domain.Song, error)
	GetSongLyricsPaginated(ctx context.Context, songID int, limit, offset int) ([]string, error)
	DeleteSong(ctx context.Context, songID int) error
	UpdateSong(ctx context.Context, song domain.Song) error
	AddSong(ctx context.Context, song domain.Song) error
}

type songService struct {
	repo   repository.SongRepository
	logger *logger.Loggers
}

func NewSongService(repo repository.SongRepository, logger *logger.Loggers) SongService {
	return &songService{
		repo:   repo,
		logger: logger,
	}
}

func (s *songService) GetSongs(ctx context.Context, filter repository.SongFilter, limit, offset int) ([]domain.Song, error) {
	s.logger.DebugLogger.Debug("Entering GetSongs service", slog.Any("filter", filter), slog.Int("limit", limit), slog.Int("offset", offset))

	songs, err := s.repo.GetSongs(ctx, filter, limit, offset)
	if err != nil {
		s.logger.ErrorLogger.Error("Error fetching songs", slog.Any("error", err))
		return nil, err
	}

	s.logger.InfoLogger.Info("Successfully fetched songs", slog.Int("count", len(songs)))
	return songs, nil
}

func (s *songService) GetSongLyricsPaginated(ctx context.Context, songID int, limit, offset int) ([]string, error) {
	s.logger.DebugLogger.Debug("Entering GetSongLyricsPaginated service", slog.Int("songID", songID), slog.Int("limit", limit), slog.Int("offset", offset))

	lyrics, err := s.repo.GetSongLyricsPaginated(ctx, songID, limit, offset)
	if err != nil {
		s.logger.ErrorLogger.Error("Error fetching lyrics", slog.Int("songID", songID), slog.Any("error", err))
		return nil, err
	}

	s.logger.InfoLogger.Info("Successfully fetched lyrics", slog.Int("songID", songID), slog.Int("versesCount", len(lyrics)))
	return lyrics, nil
}

func (s *songService) DeleteSong(ctx context.Context, songID int) error {
	s.logger.DebugLogger.Debug("Entering DeleteSong service", slog.Int("songID", songID))

	err := s.repo.DeleteSong(ctx, songID)
	if err != nil {
		s.logger.ErrorLogger.Error("Error deleting song", slog.Int("songID", songID), slog.Any("error", err))
		return err
	}

	s.logger.InfoLogger.Info("Successfully deleted song", slog.Int("songID", songID))
	return nil
}

func (s *songService) UpdateSong(ctx context.Context, song domain.Song) error {
	s.logger.DebugLogger.Debug("Entering UpdateSong service", slog.Any("song", song))

	err := s.repo.UpdateSong(ctx, song)
	if err != nil {
		s.logger.ErrorLogger.Error("Error updating song", slog.Int("songID", song.ID), slog.Any("error", err))
		return err
	}

	s.logger.InfoLogger.Info("Successfully updated song", slog.Int("songID", song.ID))
	return nil
}

func (s *songService) AddSong(ctx context.Context, song domain.Song) error {
	s.logger.DebugLogger.Debug("Entering AddSong service", slog.Any("song", song))

	err := s.repo.AddSong(ctx, song)
	if err != nil {
		s.logger.ErrorLogger.Error("Failed to store the song in the database", slog.Any("error", err))
		return err
	}

	s.logger.InfoLogger.Info("Successfully added song", slog.Any("song", song))
	return nil
}
