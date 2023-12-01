package storage

import (
	"arch-demo/layers_actors/internal/domain"
	"golang.org/x/exp/slices"
	"sort"
	"strconv"
	"strings"
)

type ActorsStorage struct {
	actors []domain.Actor
}

func NewActorStorage() *ActorsStorage {
	return &ActorsStorage{
		actors: make([]domain.Actor, 0),
	}
}

func (s *ActorsStorage) Insert(actor domain.Actor) domain.Actor {
	var lastID int
	if len(s.actors) > 0 {
		lastID = s.actors[len(s.actors)-1:][0].ID
	}

	actor.ID = lastID + 1

	s.actors = append(s.actors, actor)
	return actor
}

func (s *ActorsStorage) IsActorExists(actor domain.Actor) bool {
	for i := range s.actors {
		if strings.Contains(s.actors[i].Name, actor.Name) &&
			strings.Contains(s.actors[i].Gender, actor.Gender) &&
			strings.Contains(strconv.Itoa(s.actors[i].BirthYear), strconv.Itoa(actor.BirthYear)) &&
			strings.Contains(s.actors[i].CountryOfBirth, actor.CountryOfBirth) {
			return true
		}
	}

	return false
}

func (s *ActorsStorage) GetByID(id int) (domain.Actor, error) {
	var actor *domain.Actor
	for i := range s.actors {
		if s.actors[i].ID == id {
			actor = &s.actors[i]
		}
	}

	if actor == nil {
		return domain.Actor{}, domain.ErrNotFound
	}

	return *actor, nil
}

func (s *ActorsStorage) Delete(id int) {
	s.actors = slices.DeleteFunc(s.actors, func(l1 domain.Actor) bool {
		return l1.ID == id
	})
}

func (s *ActorsStorage) Update(actorUpdate domain.Actor) {
	for i := range s.actors {
		if s.actors[i].ID == actorUpdate.ID {
			s.actors[i] = actorUpdate
		}
	}
}

func (s *ActorsStorage) GetAll() []domain.Actor {
	return s.actors
}

func (s *ActorsStorage) SortAndOrderBy(sortBy, orderBy string) []domain.Actor {
	actors := s.GetAll()

	switch {
	case sortBy == "name":
		if orderBy == "" || orderBy == "asc" {
			sort.Slice(actors, func(i, j int) bool {
				return actors[i].Name < actors[j].Name
			})
		} else {
			sort.Slice(actors, func(i, j int) bool {
				return actors[i].Name > actors[j].Name
			})
		}
	case sortBy == "country":
		if orderBy == "" || orderBy == "asc" {
			sort.Slice(actors, func(i, j int) bool {
				return actors[i].CountryOfBirth < actors[j].CountryOfBirth
			})
		} else {
			sort.Slice(actors, func(i, j int) bool {
				return actors[i].CountryOfBirth > actors[j].CountryOfBirth
			})
		}
	case sortBy == "birthdate":
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
