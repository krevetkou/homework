package api

import (
	"arch-demo/internal/domain"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type MoviesService interface {
	Create(actor domain.Movie) (domain.Movie, error)
	Get(id int) (domain.Movie, error)
	Delete(id int) error
	Update(id int, actorUpdate domain.MovieUpdate) (domain.Movie, error)
	List(orderBy, sortBy, nameQuery, genreQuery string) []domain.Movie
	GetActorsByMovie(id int) ([]domain.Actor, error)
	CreateActorsForMovie(id int, actorsByMovie []int) (int, []int, error)
}

type MoviesHandler struct {
	Service MoviesService
}

func NewLaptopsHandler(service MoviesService) MoviesHandler {
	return MoviesHandler{
		Service: service,
	}
}

//go:generate mockgen -source movies.go -destination ../tests/api_mocks/movies.go apimocks
func (h MoviesHandler) List(w http.ResponseWriter, r *http.Request) {
	SortBy := r.URL.Query().Get("sort")
	OrderBy := r.URL.Query().Get("order")
	NameQuery := r.URL.Query().Get("name")
	GenreQuery := r.URL.Query().Get("genre")

	filteredMovies := h.Service.List(SortBy, OrderBy, NameQuery, GenreQuery)
	data, err := json.Marshal(filteredMovies)
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

func (h MoviesHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content type not allowed", http.StatusUnsupportedMediaType)
		return
	}

	var newMovie domain.Movie
	err := json.NewDecoder(r.Body).Decode(&newMovie)
	if err != nil {
		http.Error(w, "failed to unmarshall data", http.StatusBadRequest)
		log.Println(err)
		return
	}

	createdMovie, err := h.Service.Create(newMovie)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFieldsRequired):
			http.Error(w, "all required fields must have values", http.StatusUnprocessableEntity)
		case errors.Is(err, domain.ErrExists):
			http.Error(w, "movie already exists", http.StatusConflict)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	data, err := json.Marshal(createdMovie)
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

func (h MoviesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := getID(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	movie, err := h.Service.Get(id)
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

	data, err := json.Marshal(movie)
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

func (h MoviesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getID(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	// более короткая и удобная запись вместо io.ReadAll
	var movieUpdate domain.MovieUpdate
	err = json.NewDecoder(r.Body).Decode(&movieUpdate)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to unmarshall request body", http.StatusBadRequest)
		return
	}

	updatedMovie, err := h.Service.Update(id, movieUpdate)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.Error(w, "movie not found", http.StatusNotFound)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	data, err := json.Marshal(updatedMovie)
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

func (h MoviesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getID(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.Service.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.Error(w, "movie not found", http.StatusNotFound)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h MoviesHandler) GetActors(w http.ResponseWriter, r *http.Request) {
	id, err := getID(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	actorsByMovie, err := h.Service.GetActorsByMovie(id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.Error(w, "actors not found", http.StatusNotFound)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}

		log.Println(err)
		return
	}

	data, err := json.Marshal(actorsByMovie)
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

func (h MoviesHandler) CreateActorsForMovie(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content type not allowed", http.StatusUnsupportedMediaType)
		return
	}

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

	var actorsForMovie []int
	err = json.NewDecoder(r.Body).Decode(&actorsForMovie)
	if err != nil {
		http.Error(w, "failed to unmarshall data", http.StatusBadRequest)
		log.Println(err)
		return
	}

	var actorsIDs []int
	_, actorsIDs, err = h.Service.CreateActorsForMovie(id, actorsForMovie)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			http.Error(w, "actors not found", http.StatusNotFound)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	data, err := json.Marshal(actorsIDs)
	if err != nil {
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
	}

	// если не передать content-type, то клиент воспримет контент как text/plain, а не json
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func getID(w http.ResponseWriter, r *http.Request) (int, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		log.Println("id required")
		http.Error(w, "id required", http.StatusBadRequest)
		return 0, domain.ErrIDRequired
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to parse id query param", http.StatusBadRequest)
		return 0, err
	}

	return id, nil
}
