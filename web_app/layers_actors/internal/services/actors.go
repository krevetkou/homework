package services

import (
	"arch-demo/layers_actors/internal/domain"
	"errors"
	"fmt"
	"log"
	"strings"
)

type ActorsRepository interface {
	Insert(actor domain.Actor) domain.Actor
	IsActorExists(actor domain.Actor) bool
	GetByID(id int) (domain.Actor, error)
	Delete(id int)
	Update(actor domain.Actor)
	GetAll() []domain.Actor
	OrderBy(param string) []domain.Actor
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
	if actor.Name == "" || actor.Sex == "" || actor.BirthYear == "" || actor.CountryOfBirth == "" {
		return domain.Actor{}, domain.ErrFieldsRequired
	}

	isActorExist := s.Storage.IsActorExists(actor)
	if isActorExist {
		return domain.Actor{}, domain.ErrActorExists
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
		actor.Sex = *actorUpdate.Sex
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

func (s ActorsService) List(name, countryOfBirth, orderBy string) []domain.Actor {
	actors := s.Storage.GetAll()
	if name == "" && countryOfBirth == "" && orderBy == "" {
		return actors
	}

	var filteredActors []domain.Actor

	if orderBy != "" {
		filteredActors = s.Storage.OrderBy(orderBy)
		return filteredActors
	}

	for i := range actors {
		if (name != "" && strings.Contains(actors[i].Name, name)) ||
			(countryOfBirth != "" && strings.Contains(actors[i].CountryOfBirth, countryOfBirth)) {
			filteredActors = append(filteredActors, actors[i])
			log.Println(countryOfBirth)
		}
	}

	return filteredActors
}
