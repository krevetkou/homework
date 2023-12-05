package storage

import (
	"arch-demo/internal/domain"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

func (s *Storage) InsertMovie(actor domain.Movie) domain.Movie {
	var lastID int
	if len(s.movies) > 0 {
		lastID = s.movies[len(s.movies)-1:][0].ID
	}

	actor.ID = lastID + 1

	s.movies = append(s.movies, actor)
	return actor
}

func (s *Storage) IsMovieExists(movie domain.Movie) bool {
	for i := range s.movies {
		if strings.Contains(s.movies[i].Name, movie.Name) &&
			s.movies[i].ReleaseDate == movie.ReleaseDate &&
			strings.Contains(s.movies[i].Country, movie.Country) &&
			strings.Contains(s.movies[i].Genre, movie.Genre) &&
			s.movies[i].Rating == movie.Rating {
			return true
		}
	}

	return false
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

func (s *Storage) UpdateMovie(movieUpdate domain.Movie) {
	for i := range s.movies {
		if s.movies[i].ID == movieUpdate.ID {
			s.movies[i] = movieUpdate
		}
	}
}

func (s *Storage) DeleteMovie(id int) {
	s.movies = slices.DeleteFunc(s.movies, func(l1 domain.Movie) bool {
		return l1.ID == id
	})
}

func (s *Storage) GetAllMovies() []domain.Movie {
	return s.movies
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
				return movies[i].ReleaseDate.Year > movies[j].ReleaseDate.Year
			})
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].ReleaseDate.Month > movies[j].ReleaseDate.Month
			})
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].ReleaseDate.Date > movies[j].ReleaseDate.Date
			})
		} else {
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].ReleaseDate.Year < movies[j].ReleaseDate.Year
			})
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].ReleaseDate.Month < movies[j].ReleaseDate.Month
			})
			sort.Slice(movies, func(i, j int) bool {
				return movies[i].ReleaseDate.Date < movies[j].ReleaseDate.Date
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

func (s *Storage) CreateActorsByMovie(id int, actors []int) error {
	_, err := s.GetMovieByID(id)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return err
	case err != nil:
		return fmt.Errorf("unexpected error %w", err)
	}

	for _, val := range actors {
		actorExists := slices.ContainsFunc(s.actors, func(actor domain.Actor) bool {
			return actor.ID == val
		})
		if !actorExists {
			return domain.ErrNotFound
		}
	}

	//check if all actors exist

	s.actorsByMovie[id] = actors

	return nil
}
