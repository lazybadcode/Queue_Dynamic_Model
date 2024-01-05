package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"queue/usecase"
)

func NewSchemaHandler(usc usecase.UsecaseInterface, conf *Config) *handler.Handler {
	//userTypeV2 := newType("user", usc.Model())
	queueType := newType("queue", conf.Model)
	queueInputType := newInputType("queue", conf.Model)
	// Define the Query type
	var queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"queue": &graphql.Field{
				Type: graphql.NewList(queueType),
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(queueInputType),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					input := params.Args["input"].(map[string]interface{})
					return usc.GetQueue(params.Context, input)
				},
			},
		},
	})

	// Define the Mutation type
	var mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createQueue": &graphql.Field{
				Type: queueType, // The return type of the mutation
				Args: graphql.FieldConfigArgument{
					"idCard": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"mobileNo": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"input": &graphql.ArgumentConfig{
						//Type: jsonType, // Dynamic JSON data
						Type: graphql.NewNonNull(queueInputType),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					idCard := params.Args["idCard"].(string)
					mobileNo := params.Args["mobileNo"].(string)
					input := params.Args["input"].(map[string]interface{}) // input is a map

					return usc.CreateQueue(params.Context, idCard, mobileNo, input)
				},
			},

			"updateQueue": &graphql.Field{
				Type: queueType, // The return type of the mutation
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"date": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"slot": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(queueInputType),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(string)
					date := params.Args["date"].(string)
					slot := params.Args["slot"].(int)
					input := params.Args["input"].(map[string]interface{}) // input is a map

					return usc.UpdateQueue(params.Context, id, date, int(slot), input)
				},
			},

			"deleteQueue": &graphql.Field{
				Type: queueType, // The return type of the mutation
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(string)
					return usc.DeleteQueue(params.Context, id)
				},
			},
		},
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new schema, error: %v", err)
	}

	return handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		Playground: true,
	})
}
