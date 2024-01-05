package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"queue/config"
	"queue/database"
	"queue/schema"
	"queue/usecase"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	env := os.Getenv("ENV")
	conf := config.New(fmt.Sprintf("./config/%s", env))

	db := database.New(&conf.DB)
	usc := usecase.New(db, &conf.Usecase)
	go usc.Batch()

	h := schema.NewSchemaHandler(usc, &conf.Schema)
	http.Handle("/graphql", h)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
