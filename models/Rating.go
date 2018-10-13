package models

type Rating struct {
	Culture string   `json:"culture"`
	Score   float64  `json:"score"`
	Comment []string `json:"comments"`
}
