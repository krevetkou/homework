package api

import (
	"arch-demo/layers_movies/internal/domain"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
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
}

type MoviesHandler struct {
	Service MoviesService
}

func NewLaptopsHandler(service MoviesService) MoviesHandler {
	return MoviesHandler{
		Service: service,
	}
}

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
		return
	}
}

func (h MoviesHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var newMovie domain.Movie
	err = json.Unmarshal(body, &newMovie)
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
		case errors.Is(err, domain.ErrMovieExists):
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
		return
	}
}

func (h MoviesHandler) Update(w http.ResponseWriter, r *http.Request) {
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
		return
	}
}

func (h MoviesHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "movie not found", http.StatusNotFound)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	w.WriteHeader(http.StatusAccepted)
}
