package main

import (
    "log"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
)

func GetPlaces(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func CreatePlace(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "not implemented yet !")
}

func main () {
    r := mux.NewRouter()
    r.HandleFunc("/placesSearch/{phrase}", GetPlaces).Methods("GET");
    if err := http.ListenAndServe(":3000", r); err != nil {
        log.Fatal(err)
    }
}
