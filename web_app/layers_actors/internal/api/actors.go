package api

import (
	"arch-demo/layers_actors/internal/domain"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"strconv"
)

type ActorsService interface {
	Create(actor domain.Actor) (domain.Actor, error)
	Get(id int) (domain.Actor, error)
	Delete(id int) error
	Update(id int, actorUpdate domain.ActorUpdate) (domain.Actor, error)
	List(sortBy, orderBy, nameQuery, countryOfBirthQuery string) []domain.Actor
}

type ActorsHandler struct {
	Service ActorsService
}

func NewActorsHandler(service ActorsService) ActorsHandler {
	return ActorsHandler{
		Service: service,
	}
}

func (h ActorsHandler) Create(w http.ResponseWriter, r *http.Request) {
	// необходимо удостоверится, что в запросе контент нужного типа
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content type not allowed", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	var newActor domain.Actor
	err = json.Unmarshal(body, &newActor)
	if err != nil {
		http.Error(w, "failed to unmarshall data", http.StatusBadRequest)
		log.Println(err)
		return
	}

	createdActor, err := h.Service.Create(newActor)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFieldsRequired):
			http.Error(w, "all required fields must have values", http.StatusUnprocessableEntity)
		case errors.Is(err, domain.ErrActorExists):
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
	SortBy := r.URL.Query().Get("sort")
	OrderBy := r.URL.Query().Get("order")
	NameQuery := r.URL.Query().Get("name")
	CountryOfBirthQuery := r.URL.Query().Get("country")

	filteredActors := h.Service.List(SortBy, OrderBy, NameQuery, CountryOfBirthQuery)
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
		return
	}
}

func (h ActorsHandler) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		log.Println("id required")
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to parse id query param", http.StatusBadRequest)
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
		return
	}
}

func (h ActorsHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		log.Println("id required")
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to parse id query param", http.StatusBadRequest)
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
	}

	// если не передать content-type, то клиент воспримет контент как text/plain, а не json
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func (h ActorsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		log.Println("id required")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "failed to parse id query param", http.StatusBadRequest)
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
