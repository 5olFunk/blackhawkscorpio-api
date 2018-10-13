package main

import (
    "gopkg.in/mgo.v2/bson"
)

type Rating struct {
    ID bson.ObjectId `json:"id" bson:"_id"`
    Culture string `json:"culture"`
    Score float64 `json:"score" bson:"score"`
    Comment []string `json:"comments" bson:"comments"`
}