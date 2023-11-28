package domain

// Характеристики актера: полное имя, год рождения, страна рождения, пол
type Actor struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BirthYear      string `json:"birth_year"`
	CountryOfBirth string `json:"country_of_birth"`
	Sex            string `json:"sex"`
}

type ActorUpdate struct {
	Name           *string `json:"name,omitempty"`
	BirthYear      *string `json:"birth_year,omitempty"`
	CountryOfBirth *string `json:"country_of_birth,omitempty"`
	Sex            *string `json:"sex,omitempty"`
}
