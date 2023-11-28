package main

import (
	"arch-demo/layers_movies/internal/api"
	"arch-demo/layers_movies/internal/services"
	"arch-demo/layers_movies/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	moviesStorage := storage.NewMovieStorage()
	moviesService := services.NewMovieService(moviesStorage)
	moviesHandler := api.NewLaptopsHandler(moviesService)

	r := chi.NewRouter()
	r.Route("/movies", func(r chi.Router) {
		r.Post("/", moviesHandler.Create)
		r.Get("/", moviesHandler.List)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", moviesHandler.Get)
			r.Patch("/", moviesHandler.Update)
			r.Delete("/", moviesHandler.Delete)
		})
	})
}
