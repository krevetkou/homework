package domain

// Характеристики актера: полное имя, год рождения, страна рождения, пол
type Actor struct {
	ID             int    `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	BirthYear      int    `json:"birth_year" db:"birth_year"`
	CountryOfBirth string `json:"country_of_birth" db:"country_of_birth"`
	Gender         string `json:"gender" db:"gender"`
}

type ActorUpdate struct {
	Name           *string `json:"name,omitempty"`
	BirthYear      *int    `json:"birth_year,omitempty"`
	CountryOfBirth *string `json:"country_of_birth,omitempty"`
	Sex            *string `json:"sex,omitempty"`
}

//{
//"name": "a",
//"birth_year": 1234,
//"country_of_birth": "afads",
//"gender": "fv"
//}
