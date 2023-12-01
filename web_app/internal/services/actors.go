package services

import (
	"arch-demo/internal/domain"
	"errors"
	"fmt"
	"strings"
)

type ActorsRepository interface {
	Insert(actor domain.Actor) domain.Actor
	IsActorExists(actor domain.Actor) bool
	GetByID(id int) (domain.Actor, error)
	Delete(id int)
	Update(actor domain.Actor)
	GetAll() []domain.Actor
	SortAndOrderBy(sortBy, orderBy string, actors []domain.Actor) []domain.Actor
}

type ActorsService struct {
	Storage ActorsRepository
}

func NewActorService(storage ActorsRepository) ActorsService {
	return ActorsService{
		Storage: storage,
	}
}

func (s ActorsService) Create(actor domain.Actor) (domain.Actor, error) {
	// входящие параметры необходимо валидировать
	if actor.Name == "" || actor.Gender == "" || actor.BirthYear == 0 || actor.CountryOfBirth == "" {
		return domain.Actor{}, domain.ErrFieldsRequired
	}

	isActorExist := s.Storage.IsActorExists(actor)
	if isActorExist {
		return domain.Actor{}, domain.ErrExists
	}

	newActor := s.Storage.Insert(actor)

	return newActor, nil
}

func (s ActorsService) Get(id int) (domain.Actor, error) {
	actor, err := s.Storage.GetByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.Actor{}, err
	}

	if err != nil {
		return domain.Actor{}, fmt.Errorf("failed to find actor, unexpected error: %w", err)
	}

	return actor, nil
}

func (s ActorsService) Update(id int, actorUpdate domain.ActorUpdate) (domain.Actor, error) {
	actor, err := s.Storage.GetByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.Actor{}, err
	}

	if actorUpdate.Name != nil {
		actor.Name = *actorUpdate.Name
	}

	if actorUpdate.BirthYear != nil {
		actor.BirthYear = *actorUpdate.BirthYear
	}

	if actorUpdate.CountryOfBirth != nil {
		actor.CountryOfBirth = *actorUpdate.CountryOfBirth
	}

	if actorUpdate.Sex != nil {
		actor.Gender = *actorUpdate.Sex
	}

	s.Storage.Update(actor)

	return actor, nil
}

func (s ActorsService) Delete(id int) error {
	_, err := s.Storage.GetByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("actor id: %d, err: %w", id, err)
	}

	if err != nil {
		return fmt.Errorf("failed to find actor, unexpected error: %w", err)
	}

	s.Storage.Delete(id)

	return nil
}

func (s ActorsService) List(sortBy, orderBy, nameQuery, countryOfBirthQuery string) []domain.Actor {
	actors := s.Storage.GetAll()
	var filteredActors []domain.Actor

	if nameQuery != "" || countryOfBirthQuery != "" {
		for i := range actors {
			if (nameQuery != "" && strings.Contains(actors[i].Name, nameQuery)) ||
				(countryOfBirthQuery != "" && strings.Contains(actors[i].CountryOfBirth, countryOfBirthQuery)) {
				filteredActors = append(filteredActors, actors[i])
			}
		}
	} else {
		filteredActors = actors
	}

	switch {
	case sortBy == "" && orderBy == "":
		actors = s.Storage.SortAndOrderBy("name", "asc", filteredActors)
	case sortBy == "" && orderBy != "":
		actors = s.Storage.SortAndOrderBy("name", orderBy, filteredActors)
	case sortBy != "" && orderBy == "":
		actors = s.Storage.SortAndOrderBy(sortBy, "asc", filteredActors)
	default:
		actors = s.Storage.SortAndOrderBy(sortBy, orderBy, filteredActors)
	}

	return actors
}
