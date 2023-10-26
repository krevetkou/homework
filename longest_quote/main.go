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
	Date, Month, Year string
}

type Quote struct {
	Text   string
	Length int
}

type PersonData struct {
	Name      string
	BirthDate BirthDate
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

func (b *BirthDate) setBirthDate(data string) {
	s := strings.Split(data, ".")

	b.Date = s[0]
	b.Month = s[1]
	b.Year = s[2]

}

func (q *Quote) setQuote(data string) {
	q.Text = data
	q.Length = utf8.RuneCountInString(data)
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

	count, wCount, pCount := 0, 0, 0
	for scanner.Scan() {
		scan := scanner.Text()
		switch {
		case count%4 == 0:
			s := strings.Split(scan, ",")

			writers = append(writers, Writer{
				PersonData: PersonData{
					Name: s[0],

					Quote: Quote{},
				},
				Genre: s[2],
			})
			writers[wCount].BirthDate.setBirthDate(s[1])

			wCount++
		case (count-2)%4 == 0:
			s := strings.Split(scan, ",")

			politicians = append(politicians, Politician{
				PersonData: PersonData{
					Name:  s[0],
					Quote: Quote{},
				},
				JobTitle: s[2],
			})

			politicians[pCount].BirthDate.setBirthDate(s[1])

			pCount++
		case (count-1)%4 == 0 && wCount < count || count == 1:
			writers[wCount-1].Quote.setQuote(scan)
		case (count-3)%4 == 0 && pCount < count || count == 3:
			politicians[pCount-1].Quote.setQuote(scan)
		}
		count++
	}

	return writers, politicians
}
