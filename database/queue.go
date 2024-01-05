package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func (db *DB) FindOneQueue(ctx context.Context, id primitive.ObjectID) (map[string]interface{}, error) {
	var result map[string]interface{}
	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollQueue)

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (db *DB) FindAllQueue(ctx context.Context, filter map[string]interface{}, project []string) ([]map[string]interface{}, error) {
	var queues []map[string]interface{}
	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollQueue)

	var projection bson.D
	for _, s := range project {
		projection = append(projection, primitive.E{
			Key:   s,
			Value: 1,
		})
	}

	opts := options.Aggregate()

	// Aggregation pipeline
	pipeline := mongo.Pipeline{
		bson.D{{"$match", filter}},
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "userId"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$unwind", "$user"}}, // Converts user array to single object
	}
	if projection != nil {
		pipeline = append(pipeline, bson.D{{"$project", projection}})
	}

	cursor, err := collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &queues); err != nil {
		return nil, err
	}

	return queues, nil

}

func (db *DB) FindQueueInCreatedTimeRange(ctx context.Context, start, stop time.Time) ([]map[string]interface{}, error) {
	return db.FindAllQueue(ctx, bson.M{"created_time": bson.D{{"$gte", start}, {"$lt", stop}}}, nil)
}

type SlotCount struct {
	Slot  int `bson:"_id"`
	Count int `bson:"count"`
}

func (db *DB) ListSlotByDay(ctx context.Context, t time.Time) (results []SlotCount, maxNo int, err error) {
	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollQueue)

	opts := options.Aggregate()

	// Aggregation pipeline
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"date", bson.D{{"$eq", t.Format("20060102")}}}}}},
		{{"$group", bson.D{{"_id", "$slot"}, {"count", bson.D{{"$sum", 1}}}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, 0, err
	}

	//get max no of the day
	pipeline = mongo.Pipeline{
		{{"$match", bson.D{{"date", bson.D{{"$eq", t.Format("20060102")}}}}}},
		bson.D{{"$group", bson.D{{"_id", nil}, {"max_no", bson.D{{"$max", "$no"}}}}}},
	}

	// Perform the aggregation query
	cursor, err = collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// Fetch and display the result
	var results2 []struct {
		MaxNo int `bson:"max_no"`
	}
	if err = cursor.All(ctx, &results2); err != nil {
		return nil, 0, err
	}
	for _, result := range results2 {
		maxNo = result.MaxNo
	}

	return results, maxNo, nil
}
func (db *DB) InsertQueue(ctx context.Context, queue map[string]interface{}) (primitive.ObjectID, error) {
	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollQueue)
	insertResult, err := collection.InsertOne(ctx, queue)
	if err != nil {
		var writeException mongo.WriteException
		if errors.As(err, &writeException) {
			for _, err := range writeException.WriteErrors {
				if err.Code == 11000 || err.Code == 11001 {
					return primitive.NilObjectID, nil
				}
			}
		}
		return primitive.NilObjectID, err
	}
	log.Println("Inserted a single document: ", insertResult.InsertedID)
	return insertResult.InsertedID.(primitive.ObjectID), nil
}

func (db *DB) UpdateQueue(ctx context.Context, id primitive.ObjectID, updateQueue map[string]interface{}) error {
	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollQueue)

	if len(updateQueue) == 0 {
		return nil
	}
	updateQueue["updatedTime"] = time.Now()
	update := bson.M{
		"$set": updateQueue,
	}

	_, err := collection.UpdateByID(ctx, id, update)
	return err
}

func (db *DB) DeleteQueue(ctx context.Context, id primitive.ObjectID) error {
	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollQueue)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("queue not found")
	}
	return nil
}

//########### Achieve zone ###########

//func (db *DB) FindAllQueue(ctx context.Context, filter map[string]interface{}, project []string) ([]*model.Queue, error) {
//	var queues []*model.Queue
//	collection := db.client.Database(db.conf.TableName).Collection(db.conf.CollQueue)
//
//	var projection bson.D
//	for _, s := range project {
//		projection = append(projection, primitive.E{
//			Key:   s,
//			Value: 1,
//		})
//	}
//
//	opts := options.Find()
//	opts.SetProjection(projection)
//	cursor, err := collection.Find(ctx, filter, opts)
//	if err != nil {
//		return nil, err
//	}
//	defer cursor.Close(ctx)
//
//	for cursor.Next(ctx) {
//		var queue model.Queue
//		err := cursor.Decode(&queue)
//		if err != nil {
//			return nil, err
//		}
//		queues = append(queues, &queue)
//	}
//
//	if err := cursor.Err(); err != nil {
//		return nil, err
//	}
//
//	return queues, nil
//}
