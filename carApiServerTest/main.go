package main

import (
	"encoding/json"
	"net/http"
)

type Car struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int64  `json:"year,omitempty"`
	Owner  People `json:"owner"`
}

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

var cache = map[string]Car{
	"X123XX150": {
		RegNum: "X123XX150",
		Mark:   "Lada",
		Model:  "Vesta",
		Owner: People{
			Name:       "John",
			Surname:    "Doe",
			Patronymic: "Smith",
		},
	},
	"A456BC789": {
		RegNum: "A456BC789",
		Mark:   "Toyota",
		Model:  "Corolla",
		Year:   2015,
		Owner: People{
			Name:    "Alice",
			Surname: "Johnson",
		},
	},
	"H789GF123": {
		RegNum: "H789GF123",
		Mark:   "BMW",
		Model:  "X5",
		Year:   2019,
		Owner: People{
			Name:       "Bob",
			Surname:    "Brown",
			Patronymic: "Lee",
		},
	},
	"Z456BC789": {
		RegNum: "Z456BC789",
		Mark:   "Toyota",
		Model:  "Mark 2",
		Year:   1999,
		Owner: People{
			Name:    "Alice",
			Surname: "Johnson",
		},
	},
}


func main() {
	http.HandleFunc("/car", getCarByRegNum)
	http.ListenAndServe(":8086", nil)
}

func getCarByRegNum(w http.ResponseWriter, r *http.Request) {
	regNum := r.URL.Query().Get("regNum")
	car, found := cache[regNum]
	if !found {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(car)
	if err != nil {
		http.Error(w, "Failed to marshal car data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
