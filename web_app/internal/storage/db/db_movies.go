package db

import (
	"arch-demo/internal/domain"
	"database/sql"
	"errors"
	"sort"
)

func (s *StorageDB) InsertMovie(movie domain.Movie) (domain.Movie, error) {
	query := `insert into movies (name, release_date, country, genre, rating) 
				values ($1, $2, $3, $4, $5) 
				returning id, name, release_date, country, genre, rating`
	var newMovie domain.Movie
	err := s.db.QueryRow(query, movie.Name, movie.ReleaseDate, movie.Country, movie.Genre, movie.Rating).
		Scan(&newMovie.ID, &newMovie.Name, &newMovie.ReleaseDate, &newMovie.Country, &newMovie.Genre, &newMovie.Rating)
	if err != nil {
		return domain.Movie{}, err
	}

	return newMovie, nil
}

func (s *StorageDB) IsMovieExists(movie domain.Movie) (bool, error) {
	query := `select id from movies where name = $1`
	var id int8
	err := s.db.QueryRow(query, movie.Name).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.ErrNotFound
		}
		return false, err
	}

	return true, nil
}

func (s *StorageDB) GetMovieByID(id int) (domain.Movie, error) {
	query := `select id, name, release_date, country, genre, rating from movies where id = $1`
	var newMovie domain.Movie
	err := s.db.QueryRow(query, id).Scan(&newMovie.ID, &newMovie.Name, &newMovie.ReleaseDate, &newMovie.Country, &newMovie.Genre, &newMovie.Rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Movie{}, err
		}
		return domain.Movie{}, err
	}

	return newMovie, nil
}

func (s *StorageDB) UpdateMovie(movieUpdate domain.Movie) error {
	query := `update movies set name = $1, release_date = $2, country = $3, genre = $4, rating = $5 where id = $6;`
	_, err := s.db.Exec(query, movieUpdate.Name, movieUpdate.ReleaseDate, movieUpdate.Country, movieUpdate.Genre, movieUpdate.Rating)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageDB) DeleteMovie(id int) error {
	query := `DELETE FROM movies WHERE id = $1;`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageDB) GetAllMovies() ([]domain.Movie, error) {
	rows, err := s.db.Query("select * from movies")
	if err != nil {
		return []domain.Movie{}, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var movies []domain.Movie
	for rows.Next() {
		var newMovie domain.Movie
		if err = rows.Scan(&newMovie.ID, &newMovie.Name, &newMovie.ReleaseDate, &newMovie.Country, &newMovie.Genre, &newMovie.Rating); err != nil {
			return []domain.Movie{}, err
		}
		movies = append(movies, newMovie)
	}
	if err = rows.Err(); err != nil {
		return []domain.Movie{}, err
	}

	return movies, nil
}

func (s *StorageDB) SortAndOrderByMovie(sortBy, orderBy string, movies []domain.Movie) []domain.Movie {
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

func (s *StorageDB) GetActorsByMovie(id int) ([]domain.Actor, error) {
	query := `select actors_ids from "actorsInMovies" where id = $1`
	var actorsIDs []int
	err := s.db.QueryRow(query, id).Scan(&actorsIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Actor{}, err
		}
		return []domain.Actor{}, err
	}
	actors, err := s.GetAllActors()
	if err != nil {
		return []domain.Actor{}, err
	}

	var filteredActors []domain.Actor
	for _, actor := range actors {
		for _, actorID := range actorsIDs {
			if actorID == actor.ID {
				filteredActors = append(filteredActors, actor)
				break
			}
		}
	}

	return filteredActors, nil
}

func (s *StorageDB) CreateActorsByMovie(id int, actors []int) error {
	return nil
}
