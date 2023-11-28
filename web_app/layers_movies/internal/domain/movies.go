package domain

// Характеристики фильма: название фильма, дата выхода, страна, жанр, рейтинг от 1 до 5
type Movies struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	Country     string `json:"country"`
	Genre       string `json:"genre"`
	Rating      int    `json:"rating"`
}
