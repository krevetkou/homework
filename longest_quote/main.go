package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type BirthDate struct {
	Date, Month, Year uint
}

type Quote struct {
	Text   string
	Length int
}

type PersonData struct {
	Name string
	BirthDate
	Quote
}

type Writer struct {
	PersonData
	Genre string
}

type Politician struct {
	PersonData
	JobTitle string
}

func main() {
	file, err := os.Open("quotes.txt")
	if err != nil {
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	scanner := bufio.NewScanner(file)
	writers, politicians := vars(scanner)

	if err := scanner.Err(); err != nil {
		return
	}

	maximum, author := max(writers, politicians)

	fmt.Println("Самая длинная цитата принадлежит " + author + "ее длина составляет " + strconv.Itoa(maximum))
}

func (p *PersonData) setBirthDate(data string) {
	s := strings.Split(data, ".")

	if len(s) == 3 {
		x, _ := strconv.Atoi(s[0])
		y, _ := strconv.Atoi(s[1])
		z, _ := strconv.Atoi(s[2])

		p.Date = uint(x)
		p.Month = uint(y)
		p.Year = uint(z)
	} else {
		p.Date = 0
		p.Month = 0
		p.Year = 0
	}
}

func (p *PersonData) setQuote(data string) {
	p.Text = data
	p.Length = utf8.RuneCountInString(data)
}

func max(writers []Writer, politicians []Politician) (int, string) {
	aMax := writers[0].Quote.Length
	author := writers[0].Name

	for _, val := range writers {
		if val.Quote.Length > aMax {
			aMax = val.Quote.Length
			author = val.Name
		}
	}

	for _, val := range politicians {
		if val.Quote.Length > aMax {
			aMax = val.Quote.Length
			author = val.Name
		}
	}

	return aMax, author
}

func vars(scanner *bufio.Scanner) ([]Writer, []Politician) {
	var writers []Writer
	var politicians []Politician
	isWriter := true

	count, wCount, pCount := 0, 0, 0
	for scanner.Scan() {
		scan := scanner.Text()
		switch {
		case isWriter == true && count%2 == 0:
			s := strings.Split(scan, ",")

			writers = append(writers, Writer{
				PersonData: PersonData{
					Name: s[0],
				},
				Genre: s[2],
			})
			writers[wCount].setBirthDate(s[1])

			wCount++

		case isWriter == false && count%2 == 0:
			s := strings.Split(scan, ",")

			politicians = append(politicians, Politician{
				PersonData: PersonData{
					Name: s[0],
				},
				JobTitle: s[2],
			})

			politicians[pCount].setBirthDate(s[1])

			pCount++
		case isWriter == true && wCount < count || count == 1:
			writers[wCount-1].setQuote(scan)
			isWriter = !isWriter
		case isWriter == false && pCount < count || count == 3:
			politicians[pCount-1].setQuote(scan)
			isWriter = !isWriter
		}
		count++
	}

	return writers, politicians
}
