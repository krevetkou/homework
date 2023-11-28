package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
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
		case "load":
			movies.loadMovie(txt)
		case "exit":
			os.Exit(0)
		}
	}

}

func (m Movies) addMovie(text string) {
	movieAndActor := strings.Split(text, " m ")
	if len(movieAndActor) != 2 {
		fmt.Printf("Вы не ввели, какой фильм и актера необходимо добавить %s\n", movieAndActor)
		log.Println("not enough arguments")
		return
	}

	movieAndActor = strings.Split(movieAndActor[1], " a ")
	if len(movieAndActor) != 2 {
		fmt.Printf("Вы не ввели либо фильм, либо актера, которого необходимо добавить %s\n", movieAndActor)
		log.Println("not enough arguments")
		return
	}

	movie := movieAndActor[0]
	actor := movieAndActor[1]

	if _, ok := m[movie]; !ok {
		m[movie] = []string{}
	}

	m[movie] = append(m[movie], actor)

	log.Println(actor + " was added to movie " + movie)
}

func (m Movies) removeActorOrMovie(text string) {
	movieAndActor := strings.Split(text, " m ")
	if len(movieAndActor) != 2 {
		fmt.Printf("Вы не ввели, какой фильм и актера необходимо добавить %s\n", movieAndActor)
		log.Println("not enough arguments")
		return
	}

	movieAndActor = strings.Split(movieAndActor[1], " a ")
	movie := movieAndActor[0]

	if len(movieAndActor) != 2 {
		delete(m, movie)

		log.Println(movie + " was deleted")
		return
	}

	actor := movieAndActor[1]

	i := 0

	if len(movieAndActor) == 1 {
		delete(m, movie)

		log.Println(movie + " was deleted")
		return
	}

	for ind, val := range m[movie] {
		if val == actor {
			i = ind
			break
		}
	}

	remove(m[movie], i)

	log.Println(actor + " was deleted from " + movie)
}

func (m Movies) saveMovie(text string) {
	dirExists, err := exists(dir)
	if err != nil {
		return
	}

	if !dirExists {
		err = os.Mkdir(dir, 0777)
		if err != nil {
			return
		}
	}

	movies := strings.Split(text, " ") // [1] - movie

	if len(movies) > 1 {
		movie := movies[1]
		save(movie, m[movie])
		return
	}

	for ind, val := range m {
		save(ind, val)
	}

}

func (m Movies) loadMovie(text string) {
	directory := strings.Split(text, " d ")
	if len(directory) != 2 {
		fmt.Printf("Вы не ввели директорию %s\n", directory)
		log.Println("not enough arguments")
		return
	}

	entries, err := os.ReadDir("./" + directory[1])
	if err != nil {
		log.Println("can't read files from directory " + directory[1])
		return
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for _, e := range entries {
		log.Println("file " + e.Name())

		wg.Add(1)

		go func(e os.DirEntry) {
			defer wg.Done()
			file, err := os.ReadFile("./" + directory[1] + "/" + e.Name())
			if err != nil {
				log.Println("can't open file")
				return
			}
			name := strings.Split(e.Name(), ".")

			mu.Lock()
			m[name[0]] = append(m[name[0]], string(file))
			mu.Unlock()
		}(e)
		wg.Wait()
	}
}

func save(movie string, actors []string) {
	path := dir + "/" + movie + ".txt"

	exist, err := exists(path)
	if err != nil {
		return
	}

	if !exist {
		file, err := os.Create(path)
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

		log.Println(movie + " was saved")
	}
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
