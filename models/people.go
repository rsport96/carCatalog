package models

type People struct {
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic,omitempty" db:"patronymic,omitempty"`
}
