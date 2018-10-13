package models

type Place struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Lat     float64  `json:"lat"`
	Long    float64  `json:"long"`
	Ratings []Rating `json:"ratings"`
}
