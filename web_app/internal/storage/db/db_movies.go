package db

import (
	"arch-demo/internal/domain"
)

func (s *StorageDB) InsertMovie(actor domain.Movie) domain.Movie {
	return domain.Movie{}
}

func (s *StorageDB) IsMovieExists(movie domain.Movie) bool {
	return false
}

func (s *StorageDB) GetMovieByID(id int) (domain.Movie, error) {
	return domain.Movie{}, nil
}

func (s *StorageDB) UpdateMovie(movieUpdate domain.Movie) {

}

func (s *StorageDB) DeleteMovie(id int) {

}

func (s *StorageDB) GetAllMovies() []domain.Movie {
	return []domain.Movie{}
}

func (s *StorageDB) SortAndOrderByMovie(sortBy, orderBy string, movies []domain.Movie) []domain.Movie {
	return []domain.Movie{}
}

func (s *StorageDB) GetActorsByMovie(id int) ([]domain.Actor, error) {
	return []domain.Actor{}, nil
}

func (s *StorageDB) CreateActorsByMovie(id int, actors []int) error {
	return nil
}
