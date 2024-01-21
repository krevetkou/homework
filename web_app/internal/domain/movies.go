package domain

import "time"

type Movie struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	ReleaseDate time.Time `json:"release_date"`
	Country     string    `json:"country"`
	Genre       string    `json:"genre"`
	Rating      int8      `json:"rating"`
}

type MovieUpdate struct {
	Name        *string    `json:"name,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Country     *string    `json:"country,omitempty"`
	Genre       *string    `json:"genre,omitempty"`
	Rating      *int8      `json:"rating,omitempty"`
}

//{
//"name": "a",
//"release_date": "2021-02-18T21:54:42.123Z",
//"country": "afads",
//"genre": "fv",
//"rating": 5
//}
