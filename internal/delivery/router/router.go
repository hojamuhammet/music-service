package router

import (
	"music-service/internal/delivery/handler"

	_ "music-service/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(songHandler *handler.SongHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/songs", func(r chi.Router) {
		r.Get("/", songHandler.GetSongs)
		r.Get("/{id}/lyrics", songHandler.GetSongLyricsPaginated)
		r.Delete("/{id}", songHandler.DeleteSong)
		r.Put("/{id}", songHandler.UpdateSong)
		r.Post("/", songHandler.AddSong)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
