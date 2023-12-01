package main

import (
	api2 "arch-demo/internal/api"
	domain2 "arch-demo/internal/domain"
	services2 "arch-demo/internal/services"
	storage2 "arch-demo/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	storage := storage2.NewStorage()
	actorsService := services2.NewActorService(storage)
	actorsHandler := api2.NewActorsHandler(actorsService)
	moviesService := services2.NewMovieService(storage)
	moviesHandler := api2.NewLaptopsHandler(moviesService)

	actors := []domain2.Actor{
		{
			Name:           "Lol",
			BirthYear:      1909,
			CountryOfBirth: "Sos",
			Gender:         "female",
		}, {
			Name:           "Kek",
			BirthYear:      2012,
			CountryOfBirth: "Lock",
			Gender:         "male",
		}, {
			Name:           "Cheburek",
			BirthYear:      900,
			CountryOfBirth: "Klkd",
			Gender:         "elephant",
		}, {
			Name:           "Holk",
			BirthYear:      2765,
			CountryOfBirth: "Dada",
			Gender:         "fox",
		},
	}

	movies := []domain2.Movie{
		{
			Name: "Lol",
			ReleaseDate: domain2.ReleaseDate{
				Date:  1,
				Month: 2,
				Year:  1995,
			},
			Country: "Poland",
			Genre:   "hehe",
			Rating:  5,
		}, {
			Name: "Kek",
			ReleaseDate: domain2.ReleaseDate{
				Date:  1,
				Month: 3,
				Year:  1995,
			},
			Country: "Russland",
			Genre:   "nehehe",
			Rating:  2,
		}, {
			Name: "Cheburek",
			ReleaseDate: domain2.ReleaseDate{
				Date:  1,
				Month: 2,
				Year:  2013,
			},
			Country: "Flugegenheimen",
			Genre:   "strashno",
			Rating:  909,
		}, {
			Name: "Shashlik",
			ReleaseDate: domain2.ReleaseDate{
				Date:  5,
				Month: 5,
				Year:  2020,
			},
			Country: "UK",
			Genre:   "ploho",
			Rating:  1,
		},
	}

	for _, val := range movies {
		moviesStorage.Insert(val)
	}

	for _, val := range actors {
		actorsStorage.InsertActor(val)
	}

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Route("/actors", func(r chi.Router) {
			r.Post("/", actorsHandler.Create) //добавление нового актера
			r.Get("/", actorsHandler.List)    //получение списка актеров

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", actorsHandler.Get)       //получение одного актера по id
				r.Patch("/", actorsHandler.Update)  //частичное обновление актера, можно обновить любое значение
				r.Delete("/", actorsHandler.Delete) //удаление актера по id
			})
		})

		r.Route("/movies", func(r chi.Router) {
			r.Post("/", moviesHandler.Create)
			r.Get("/", moviesHandler.List)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", moviesHandler.Get)
				r.Patch("/", moviesHandler.Update)
				r.Delete("/", moviesHandler.Delete)
				r.Route("/actors", func(r chi.Router) {
					r.Get("/", moviesHandler.GetActors)
					r.Post("/", moviesHandler.CreateActorsForMovie)
				})
			})

		})

	})

	//POST /movies/{movie_id}/actors - добавление в фильм списка актеров - в теле запроса необходимо передать массив id актеров
	//GET /movies/{movie_id}/actors - получение списка актеров в фильме, возвращается полная информация о всех актерах

	err := http.ListenAndServe(":8080", r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
		return
	}
	if err != nil {
		log.Printf("server error: %s", err)
	}
}
