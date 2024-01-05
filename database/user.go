package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"queue/model"
	"time"
)

func (db *DB) FindOneUserByIdCard(ctx context.Context, idCard string) (*model.User, error) {
	var result *model.User
	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollUser)

	err := collection.FindOne(ctx, bson.M{"id_card": idCard}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil

}
func (db *DB) FindOneUser(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var result *model.User
	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollUser)

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (db *DB) InsertUser(ctx context.Context, user map[string]interface{}) (primitive.ObjectID, error) {
	userCollection := db.client.Database(db.conf.TableName).Collection("users")
	result, err := userCollection.InsertOne(ctx, user)
	return result.InsertedID.(primitive.ObjectID), err
}

func (db *DB) UpdateUser(ctx context.Context, id primitive.ObjectID, updateUser map[string]interface{}) error {
	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollUser)

	if len(updateUser) == 0 {
		return nil
	}
	updateUser["updatedTime"] = time.Now()
	update := bson.M{
		"$set": updateUser,
	}

	_, err := collection.UpdateByID(ctx, id, update)
	return err
}
