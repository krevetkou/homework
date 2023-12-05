package inmemory

import (
	"arch-demo/internal/domain"
	"golang.org/x/exp/slices"
	"sort"
	"strconv"
	"strings"
)

type Storage struct {
	actors        []domain.Actor
	movies        []domain.Movie
	actorsByMovie map[int][]int
}

func NewStorage() *Storage {
	return &Storage{
		actors:        make([]domain.Actor, 0),
		movies:        make([]domain.Movie, 0),
		actorsByMovie: make(map[int][]int),
	}
}

func (s *Storage) InsertActor(actor domain.Actor) (domain.Actor, error) {
	var lastID int
	if len(s.actors) > 0 {
		lastID = s.actors[len(s.actors)-1:][0].ID
	}

	actor.ID = lastID + 1

	s.actors = append(s.actors, actor)
	return actor, nil
}

func (s *Storage) IsActorExists(actor domain.Actor) (bool, error) {
	for i := range s.actors {
		if strings.Contains(s.actors[i].Name, actor.Name) &&
			strings.Contains(s.actors[i].Gender, actor.Gender) &&
			strings.Contains(strconv.Itoa(s.actors[i].BirthYear), strconv.Itoa(actor.BirthYear)) &&
			strings.Contains(s.actors[i].CountryOfBirth, actor.CountryOfBirth) {
			return true, nil
		}
	}

	return false, domain.ErrExists
}

func (s *Storage) GetActorByID(id int) (domain.Actor, error) {
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

func (s *Storage) DeleteActor(id int) {
	s.actors = slices.DeleteFunc(s.actors, func(l1 domain.Actor) bool {
		return l1.ID == id
	})
}

func (s *Storage) UpdateActor(actorUpdate domain.Actor) error {
	for i := range s.actors {
		if s.actors[i].ID == actorUpdate.ID {
			s.actors[i] = actorUpdate
		}
	}

	return nil
}

func (s *Storage) GetAllActor() ([]domain.Actor, error) {
	return s.actors, nil
}

func (s *Storage) SortAndOrderByActor(sortBy, orderBy string, actors []domain.Actor) []domain.Actor {
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
