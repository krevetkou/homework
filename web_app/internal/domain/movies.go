package domain

type Movie struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	ReleaseDate string `json:"release_date" db:"release_date"`
	Country     string `json:"country" db:"country"`
	Genre       string `json:"genre" db:"genre"`
	Rating      int    `json:"rating" db:"rating"`
}

type MovieUpdate struct {
	Name        *string `json:"name,omitempty"`
	ReleaseDate *string `json:"release_date,omitempty"`
	Country     *string `json:"country,omitempty"`
	Genre       *string `json:"genre,omitempty"`
	Rating      *int    `json:"rating,omitempty"`
}

//{
//"name": "a",
//"release_date": "2021-02-18T21:54:42.123Z",
//"country": "afads",
//"genre": "fv",
//"rating": 5
//}
