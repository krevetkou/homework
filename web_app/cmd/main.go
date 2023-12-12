package main

import (
	"arch-demo/internal/api"
	"arch-demo/internal/services"
	"arch-demo/internal/storage/db"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
)

func main() {
	dbCon, err := sql.Open("pgx", "postgres://postgres:356597@localhost:5432/test_db")
	if err != nil {
		log.Println(err)
		return
	}
	err = dbCon.Ping()
	if err != nil {
		log.Printf("failed to connect to db %v", err)
		return
	}

	dbStorage := db.NewDbStorage(dbCon)
	actorsService := services.NewActorService(dbStorage)
	actorsHandler := api.NewActorsHandler(actorsService)
	moviesService := services.NewMovieService(dbStorage)
	moviesHandler := api.NewLaptopsHandler(moviesService)

	//actors := []domain.Actor{
	//	{
	//		Name:           "Lol",
	//		BirthYear:      1909,
	//		CountryOfBirth: "Sos",
	//		Gender:         "female",
	//	}, {
	//		Name:           "Kek",
	//		BirthYear:      2012,
	//		CountryOfBirth: "Lock",
	//		Gender:         "male",
	//	}, {
	//		Name:           "Cheburek",
	//		BirthYear:      900,
	//		CountryOfBirth: "Klkd",
	//		Gender:         "elephant",
	//	}, {
	//		Name:           "Holk",
	//		BirthYear:      2765,
	//		CountryOfBirth: "Dada",
	//		Gender:         "fox",
	//	},
	//}
	//
	//movies := []domain.Movie{
	//	{
	//		Name: "Lol",
	//		ReleaseDate: domain.ReleaseDate{
	//			Date:  1,
	//			Month: 2,
	//			Year:  1995,
	//		},
	//		Country: "Poland",
	//		Genre:   "hehe",
	//		Rating:  5,
	//	}, {
	//		Name: "Kek",
	//		ReleaseDate: domain.ReleaseDate{
	//			Date:  1,
	//			Month: 3,
	//			Year:  1995,
	//		},
	//		Country: "Russland",
	//		Genre:   "nehehe",
	//		Rating:  2,
	//	}, {
	//		Name: "Cheburek",
	//		ReleaseDate: domain.ReleaseDate{
	//			Date:  1,
	//			Month: 2,
	//			Year:  2013,
	//		},
	//		Country: "Flugegenheimen",
	//		Genre:   "strashno",
	//		Rating:  909,
	//	}, {
	//		Name: "Shashlik",
	//		ReleaseDate: domain.ReleaseDate{
	//			Date:  5,
	//			Month: 5,
	//			Year:  2020,
	//		},
	//		Country: "UK",
	//		Genre:   "ploho",
	//		Rating:  1,
	//	},
	//}
	//
	//actorsForMovie := map[int][]int{
	//	1: {1, 2, 3},
	//	2: {2, 1, 4},
	//}
	//
	//for _, val := range movies {
	//	serviceStorage.InsertMovie(val)
	//}
	//
	//for _, val := range actors {
	//	serviceStorage.InsertActor(val)
	//}
	//
	//for ind, val := range actorsForMovie {
	//	err := serviceStorage.CreateActorsByMovie(ind, val)
	//	if err != nil {
	//		return
	//	}
	//}

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

	err = http.ListenAndServe(":8080", r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
		return
	}
	if err != nil {
		log.Printf("server error: %s", err)
	}
}
