package main

import (
    "gopkg.in/mgo.v2/bson"
)

type Place struct {
    ID bson.ObjectId `json:"id" bson:"_id"`
    Lat float64 `json:"lat" bson:"lat"`
    Long float64 `json:"long" bson:"long"`
    Ratings []Rating `json:"ratings" bson:"ratings"`
}
