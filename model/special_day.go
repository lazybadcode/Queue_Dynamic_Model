package model

type SpecialDay struct {
	Date        string `bson:"date" json:"date"`
	Description string `bson:"description" json:"description"`
}
