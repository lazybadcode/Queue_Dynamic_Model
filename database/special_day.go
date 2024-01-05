package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"queue/model"
	"queue/utils"
	"time"
)

func (db *DB) QuerySpecialDays(ctx context.Context, t time.Time) (*model.SpecialDay, error) {

	var result *model.SpecialDay
	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollSpecialDay)

	err := collection.FindOne(ctx, bson.M{"date": t.Format("20060102")}).Decode(&result) //TODO maybe filter not work
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (db *DB) QuerySpecialDaysInRange(ctx context.Context, start, end time.Time) ([]model.SpecialDay, error) {
	specialDayCollection := db.client.Database(db.conf.TableName).Collection(db.conf.CollQueue)

	filter := bson.D{
		{"date", bson.D{
			{"$gte", utils.BeginOfDay(start)},
			{"$lte", utils.BeginOfDay(end)},
		}},
	}

	cursor, err := specialDayCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var specialDays []model.SpecialDay
	if err = cursor.All(ctx, &specialDays); err != nil {
		return nil, err
	}

	return specialDays, nil
}
