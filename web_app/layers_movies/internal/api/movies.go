package api

import "net/http"

type MoviesService interface {
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

}

func (h MoviesHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h MoviesHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h MoviesHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h MoviesHandler) Delete(w http.ResponseWriter, r *http.Request) {

}
