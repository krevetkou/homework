package main

import (
	"arch-demo/layers_movies/internal/api"
	"arch-demo/layers_movies/internal/domain"
	"arch-demo/layers_movies/internal/services"
	"arch-demo/layers_movies/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	moviesStorage := storage.NewMovieStorage()
	moviesService := services.NewMovieService(moviesStorage)
	moviesHandler := api.NewLaptopsHandler(moviesService)

	actors := []domain.Movie{
		{
			Name: "Lol",
			ReleaseDate: domain.ReleaseDate{
				Date:  1,
				Month: 2,
				Year:  1995,
			},
			Country: "Poland",
			Genre:   "hehe",
			Rating:  5,
		}, {
			Name: "Kek",
			ReleaseDate: domain.ReleaseDate{
				Date:  1,
				Month: 3,
				Year:  1995,
			},
			Country: "Russland",
			Genre:   "nehehe",
			Rating:  2,
		}, {
			Name: "Cheburek",
			ReleaseDate: domain.ReleaseDate{
				Date:  1,
				Month: 2,
				Year:  2013,
			},
			Country: "Flugegenheimen",
			Genre:   "strashno",
			Rating:  909,
		}, {
			Name: "Shashlik",
			ReleaseDate: domain.ReleaseDate{
				Date:  5,
				Month: 5,
				Year:  2020,
			},
			Country: "UK",
			Genre:   "ploho",
			Rating:  1,
		},
	}

	for _, val := range actors {
		moviesStorage.Insert(val)
	}

	//POST /movies - добавление нового фильма
	//PATCH /movies/{id} - частичное обновление фильма, можно обновить любое значение
	//DELETE /movies/{id} - удаление фильма по id
	//GET /movies/{id} - получение одного фильма по id
	//GET /movies - получение списка фильмов, добавить возможность фильтровать по названию фильма и жанру
	//с помощью query параметров /movies?name и /movies?genre. Добавить возможность упорядочить вывод по
	//названию, жанру и дате выхода с помощью query параметра /movies?order="name", /movies?order="name",
	///movies?order="genre" и /actors?order="date". Добавить возможность отсортировать вывод в порядке
	//убывания или возростания с помощью query параметра /actors?sort="asc" , /actors?sort="desc"

	//POST /movies/{movie_id}/actors - добавление в фильм списка актеров - в теле запроса необходимо передать массив id актеров
	//GET /movies/{movie_id}/actors - получение списка актеров в фильме, возвращается полная информация о всех актерах

	r := chi.NewRouter()
	r.Route("/movies", func(r chi.Router) {
		r.Post("/", moviesHandler.Create)
		r.Get("/", moviesHandler.List)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", moviesHandler.Get)
			r.Patch("/", moviesHandler.Update)
			r.Delete("/", moviesHandler.Delete)
		})
	})

	err := http.ListenAndServe(":8080", r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
		return
	}
	if err != nil {
		log.Printf("server error: %s", err)
	}
}
