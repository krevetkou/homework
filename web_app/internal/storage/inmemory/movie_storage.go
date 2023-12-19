package inmemory

import (
	"arch-demo/internal/domain"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

func (s *Storage) InsertMovie(movie domain.Movie) (domain.Movie, error) {
	var lastID int
	if len(s.movies) > 0 {
		lastID = s.movies[len(s.movies)-1:][0].ID
	}

	movie.ID = lastID + 1

	s.movies = append(s.movies, movie)
	return movie, nil
}

func (s *Storage) IsMovieExists(movie domain.Movie) (bool, error) {
	for i := range s.movies {
		if strings.Contains(s.movies[i].Name, movie.Name) &&
			s.movies[i].ReleaseDate == movie.ReleaseDate &&
			strings.Contains(s.movies[i].Country, movie.Country) &&
			strings.Contains(s.movies[i].Genre, movie.Genre) &&
			s.movies[i].Rating == movie.Rating {
			return true, nil
		}
	}

	return false, nil
}

func (s *Storage) GetMovieByID(id int) (domain.Movie, error) {
	var movie *domain.Movie
	for i := range s.movies {
		if s.movies[i].ID == id {
			movie = &s.movies[i]
		}
	}

	if movie == nil {
		return domain.Movie{}, domain.ErrNotFound
	}

	return *movie, nil
}

func (s *Storage) UpdateMovie(movieUpdate domain.Movie) error {
	for i := range s.movies {
		if s.movies[i].ID == movieUpdate.ID {
			s.movies[i] = movieUpdate
		}
	}

	return nil
}

func (s *Storage) DeleteMovie(id int) error {
	s.movies = slices.DeleteFunc(s.movies, func(l1 domain.Movie) bool {
		return l1.ID == id
	})

	return nil
}

func (s *Storage) GetAllMovies() ([]domain.Movie, error) {
	return s.movies, nil
}

func (s *Storage) SortAndOrderByMovie(sortBy, orderBy string, movies []domain.Movie) []domain.Movie {
	switch {
	case sortBy == "name" || sortBy == "":
		if orderBy == "" || orderBy == "asc" {
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].Name < movies[j].Name
			})
		} else {
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].Name > movies[j].Name
			})
		}

	case sortBy == "genre":
		if orderBy == "" || orderBy == "asc" {
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].Genre < movies[j].Genre
			})
		} else {
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].Genre > movies[j].Genre
			})
		}

	case sortBy == "date":
		if orderBy == "" || orderBy == "asc" {
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].ReleaseDate.Before(movies[j].ReleaseDate)
			})
		} else {
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].ReleaseDate.After(movies[j].ReleaseDate)
			})
		}
	}

	return movies
}

func (s *Storage) GetActorsByMovie(id int) ([]domain.Actor, error) {
	var actors []domain.Actor
	actorsIDs := s.actorsByMovie[id]
	for _, val := range actorsIDs {
		actorExists := slices.ContainsFunc(s.actors, func(actor domain.Actor) bool {
			return actor.ID == val
		})
		if !actorExists {
			return []domain.Actor{}, domain.ErrNotFound
		}
		if actorExists {
			actor, _ := s.GetActorByID(val) //!!
			actors = append(actors, actor)
		}
	}

	if len(actors) == 0 {
		return []domain.Actor{}, domain.ErrNotFound
	}

	return actors, nil
}

func (s *Storage) CreateActorsByMovie(id int, actors []int) (int, []int, error) {
	_, err := s.GetMovieByID(id)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return 0, []int{0}, domain.ErrNotExists
	case err != nil:
		return 0, []int{0}, fmt.Errorf("unexpected error %w", err)
	}

	for _, val := range actors {
		actorExists := slices.ContainsFunc(s.actors, func(actor domain.Actor) bool {
			return actor.ID == val
		})
		if !actorExists {
			return 0, []int{0}, domain.ErrNotFound
		}
	}

	//check if all actors exist

	s.actorsByMovie[id] = actors

	return id, actors, nil
}
