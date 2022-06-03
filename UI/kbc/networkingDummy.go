package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const httpPort = 8080

type user struct {
	ID        int      `json:"id"`
	Moniker   string   `json:"moniker"`
	Bio       string   `json:"bio"`
	Languages []string `json:"languages`
}

var allUsers []*user

func init() {
	allUsers = []*user{
		{ID: 1, Moniker: "Hades", Bio: `god of the underworld, ruler of the dead
		and brother to the supreme ruler of the gods, Zeus`, Languages: []string{"Greek"}},
		{ID: 2, Moniker: "Horus", Bio: `god of the sun, sky and war`, Languages: []string{"Arabic"}},
		{ID: 3, Moniker: "Apollo", Bio: `god of light, music, manly beauty, dance, prophecy, medicine,
		poetry and almost every other thing. Son of Zeus`, Languages: []string{"Greek"}},
		{ID: 4, Moniker: "Artemis", Bio: `goddess of the wilderness and wild animals.
		Sister to Apollo and daughter of Zeus`, Languages: []string{"Greek"}},
	}
}

type users struct {
}

func (u users) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	j, _ := json.Marshal(allUsers)

	fmt.Fprintf(w, string(j))
}

func RunDummy() {
	http.Handle("/users/", users{})

	http.ListenAndServe(":8080", nil)
}
