package main

import (
	"encoding/json"
	_ "fmt"
	"io/ioutil"
	"log"
	"net/http" //dao "github.com/5olFunk/blackhawkscorpio-api/dao"
	//. "github.com/5olFunk/blackhawkscorpio-api/models"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	//"gopkg.in/mgo.v2/bson"
)

var placesStore []Place

type Place struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Lat     float64     `json:"lat"`
	Long    float64     `json:"long"`
	Images  []string    `json:"images"`
	Ratings []Rating    `json:"ratings"`
	Blob    interface{} `json:"blob"`
}

type Rating struct {
	Culture  string   `json:"culture"`
	Score    float64  `json:"score"`
	Comments []string `json:"comments"`
}

type GoogleResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Geometry Geometry `json:"geometry"`
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Rating   float32  `json:"rating"`
	Photos   []Photo  `json:"photos"`
}

type Photo struct {
	HtmlAttributions []string `json:"html_attributions"`
}

type Geometry struct {
	Location Location `json:"location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func findById(xs []Place, id string) *Place {
	for _, x := range xs {
		if x.ID == id {
			return &x
		}
	}
	return nil
}

func googleSearchify(phrase string) ([]Result, error) {
	here := "https://maps.googleapis.com/maps/api/place/nearbysearch/json?key=AIzaSyBjY8lU-8KkRNYHP6fqCunBqYhDyGgdz0A&location=38.632069,-90.227531&radius=16093.4"
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

func unquoteAndUntag(str string) string {
	splitRes := strings.Split(str, "\"")

	if len(splitRes) >= 2 {
		return splitRes[1]
	} else {
		return ""
	}

}

func hydrateResults(googleResults []Result) []Place {
	var places []Place
	for _, res := range googleResults {

		match := findById(placesStore, res.ID)
		usrating := Rating{
			Culture: "US",
			Score:   float64(res.Rating),
		}

		var ratings []Rating
		var blob interface{}
		if match != nil {
			ratings = append(match.Ratings, usrating)
			blob = match.Blob
		} else {
			ratings = append(ratings, usrating)
		}

		var unquotedImages []string
		for _, x := range res.Photos {

			unquotedImages = append(
				unquotedImages,
				unquoteAndUntag(x.HtmlAttributions[0]))
		}

		places = append(places, Place{
			ID:      res.ID,
			Name:    res.Name,
			Lat:     res.Geometry.Location.Lat,
			Long:    res.Geometry.Location.Lng,
			Images:  unquotedImages,
			Ratings: ratings,
			Blob:    blob,
		})

	}
	return places
}

func GetPlaceById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var match Place
	for _, x := range placesStore {
		if x.ID == id {
			match = x
		}
	}
	respondWithJson(w, 200, match)
}

// this will eventually take a search phrase and return
// results hydrated with our internal data
func SearchPlacesEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	searchResults, _ := googleSearchify(params["phrase"])

	hydratedResults := hydrateResults(searchResults)
	respondWithJson(w, 200, hydratedResults)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

const port = ":3000"

func init() {
	jsonFile, err := os.Open("data.json")

	if err != nil {
		log.Print("Error opening data.json: " + err.Error())
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &placesStore)
	if err != nil {
		log.Print(err)
	}
}

func main() {
	r := mux.NewRouter()
	corsObj := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/placesSearch/{phrase}", SearchPlacesEndpoint).Methods("GET")
	r.HandleFunc("/places/{id}", GetPlaceById).Methods("GET")
	log.Print("Listening on localhost" + port)
	if err := http.ListenAndServe(port, handlers.CORS(corsObj)(r)); err != nil {
		log.Fatal(err)
	}
}
