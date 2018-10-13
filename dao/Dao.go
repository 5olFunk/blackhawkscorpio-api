package dao

import (
	"log"

	. "github.com/5olFunk/blackhawkscorpio-api/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PlacesDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "places"
)

// Establish a connection to database
func (m *PlacesDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of places
func (m *PlacesDAO) FindAll() ([]Place, error) {
	var places []Place
	err := db.C(COLLECTION).Find(bson.M{}).All(&places)
	return places, err
}

// Find a place by its id
func (m *PlacesDAO) FindById(id string) (Place, error) {
	var place Place
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&place)
	return place, err
}

// Insert a place into database
func (m *PlacesDAO) Insert(place Place) error {
	err := db.C(COLLECTION).Insert(&place)
	return err
}

// Delete an existing place
func (m *PlacesDAO) Delete(place Place) error {
	err := db.C(COLLECTION).Remove(&place)
	return err
}

// Update an existing place
func (m *PlacesDAO) Update(place Place) error {
	err := db.C(COLLECTION).UpdateId(place.ID, &place)
	return err
}
