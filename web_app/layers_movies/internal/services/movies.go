package services

import (
	"arch-demo/layers_movies/internal/domain"
	"errors"
	"fmt"
	"strings"
)

type MoviesRepository interface {
	Insert(actor domain.Movie) domain.Movie
	IsMovieExists(actor domain.Movie) bool
	GetByID(id int) (domain.Movie, error)
	Update(actor domain.Movie)
	Delete(id int)
	GetAll() []domain.Movie
	SortAndOrderBy(sortBy, orderBy string) []domain.Movie
}

type MoviesService struct {
	Storage MoviesRepository
}

func NewMovieService(storage MoviesRepository) MoviesService {
	return MoviesService{
		Storage: storage,
	}
}

func (s MoviesService) Create(movie domain.Movie) (domain.Movie, error) {
	// входящие параметры необходимо валидировать
	if movie.Name == "" || movie.ReleaseDate.Date == 0 || movie.ReleaseDate.Month == 0 || movie.ReleaseDate.Year == 0 ||
		movie.Country == "" || movie.Genre == "" || movie.Rating == 0 {
		return domain.Movie{}, domain.ErrFieldsRequired
	}

	isMovieExist := s.Storage.IsMovieExists(movie)
	if isMovieExist {
		return domain.Movie{}, domain.ErrMovieExists
	}

	newMovie := s.Storage.Insert(movie)

	return newMovie, nil
}

func (s MoviesService) Get(id int) (domain.Movie, error) {
	movie, err := s.Storage.GetByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.Movie{}, err
	}

	if err != nil {
		return domain.Movie{}, fmt.Errorf("failed to find movie, unexpected error: %w", err)
	}

	return movie, nil
}

func (s MoviesService) Update(id int, movieUpdate domain.MovieUpdate) (domain.Movie, error) {
	movie, err := s.Storage.GetByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.Movie{}, err
	}

	if movieUpdate.Name != nil {
		movie.Name = *movieUpdate.Name
	}

	if movieUpdate.ReleaseDate != nil {
		movie.ReleaseDate = *movieUpdate.ReleaseDate
	}

	if movieUpdate.Country != nil {
		movie.Country = *movieUpdate.Country
	}

	if movieUpdate.Genre != nil {
		movie.Genre = *movieUpdate.Genre
	}

	if movieUpdate.Rating != nil {
		movie.Rating = *movieUpdate.Rating
	}

	s.Storage.Update(movie)

	return movie, nil
}

func (s MoviesService) Delete(id int) error {
	_, err := s.Storage.GetByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("movie id: %d, err: %w", id, err)
	}

	if err != nil {
		return fmt.Errorf("failed to find movie, unexpected error: %w", err)
	}

	s.Storage.Delete(id)

	return nil
}

func (s MoviesService) List(orderBy, sortBy, nameQuery, genreQuery string) []domain.Movie {
	movies := s.Storage.GetAll()
	var filteredMovies []domain.Movie

	if nameQuery != "" || genreQuery != "" {
		for i := range movies {
			if (strings.Contains(movies[i].Name, nameQuery)) ||
				(strings.Contains(movies[i].Genre, genreQuery)) {
				filteredMovies = append(filteredMovies, movies[i])
			}
		}
		return filteredMovies
	}

	switch {
	case sortBy == "" && orderBy == "":
		movies = s.Storage.SortAndOrderBy("name", "asc")
	case sortBy == "" && orderBy != "":
		movies = s.Storage.SortAndOrderBy("name", orderBy)
	case sortBy != "" && orderBy == "":
		movies = s.Storage.SortAndOrderBy(sortBy, "asc")
	default:
		movies = s.Storage.SortAndOrderBy(sortBy, orderBy)
	}

	return movies
}
