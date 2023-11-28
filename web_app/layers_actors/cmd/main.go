package main

import (
	"arch-demo/layers_actors/internal/api"
	"arch-demo/layers_actors/internal/services"
	"arch-demo/layers_actors/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	actorsStorage := storage.NewActorStorage()
	actorsService := services.NewActorService(actorsStorage)
	actorsHandler := api.NewActorsHandler(actorsService)

	//POST /actors - добавление нового актера
	//PATCH /actors/{id} - частичное обновление актера, можно обновить любое значение
	//DELETE /actors/{id} - удаление актера по id
	//GET /actors/{id} - получение одного актера по id
	//GET /actors - получение списка актеров, добавить возможность фильтровать по имени актера и стране рождения с помощью query параметров /actors?name и /actors?country. Добавить возможность упорядочить вывод по имени, стране и году рождения с помощью query параметра /actors?order="name", /actors?order="country" /actors?order="birthdate". Добавить возможность отсортировать вывод в порядке убывания или возростания с помощью query параметра /actors?sort="asc" , /actors?sort="desc"

	r := chi.NewRouter()
	r.Route("/actors", func(r chi.Router) {
		r.Post("/", actorsHandler.Create) //добавление нового актера
		r.Get("/", actorsHandler.List)    //получение списка актеров,
		// добавить возможность фильтровать по имени актера и стране рождения с помощью query параметров
		///actors?name=a и /actors?country=a. Добавить возможность упорядочить вывод по имени,
		//стране и году рождения с помощью query параметра /actors?order=name, /actors?order=country
		///actors?order=birthdate. Добавить возможность отсортировать вывод в порядке убывания или возрастания
		//с помощью query параметра /actors?sort="asc" , /actors?sort="desc"

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", actorsHandler.Get)       //получение одного актера по id
			r.Patch("/", actorsHandler.Update)  //частичное обновление актера, можно обновить любое значение
			r.Delete("/", actorsHandler.Delete) //удаление актера по id
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
