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
	Date, Month, Year int
}

type Quote struct {
	Text   string
	Length int
}

type PersonData struct {
	Name      string
	BirthDate BirthDate
}

type Writer struct {
	PersonData
	Genre string
	Quote
}

type Politician struct {
	PersonData
	JobTitle string
	Quote
}

func main() {
	var writers [3]Writer
	var politicians [2]Politician

	count, wCount, pCount := 0, 0, 0

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
	for scanner.Scan() {
		scan := scanner.Text()
		switch {
		case count == 0 || count == 4 || count == 8:
			s := strings.Split(scan, ",")
			b := strings.Split(s[1], ".")
			writers[wCount].Name = s[0]
			writers[wCount].BirthDate.Date, _ = strconv.Atoi(b[0])
			writers[wCount].BirthDate.Month, _ = strconv.Atoi(b[1])
			writers[wCount].BirthDate.Year, _ = strconv.Atoi(b[2])
			writers[wCount].Genre = s[2]
			wCount++
		case count == 2 || count == 6:
			s := strings.Split(scan, ",")
			b := strings.Split(s[1], ".")
			politicians[pCount].Name = s[0]
			politicians[pCount].BirthDate.Date, _ = strconv.Atoi(b[0])
			politicians[pCount].BirthDate.Month, _ = strconv.Atoi(b[1])
			politicians[pCount].BirthDate.Year, _ = strconv.Atoi(b[2])
			politicians[pCount].JobTitle = s[2]
			pCount++
		case wCount < 4 && (count == 1 || count == 5 || count == 9):
			writers[wCount-1].Quote.Text = scan
			writers[wCount-1].Quote.Length = utf8.RuneCountInString(scan)
		case pCount < 3 && (count == 3 || count == 7):
			politicians[pCount-1].Quote.Text = scan
			politicians[pCount-1].Quote.Length = utf8.RuneCountInString(scan)
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		return
	}

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

	fmt.Println("Самая длинная цитата принадлежит " + author + "ее длина составляет " + strconv.Itoa(aMax))
}
