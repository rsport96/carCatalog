package main

import (
	"encoding/json"
	"net/http"
)

type Car struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Owner  People `json:"owner"`
}

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

func main() {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		regNum := r.URL.Query().Get("regNum")
		if regNum == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		car := Car{
			RegNum: regNum,
			Mark:   "Lada",
			Model:  "Vesta",
			Year:   2002,
			Owner: People{
				Name:       "John",
				Surname:    "Doe",
				Patronymic: "Smith",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(car)
	})

	http.ListenAndServe(":8081", nil)
}
