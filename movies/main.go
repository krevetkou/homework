package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Movies map[string][]string

func main() {
	movies := Movies{}
	var txt string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		txt = scanner.Text()

		t := strings.Split(txt, " ")
		switch t[0] {
		case "add":
			movies.addMovie(txt)
		case "remove":
			movies.removeActorOrMovie(txt)
		case "dump":
			movies.saveMovie(txt)
		case "exit":
			os.Exit(0)
		}
	}

}

func (m Movies) addMovie(text string) {
	movie := strings.Split(text, " m ")
	movie = strings.Split(movie[1], " a ") // [0] - movie, [1] - actor

	if _, ok := m[movie[0]]; !ok {
		m[movie[0]] = []string{}
	}

	m[movie[0]] = append(m[movie[0]], movie[1])
}

func (m Movies) removeActorOrMovie(text string) {
	movieAndActor := strings.Split(text, " m ")
	movieAndActor = strings.Split(movieAndActor[1], " a ") // [0] - movie, [1] - actor
	actors := m[movieAndActor[0]]
	i := 0

	if len(movieAndActor) == 1 {
		delete(m, movieAndActor[0])
	} else {
		for ind, val := range actors {
			if val == movieAndActor[1] {
				i = ind
				break
			}
		}

		actors[i] = actors[len(actors)-1]
		actors[len(actors)-1] = ""
		actors = actors[:len(actors)-1]
	}
}

func (m Movies) saveMovie(text string) {
	err := os.Mkdir("moviesDB", 0755)
	if err != nil {
		os.Exit(0)
	}

	movie := strings.Split(text, " m ") // [1] - movie

	if len(movie) > 1 {
		save(movie[1], m[movie[1]])
	} else {
		for ind, val := range m {
			save(ind, val)
		}
	}
}

func save(movie string, actors []string) {
	txt := "moviesDB/" + movie
	file, err := os.Create(txt)
	if err != nil {
		os.Exit(1)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	_, err = io.WriteString(file, strings.Join(actors, ", "))
	if err != nil {
		os.Exit(1)
	}
}
