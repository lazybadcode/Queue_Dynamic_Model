package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	IDCard      string             `bson:"idCard" json:"idCard"`
	MobileNo    string             `bson:"mobileNo" json:"mobileNo"`
	UpdatedTime time.Time          `bson:"updatedTime" json:"updatedTime"`
	CreatedTime time.Time          `bson:"createdTime" json:"createdTime"`
}
