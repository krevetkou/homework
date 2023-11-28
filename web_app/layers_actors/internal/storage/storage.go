package storage

import (
	"arch-demo/layers_actors/internal/domain"
	"golang.org/x/exp/slices"
	"sort"
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
			strings.Contains(s.actors[i].Sex, actor.Sex) &&
			strings.Contains(s.actors[i].BirthYear, actor.BirthYear) &&
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

func (s *ActorsStorage) OrderBy(param string) []domain.Actor {
	var filteredActors []domain.Actor

	switch {
	case param == "name":
		names := make([]string, 0, len(s.actors))

		for _, val := range s.actors {
			names = append(names, val.Name)
		}
		sort.Strings(names)

		for _, valName := range names {
			for _, valActor := range s.actors {
				if valActor.Name == valName {
					filteredActors = append(filteredActors, valActor)
					break
				}
			}
		}

	case param == "country":
		countries := make([]string, 0, len(s.actors))

		for _, val := range s.actors {
			countries = append(countries, val.CountryOfBirth)
		}
		sort.Strings(countries)

		for _, valCountry := range countries {
			for _, valActor := range s.actors {
				if valActor.CountryOfBirth == valCountry {
					filteredActors = append(filteredActors, valActor)
					break
				}
			}
		}
	case param == "birthdate":
		birthDates := make([]string, 0, len(s.actors))

		for _, val := range s.actors {
			birthDates = append(birthDates, val.BirthYear)
		}
		sort.Strings(birthDates)

		for _, valBirthDate := range birthDates {
			for _, valActor := range s.actors {
				if valActor.BirthYear == valBirthDate {
					filteredActors = append(filteredActors, valActor)
					break
				}
			}
		}
	}

	return filteredActors
}
