package domain

// Характеристики актера: полное имя, год рождения, страна рождения, пол
type Actor struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BirthYear      int    `json:"birth_year"`
	CountryOfBirth string `json:"country_of_birth"`
	Gender         string `json:"gender"`
}

type ActorUpdate struct {
	Name           *string `json:"name,omitempty"`
	BirthYear      *int    `json:"birth_year,omitempty"`
	CountryOfBirth *string `json:"country_of_birth,omitempty"`
	Sex            *string `json:"sex,omitempty"`
}
