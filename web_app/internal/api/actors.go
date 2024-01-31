package api

import (
	"arch-demo/internal/domain"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ActorsService interface {
	Create(actor domain.Actor) (domain.Actor, error)
	Get(id int) (domain.Actor, error)
	Delete(id int) error
	Update(id int, actorUpdate domain.ActorUpdate) (domain.Actor, error)
	List(sortBy, orderBy, nameQuery, countryOfBirthQuery string) ([]domain.Actor, error)
}

type ActorsHandler struct {
	Service ActorsService
}

func NewActorsHandler(service ActorsService) ActorsHandler {
	return ActorsHandler{
		Service: service,
	}
}

//go:generate mockgen -source actors.go -destination ../tests/api_mocks/users.go apimocks

func (h ActorsHandler) Create(w http.ResponseWriter, r *http.Request) {
	// необходимо удостоверится, что в запросе контент нужного типа
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content type not allowed", http.StatusUnsupportedMediaType)
		return
	}

	var newActor domain.Actor
	err := json.NewDecoder(r.Body).Decode(&newActor)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	createdActor, err := h.Service.Create(newActor)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFieldsRequired):
			http.Error(w, "all required fields must have values", http.StatusUnprocessableEntity)
		case errors.Is(err, domain.ErrExists):
			http.Error(w, "actor already exists", http.StatusConflict)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	data, err := json.Marshal(createdActor)
	if err != nil {
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h ActorsHandler) List(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	orderBy := r.URL.Query().Get("order")
	nameQuery := r.URL.Query().Get("name")
	countryOfBirthQuery := r.URL.Query().Get("country")

	filteredActors, err := h.Service.List(sortBy, orderBy, nameQuery, countryOfBirthQuery)
	if err != nil {
		http.Error(w, "failed to get actors", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(filteredActors)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		return
	}

	// если не передать content-type, то клиент воспримет контент как text/plain, а не json
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h ActorsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := getID(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	actor, err := h.Service.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.Error(w, "actor not found", http.StatusNotFound)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}

		log.Println(err)
		return
	}

	data, err := json.Marshal(actor)
	if err != nil {
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
	}

	// если не передать content-type, то клиент воспримет контент как text/plain, а не json
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h ActorsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getID(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	// более короткая и удобная запись вместо io.ReadAll
	var actorUpdate domain.ActorUpdate
	err = json.NewDecoder(r.Body).Decode(&actorUpdate)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to unmarshall request body", http.StatusBadRequest)
		return
	}

	updatedActor, err := h.Service.Update(id, actorUpdate)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.Error(w, "actor not found", http.StatusNotFound)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	data, err := json.Marshal(updatedActor)
	if err != nil {
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		return
	}

	// если не передать content-type, то клиент воспримет контент как text/plain, а не json
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h ActorsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getID(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.Service.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.Error(w, "actor not found", http.StatusNotFound)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	w.WriteHeader(http.StatusAccepted)
}
