package db

import (
	"arch-demo/internal/domain"
	"database/sql"
	"errors"
	"sort"
	"strings"
)

type StorageDB struct {
	db *sql.DB
}

func NewDbStorage(dbCon *sql.DB) *StorageDB {
	return &StorageDB{
		db: dbCon,
	}
}

func (s *StorageDB) InsertActor(actor domain.Actor) (domain.Actor, error) {
	query := `insert into actors (name, birth_year, country_of_birth, gender) values ($1, $2, $3, $4) returning id, name, birth_year, country_of_birth, gender`
	var newActor domain.Actor
	err := s.db.QueryRow(query, actor.Name, actor.BirthYear, actor.CountryOfBirth, actor.Gender).Scan(&newActor.ID, &newActor.Name, &newActor.BirthYear, &newActor.CountryOfBirth, &newActor.Gender)
	if err != nil {
		return domain.Actor{}, err
	}

	return newActor, nil
}

func (s *StorageDB) IsActorExists(actor domain.Actor) (bool, error) {
	query := `select id from actors where name = $1`
	var id int8
	err := s.db.QueryRow(query, actor.Name).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, domain.ErrNotFound
		}
		return false, err
	}

	return true, nil
}

func (s *StorageDB) GetActorByID(id int) (domain.Actor, error) {
	query := `select id, name, birth_year, country_of_birth, gender from actors where id = $1`
	var newActor domain.Actor
	err := s.db.QueryRow(query, id).Scan(&newActor.ID, &newActor.Name, &newActor.BirthYear, &newActor.CountryOfBirth, &newActor.Gender)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Actor{}, err
		}
		return domain.Actor{}, err
	}

	return newActor, nil
}

func (s *StorageDB) DeleteActor(id int) error {
	query := `DELETE FROM actors WHERE id = $1;`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageDB) UpdateActor(actorUpdate domain.Actor) error {
	query := `update users set name = $1, birth_year = $2, country_of_birth = $3, gender = $4 where id = $5;`
	_, err := s.db.Exec(query, actorUpdate.Name, actorUpdate.BirthYear, actorUpdate.CountryOfBirth, actorUpdate.Gender, actorUpdate.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageDB) GetAllActors() ([]domain.Actor, error) {
	rows, err := s.db.Query("select * from actors")
	if err != nil {
		return []domain.Actor{}, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var actors []domain.Actor
	for rows.Next() {
		var actor domain.Actor
		if err = rows.Scan(&actor.ID, &actor.Name, &actor.BirthYear, &actor.CountryOfBirth, &actor.Gender); err != nil {
			return []domain.Actor{}, err
		}
		actors = append(actors, actor)
	}
	if err = rows.Err(); err != nil {
		return []domain.Actor{}, err
	}

	return actors, nil
}

func (s *StorageDB) SortAndOrderByActor(sortOrder, orderBy string, actors []domain.Actor) []domain.Actor {
	switch {
	case sortOrder == "name":
		if orderBy == "" || orderBy == "asc" {
			sort.Slice(actors, func(i, j int) bool {
				return actors[i].Name < actors[j].Name
			})
		} else {
			sort.Slice(actors, func(i, j int) bool {
				return actors[i].Name > actors[j].Name
			})
		}
	case sortOrder == "country":
		if orderBy == "" || orderBy == "asc" {
			sort.Slice(actors, func(i, j int) bool {
				return actors[i].CountryOfBirth < actors[j].CountryOfBirth
			})
		} else {
			sort.Slice(actors, func(i, j int) bool {
				return actors[i].CountryOfBirth > actors[j].CountryOfBirth
			})
		}
	case sortOrder == "birthdate":
		if orderBy == "" || orderBy == "asc" {
			sort.Slice(actors, func(i, j int) bool {
				return actors[i].BirthYear < actors[j].BirthYear
			})
		} else {
			sort.Slice(actors, func(i, j int) bool {
				return actors[i].BirthYear > actors[j].BirthYear
			})
		}
	}

	return actors
}

func (s *StorageDB) FilterActors(nameQuery, countryOfBirthQuery string) ([]domain.Actor, error) {
	var filteredActors []domain.Actor
	actors, err := s.GetAllActors()
	if err != nil {
		return []domain.Actor{}, err
	}

	for i := range actors {
		if (nameQuery != "" && strings.Contains(actors[i].Name, nameQuery)) ||
			(countryOfBirthQuery != "" && strings.Contains(actors[i].CountryOfBirth, countryOfBirthQuery)) {
			filteredActors = append(filteredActors, actors[i])
		}
	}

	return filteredActors, nil
}
