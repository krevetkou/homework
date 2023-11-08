package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Movies map[string][]string

const dir = "moviesDB"

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
	movieAndActor, err := checkErrorsAfterSplit(text)

	if err != 0 {
		return
	}

	movie := movieAndActor[0]
	actor := movieAndActor[1]

	if _, ok := m[movie]; !ok {
		m[movie] = []string{}
	}

	m[movie] = append(m[movie], actor)
}

func (m Movies) removeActorOrMovie(text string) {
	movieAndActor, err := checkErrorsAfterSplit(text)

	if err != 0 {
		return
	}

	movie := movieAndActor[0]
	actor := movieAndActor[1]

	actors := m[movie]
	i := 0

	if len(movieAndActor) == 1 {
		delete(m, movie)
		return
	}

	for ind, val := range actors {
		if val == actor {
			i = ind
			break
		}
	}

	remove(actors, i)
}

func (m Movies) saveMovie(text string) {
	_, err := exists(dir)
	if err != nil {
		return
	}

	_, err = os.Create(dir)
	if err != nil {
		return
	}

	movies := strings.Split(text, " m ") // [1] - movie
	movie := movies[1]

	if len(movies) > 1 {
		save(movie, m[movie])
		return
	}
	for ind, val := range m {
		save(ind, val)
	}

}

func save(movie string, actors []string) {
	_, err := exists(dir + "/" + movie)
	if err != nil {
		return
	}

	txt := dir + "/" + movie
	file, err := os.Create(txt)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	_, err = io.WriteString(file, strings.Join(actors, ", "))
	if err != nil {
		return
	}
}

func checkErrorsAfterSplit(text string) ([]string, int) {
	movieAndActor := strings.Split(text, " m ")
	err := 0
	if len(movieAndActor) != 2 {
		fmt.Printf("Вы не ввели, какой фильм и актера необходимо добавить %s\n", movieAndActor)
		err = 1
	}

	movieAndActor = strings.Split(movieAndActor[1], " a ")
	if len(movieAndActor) != 2 {
		fmt.Printf("Вы не ввели либо фильм, либо актера, которого необходимо добавить %s\n", movieAndActor)
		err = 1
	}

	return movieAndActor, err
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
