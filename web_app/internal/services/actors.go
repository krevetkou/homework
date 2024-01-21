package services

import (
	"arch-demo/internal/domain"
	"errors"
	"fmt"
)

type ActorsRepository interface {
	InsertActor(actor domain.Actor) (domain.Actor, error)
	IsActorExists(actor domain.Actor) (bool, error)
	GetActorByID(id int) (domain.Actor, error)
	DeleteActor(id int) error
	UpdateActor(actor domain.Actor) error
	GetAllActors() ([]domain.Actor, error)
	SortAndOrderByActor(sortBy, orderBy string, actors []domain.Actor) []domain.Actor
	FilterActors(nameQuery, countryOfBirthQuery string) ([]domain.Actor, error)
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

	isActorExist, err := s.Storage.IsActorExists(actor)
	if isActorExist {
		return domain.Actor{}, err
	}

	newActor, err := s.Storage.InsertActor(actor)
	if err != nil {
		return domain.Actor{}, err
	}

	return newActor, nil
}

func (s ActorsService) Get(id int) (domain.Actor, error) {
	actor, err := s.Storage.GetActorByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.Actor{}, err
	}

	if err != nil {
		return domain.Actor{}, fmt.Errorf("failed to find actor, unexpected error: %w", err)
	}

	return actor, nil
}

func (s ActorsService) Update(id int, actorUpdate domain.ActorUpdate) (domain.Actor, error) {
	actor, err := s.Storage.GetActorByID(id)
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

	s.Storage.UpdateActor(actor)

	return actor, nil
}

func (s ActorsService) Delete(id int) error {
	_, err := s.Storage.GetActorByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("actor id: %d, err: %w", id, err)
	}

	if err != nil {
		return fmt.Errorf("failed to find actor, unexpected error: %w", err)
	}

	err = s.Storage.DeleteActor(id)
	if err != nil {
		return err
	}

	return nil
}

func (s ActorsService) List(sortBy, orderBy, nameQuery, countryOfBirthQuery string) ([]domain.Actor, error) {
	actors, err := s.Storage.GetAllActors()
	if err != nil {
		return []domain.Actor{}, err
	}
	var filteredActors []domain.Actor

	if nameQuery != "" || countryOfBirthQuery != "" {
		filteredActors, err = s.Storage.FilterActors(nameQuery, countryOfBirthQuery)
		if err != nil {
			return []domain.Actor{}, err
		}
	} else {
		filteredActors = actors
	}

	switch {
	case sortBy == "" && orderBy == "":
		actors = s.Storage.SortAndOrderByActor("name", "asc", filteredActors)
	case sortBy == "" && orderBy != "":
		actors = s.Storage.SortAndOrderByActor("name", orderBy, filteredActors)
	case sortBy != "" && orderBy == "":
		actors = s.Storage.SortAndOrderByActor(sortBy, "asc", filteredActors)
	default:
		actors = s.Storage.SortAndOrderByActor(sortBy, orderBy, filteredActors)
	}

	return actors, nil
}
