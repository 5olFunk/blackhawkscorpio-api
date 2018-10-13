package main

import (
	"encoding/json"
	_ "fmt"
	"io/ioutil"
	"log"
	"net/http" //dao "github.com/5olFunk/blackhawkscorpio-api/dao"
	//. "github.com/5olFunk/blackhawkscorpio-api/models"

	"github.com/gorilla/mux"
	//"gopkg.in/mgo.v2/bson"
)

var places []Place

type Place struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Lat     float64  `json:"lat"`
	Long    float64  `json:"long"`
	Ratings []Rating `json:"ratings"`
}

type Rating struct {
	Culture string   `json:"culture"`
	Score   float64  `json:"score"`
	Comment []string `json:"comments"`
}

type GoogleResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Geometry Geometry `json:"geometry"`
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Rating   float32  `json:"rating"`
}

type Geometry struct {
	Location Location `json:"location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func initPlaces() {
	places = append(places, Place{
		ID:   "1234567788",
		Name: "super cool place 1",
		Lat:  38.774349,
		Long: -90.166409,
		Ratings: []Rating{Rating{
			"US",
			6.7,
			[]string{
				"This place is great",
				"naah, bad hot wings",
			},
		}}})
}

func googleSearchify(phrase string) ([]Result, error) {
	here := "https://maps.googleapis.com/maps/api/place/nearbysearch/json?key=AIzaSyCOQ1mHzFff_OkGigI4RgZ6pOQbFufExnI&location=38.632069,-90.227531&radius=16093.4"
	resp, err := http.Get(here + "&keyword=" + phrase)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Failed to read body.")
		return nil, err
	}

	var tmp GoogleResponse
	if err := json.Unmarshal(body, &tmp); err != nil {
		log.Print("Failed to parse body to results")
		return nil, err
	}
	return tmp.Results, nil
}

func googleSearchifyRaw(phrase string) (string, error) {
	here := "https://maps.googleapis.com/maps/api/place/nearbysearch/json?key=AIzaSyCOQ1mHzFff_OkGigI4RgZ6pOQbFufExnI&location=38.632069,-90.227531&radius=16093.4"
	resp, err := http.Get(here + "&keyword=" + phrase)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// var tmp map[string]interface{}
	// if err := json.Unmarshal(body, &tmp); err != nil {
	// 	return "", err
	// }

	// log.Print(string(tmp["results"].(string)))

	return string(body), nil
}

// this will eventually take a search phrase and return
// results hydrated with our internal data
func SearchPlacesEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	searchResults, _ := googleSearchify(params["phrase"])
	respondWithJson(w, 200, searchResults)
	//var objMap = map[string]*json.RawMessage
	//results, err := json.Unmarshal([])

	//results := searchResults["results"]
	//log.Print(json.M

	// res, err := googleSearchifyRaw(params["phrase"])
	// if err != nil {
	// 	respondWithError(w, 400, "narf.")
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(200)
	// w.Write([]byte(res))

}

// func GetPlaceByIdEndpoint(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	place, err := dao.FindById(params["id"])
// 	if err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid Place ID")
// 		return
// 	}
// 	respondWithJson(w, http.StatusOK, place)
// }

// func CreatePlaceEndpoint(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var place Place
// 	if err := json.NewDecoder(r.Body).Decode(&place); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	place.ID = bson.NewObjectId()
// 	if err := dao.Insert(place); err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondWithJson(w, http.StatusCreated, place)
// }

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	r := mux.NewRouter()
	initPlaces()
	r.HandleFunc("/placesSearch/{phrase}", SearchPlacesEndpoint).Methods("GET")
	//r.HandleFunc("/places/{id}", GetPlaceByIdEndpoint).Methods("GET")
	//r.HandleFunc("/places", CreatePlaceEndpoint).Methods("POST")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
