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
	dbCon, err := sql.Open("pgx", "postgres://postgres:356597@localhost:5432/postgres")
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
	moviesHandler := api.NewMoviesHandler(moviesService)

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
