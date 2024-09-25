package handler

import (
	"encoding/json"
	"log/slog"
	"music-service/internal/domain"
	"music-service/internal/repository"
	"music-service/internal/service"
	"music-service/pkg/logger"
	"music-service/pkg/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SongHandler struct {
	songService service.SongService
	loggers     *logger.Loggers
}

func NewSongHandler(songService service.SongService, loggers *logger.Loggers) *SongHandler {
	return &SongHandler{songService: songService, loggers: loggers}
}

// GetSongs godoc
// @Summary Get songs with optional filtering and pagination
// @Description Retrieve songs filtered by group, song name, and/or release date with pagination.
// @Tags songs
// @Accept json
// @Produce json
// @Param group_name query string false "Filter by group name"
// @Param song_name query string false "Filter by song name"
// @Param release_date query string false "Filter by release date"
// @Param limit query int false "Pagination limit"
// @Param offset query int false "Pagination offset"
// @Success 200 {array} domain.Song
// @Failure 500 {object} utils.JSONError "Failed to fetch songs"
// @Router /songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.loggers.DebugLogger.Debug("Handling GetSongs request")

	groupName := r.URL.Query().Get("group_name")
	songName := r.URL.Query().Get("song_name")
	releaseDate := r.URL.Query().Get("release_date")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	filter := repository.SongFilter{
		Group:       groupName,
		Song:        songName,
		ReleaseDate: releaseDate,
	}

	songs, err := h.songService.GetSongs(ctx, filter, limit, offset)
	if err != nil {
		h.loggers.ErrorLogger.Error("Failed to fetch songs", utils.Err(err))
		utils.RespondWithErrorJSON(w, http.StatusInternalServerError, "Failed to fetch songs")
		return
	}

	h.loggers.InfoLogger.Info("Fetched songs successfully", slog.Int("count", len(songs)))
	utils.RespondWithJSON(w, http.StatusOK, songs)
}

// GetSongLyricsPaginated godoc
// @Summary Get paginated lyrics of a song
// @Description Get song lyrics paginated by verses.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param limit query int false "Pagination limit"
// @Param offset query int false "Pagination offset"
// @Success 200 {array} string
// @Failure 400 {object} utils.JSONError "Invalid song ID"
// @Failure 500 {object} utils.JSONError "Failed to fetch lyrics"
// @Router /songs/{id}/lyrics [get]
func (h *SongHandler) GetSongLyricsPaginated(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.loggers.DebugLogger.Debug("Handling GetSongLyricsPaginated request")

	songID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.loggers.ErrorLogger.Error("Invalid song ID", utils.Err(err))
		utils.RespondWithErrorJSON(w, http.StatusBadRequest, "Invalid song ID")
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	lyrics, err := h.songService.GetSongLyricsPaginated(ctx, songID, limit, offset)
	if err != nil {
		h.loggers.ErrorLogger.Error("Failed to fetch lyrics", utils.Err(err))
		utils.RespondWithErrorJSON(w, http.StatusInternalServerError, "Failed to fetch lyrics")
		return
	}

	h.loggers.InfoLogger.Info("Fetched song lyrics successfully", slog.Int("songID", songID))
	utils.RespondWithJSON(w, http.StatusOK, lyrics)
}

// DeleteSong godoc
// @Summary Delete a song by ID
// @Description Delete a song from the library and return a status and message.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]string "status and message"
// @Failure 400 {object} utils.JSONError "Invalid song ID"
// @Failure 500 {object} utils.JSONError "Failed to delete song"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.loggers.DebugLogger.Debug("Handling DeleteSong request")

	songID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.loggers.ErrorLogger.Error("Invalid song ID", utils.Err(err))
		utils.RespondWithErrorJSON(w, http.StatusBadRequest, "Invalid song ID")
		return
	}

	if err := h.songService.DeleteSong(ctx, songID); err != nil {
		h.loggers.ErrorLogger.Error("Failed to delete song", utils.Err(err))
		utils.RespondWithErrorJSON(w, http.StatusInternalServerError, "Failed to delete song")
		return
	}

	h.loggers.InfoLogger.Info("Deleted song successfully", slog.Int("songID", songID))
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Song deleted successfully",
	})
}

// UpdateSong godoc
// @Summary Update a song's details
// @Description Update an existing song's details in the library and return the updated song data.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body domain.Song true "Updated song"
// @Success 200 {object} domain.Song "Updated song details"
// @Failure 400 {object} utils.JSONError "Invalid song ID or payload"
// @Failure 500 {object} utils.JSONError "Failed to update song"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.loggers.DebugLogger.Debug("Handling UpdateSong request")

	songID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.loggers.ErrorLogger.Error("Invalid song ID", utils.Err(err))
		utils.RespondWithErrorJSON(w, http.StatusBadRequest, "Invalid song ID")
		return
	}

	var song domain.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		h.loggers.ErrorLogger.Error("Invalid request payload", utils.Err(err))
		utils.RespondWithErrorJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	song.ID = songID

	if err := h.songService.UpdateSong(ctx, song); err != nil {
		h.loggers.ErrorLogger.Error("Failed to update song", utils.Err(err))
		utils.RespondWithErrorJSON(w, http.StatusInternalServerError, "Failed to update song")
		return
	}

	updatedSong, err := h.songService.GetSongByID(ctx, songID)
	if err != nil {
		h.loggers.ErrorLogger.Error("Failed to retrieve updated song", utils.Err(err))
		utils.RespondWithErrorJSON(w, http.StatusInternalServerError, "Failed to retrieve updated song")
		return
	}

	// Return the updated song as response
	h.loggers.InfoLogger.Info("Updated song successfully", slog.Int("songID", songID))
	utils.RespondWithJSON(w, http.StatusOK, updatedSong)
}

// AddSong godoc
// @Summary Add a new song
// @Description Adds a new song to the library.
// @Tags songs
// @Accept json
// @Produce json
// @Param song body domain.Song true "New song to add"
// @Success 201 {object} map[string]interface{} "status and message"
// @Failure 400 {object} utils.JSONError "Invalid request payload"
// @Failure 500 {object} utils.JSONError "Failed to add song"
// @Router /songs [post]
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.loggers.DebugLogger.Debug("Handling AddSong request")

	var song domain.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		h.loggers.ErrorLogger.Error("Invalid request payload", slog.Any("error", err))
		utils.RespondWithErrorJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if song.Group == "" || song.Song == "" {
		h.loggers.ErrorLogger.Error("Group and song fields are required")
		utils.RespondWithErrorJSON(w, http.StatusBadRequest, "Group and song fields are required")
		return
	}

	if err := h.songService.AddSong(ctx, song); err != nil {
		h.loggers.ErrorLogger.Error("Failed to add song", slog.Any("error", err))
		utils.RespondWithErrorJSON(w, http.StatusInternalServerError, "Failed to add song")
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Song added successfully",
	}
	h.loggers.InfoLogger.Info("Added song successfully", slog.String("group", song.Group), slog.String("song", song.Song))
	utils.RespondWithJSON(w, http.StatusCreated, response)
}
