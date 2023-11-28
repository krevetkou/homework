package storage

import "arch-demo/layers_movies/internal/domain"

type MoviesStorage struct {
	actors []domain.Movies
}

func NewMovieStorage() *MoviesStorage {
	return &MoviesStorage{
		actors: make([]domain.Movies, 0),
	}
}
