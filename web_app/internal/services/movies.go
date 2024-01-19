package services

import (
	"arch-demo/internal/domain"
	"errors"
	"fmt"
	"strings"
)

type MoviesRepository interface {
	InsertMovie(actor domain.Movie) (domain.Movie, error)
	IsMovieExists(actor domain.Movie) (bool, error)
	GetMovieByID(id int) (domain.Movie, error)
	UpdateMovie(actor domain.Movie) error
	DeleteMovie(id int) error
	GetAllMovies() ([]domain.Movie, error)
	SortAndOrderByMovie(sortBy, orderBy string, movies []domain.Movie) []domain.Movie
	GetActorsByMovie(id int) ([]domain.Actor, error)
	CreateActorsByMovie(id int, actors []int) (int, []int, error)
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
	if movie.Name == "" || movie.ReleaseDate.String() == "" ||
		movie.Country == "" || movie.Genre == "" || movie.Rating == 0 {
		return domain.Movie{}, domain.ErrFieldsRequired
	}

	newMovie, err := s.Storage.InsertMovie(movie)
	if err != nil {
		return domain.Movie{}, domain.ErrExists
	}

	return newMovie, nil
}

func (s MoviesService) Get(id int) (domain.Movie, error) {
	movie, err := s.Storage.GetMovieByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.Movie{}, fmt.Errorf("movie id: %d, err: %w", id, err)
	}

	if err != nil {
		return domain.Movie{}, fmt.Errorf("failed to find movie, unexpected error: %w", err)
	}

	return movie, nil
}

func (s MoviesService) Update(id int, movieUpdate domain.MovieUpdate) (domain.Movie, error) {
	movie, err := s.Storage.GetMovieByID(id)
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

	s.Storage.UpdateMovie(movie)

	return movie, nil
}

func (s MoviesService) Delete(id int) error {
	_, err := s.Storage.GetMovieByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("movie id: %d, err: %w", id, err)
	}

	if err != nil {
		return fmt.Errorf("failed to find movie, unexpected error: %w", err)
	}

	s.Storage.DeleteMovie(id)

	return nil
}

func (s MoviesService) List(orderBy, sortBy, nameQuery, genreQuery string) []domain.Movie {
	movies, err := s.Storage.GetAllMovies()
	if err != nil {
		return []domain.Movie{}
	}
	var filteredMovies []domain.Movie

	if nameQuery != "" || genreQuery != "" {
		for i := range movies {
			if (nameQuery != "" && strings.Contains(movies[i].Name, nameQuery)) ||
				(genreQuery != "" && strings.Contains(movies[i].Genre, genreQuery)) {
				filteredMovies = append(filteredMovies, movies[i])
			}
		}

	} else {
		filteredMovies = movies
	}

	switch {
	case sortBy == "" && orderBy == "":
		movies = s.Storage.SortAndOrderByMovie("name", "asc", filteredMovies)
	case sortBy == "" && orderBy != "":
		movies = s.Storage.SortAndOrderByMovie("name", orderBy, filteredMovies)
	case sortBy != "" && orderBy == "":
		movies = s.Storage.SortAndOrderByMovie(sortBy, "asc", filteredMovies)
	default:
		movies = s.Storage.SortAndOrderByMovie(sortBy, orderBy, filteredMovies)
	}

	return movies
}

func (s MoviesService) GetActorsByMovie(id int) ([]domain.Actor, error) {
	actors, err := s.Storage.GetActorsByMovie(id) //add error
	if errors.Is(err, domain.ErrNotFound) {
		return []domain.Actor{}, err
	}

	if err != nil {
		return []domain.Actor{}, fmt.Errorf("failed to find actors, unexpected error: %w", err)
	}

	return actors, nil
}

func (s MoviesService) CreateActorsForMovie(id int, actorsByMovie []int) (int, []int, error) {
	var movieID int
	var actorsIDs []int
	movieID, actorsIDs, err := s.Storage.CreateActorsByMovie(id, actorsByMovie)

	if err != nil {
		return 0, []int{0}, err
	}

	return movieID, actorsIDs, nil
}
