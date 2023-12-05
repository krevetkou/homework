package domain

// Характеристики фильма: название фильма, дата выхода, страна, жанр, рейтинг от 1 до 5
type Movie struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	ReleaseDate ReleaseDate `json:"release_date"`
	Country     string      `json:"country"`
	Genre       string      `json:"genre"`
	Rating      int         `json:"rating"`
}

type ReleaseDate struct {
	Date  int `json:"date"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type MovieUpdate struct {
	Name        *string      `json:"name,omitempty"`
	ReleaseDate *ReleaseDate `json:"release_date,omitempty"`
	Country     *string      `json:"country,omitempty"`
	Genre       *string      `json:"genre,omitempty"`
	Rating      *int         `json:"rating,omitempty"`
}

//{
//"name": "a",
//"release_date": {
//"date": 1,
//"month": 2,
//"year": 1905
//},
//"country": "afads",
//"genre": "fv",
//"rating": 5
//}
