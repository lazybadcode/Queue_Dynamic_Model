package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Config struct {
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	URL            string `mapstructure:"url"`
	PORT           string `mapstructure:"port"`
	TableName      string `mapstructure:"table_name"`
	CollUser       string `mapstructure:"coll_user"`
	CollQueue      string `mapstructure:"coll_queue"`
	CollSpecialDay string `mapstructure:"coll_special_day"`
}

type DB struct {
	client *mongo.Client
	conf   *Config
}

func New(conf *Config) *DB {
	//loggerOptions := options.
	//	Logger().
	//	SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	//clientOptions := options.Client().ApplyURI("mongodb://user:password@localhost:27017/")
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", conf.Username, conf.Password, conf.URL, conf.PORT))
	//clientOptions.SetLoggerOptions(loggerOptions)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	return &DB{client: client, conf: conf}
}
