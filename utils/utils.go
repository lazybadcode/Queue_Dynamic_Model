package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

func ObjectIDFromHex(hex string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}

func BeginOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func RemoveFirstLayer(data []string) []string {
	for i, v := range data {
		l := strings.Split(v, ".")
		if len(l) > 1 {
			l = l[1:]
		}

		data[i] = strings.Join(l, ".")
	}
	return data
}
