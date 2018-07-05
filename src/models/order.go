package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Order struct {
	Id       bson.ObjectId `bson:"_id"`
	UserId   bson.ObjectId `bson:"userId"`
	SellerId bson.ObjectId `bson:"sellerId"`
	State    int           `bson:"state"`
	Code     string        `bson:"code"`
	Time     int64         `bson:"time"`
	Place    string        `bson:"place"`
}
