package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Seller struct {
	Id        bson.ObjectId `bson:"_id"`
	Name      string        `bson:"name"`
	Telephone string        `bson:"telephone"`
}
