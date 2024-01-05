package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Queue struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID      primitive.ObjectID `bson:"userId" json:"UserId"`
	No          int                `bson:"no" json:"no"`
	Date        string             `bson:"date" json:"date"`
	Slot        int                `bson:"slot" json:"slot"`
	Note        string             `bson:"note" json:"note"`
	UpdatedTime time.Time          `bson:"updatedTime" json:"updatedTime"`
	CreatedTime time.Time          `bson:"createdTime" json:"created_time"`
	User        *User              `bson:"user" json:"user"`
}
