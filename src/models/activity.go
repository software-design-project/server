package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Activity struct {
	Id          bson.ObjectId `bson:"_id"`
	SellerId    bson.ObjectId `bson:"sellerId"`
	Name        string        `bson:"name"`
	Slogan      string        `bson:"slogan"`
	Price       int           `bson:"price"`
	Time        int64         `bson:"time"`
	Place       string        `bson:"place"`
	Description string        `bson:"description"`
	Notice      string        `bson:"notice"`
	Telephone   string        `bson:"telephone"`
}
