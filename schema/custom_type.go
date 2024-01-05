package schema

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
	"time"
)

// Define a JSONType as a Scalar type for dynamic JSON
var jsonType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "JSON",
	Description: "Arbitrary JSON data",
	// Define serialization and parsing methods for the JSON type
	Serialize: func(value interface{}) interface{} {
		log.Printf("Serialize called with: %#v", value)
		return value
	},
	ParseValue: func(value interface{}) interface{} {
		log.Printf("ParseValue called with: %#v", value)
		return value
	},
	ParseLiteral: ParseLiteral,
})

func ParseLiteral(valueAST ast.Value) interface{} {
	switch valueAST := valueAST.(type) {
	case *ast.ObjectValue:
		// Handle object values - typically used for JSON objects
		obj := make(map[string]interface{})
		for _, field := range valueAST.Fields {
			obj[field.Name.Value] = ParseLiteral(field.Value)
		}
		return obj
	case *ast.ListValue:
		// Handle list values - typically used for JSON arrays
		arr := make([]interface{}, len(valueAST.Values))
		for i, value := range valueAST.Values {
			arr[i] = ParseLiteral(value)
		}
		return arr
	case *ast.StringValue:
		return valueAST.Value
	case *ast.IntValue:
		// Convert string to int64
		if intValue, err := strconv.ParseInt(valueAST.Value, 10, 64); err == nil {
			return intValue
		}
		return nil // or handle the error as you see fit
	case *ast.FloatValue:
		// Convert string to float64
		if floatValue, err := strconv.ParseFloat(valueAST.Value, 64); err == nil {
			return floatValue
		}
		return nil
	case *ast.BooleanValue:
		return valueAST.Value
	default:
		return valueAST.GetValue()
	}
}

func newType(name string, model Model) *graphql.Object {
	fields := graphql.Fields{}

	data, ok := model[name]
	if !ok {
		panic("not found model " + name)
	}
	for k, v := range data {
		var t graphql.Type
		switch v.Type {
		case "Int":
			t = graphql.Int
		case "String":
			t = graphql.String
		case "Float":
			t = graphql.Float
		case "Boolean":
			t = graphql.Boolean
		case "DateTime":
			t = DateTime
		case "User":
			k = "user"
			t = newType("user", model)
		default:
			panic("not found type " + v.Type)
		}
		fields[k] = &graphql.Field{
			Type: t,
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: fields,
	})
}

func newInputType(name string, model Model) *graphql.InputObject {
	fields := graphql.InputObjectConfigFieldMap{}

	data, ok := model[name]
	if !ok {
		panic("not found model " + name)
	}
	for k, v := range data {
		var t graphql.Type
		switch v.Type {
		case "Int":
			t = graphql.Int
		case "String":
			t = graphql.String
		case "Float":
			t = graphql.Float
		case "Boolean":
			t = graphql.Boolean
		case "DateTime":
			t = DateTime
		case "User":
			t = newInputType("user", model)
		default:
			panic("not found type " + v.Type)
		}
		fields[k] = &graphql.InputObjectFieldConfig{
			Type: t,
		}
	}

	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   fmt.Sprintf("input%v", name),
		Fields: fields,
	})
}

var DateTime = graphql.NewScalar(graphql.ScalarConfig{
	Name: "DateTime",
	Description: "The `DateTime` scalar type represents a DateTime." +
		" The DateTime is serialized as an RFC 3339 quoted string",
	Serialize:  serializeDateTime,
	ParseValue: unserializeDateTime,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return unserializeDateTime(valueAST.Value)
		}
		return nil
	},
})

func serializeDateTime(value interface{}) interface{} {
	switch value := value.(type) {
	case time.Time:
		buff, err := value.MarshalText()
		if err != nil {
			return nil
		}

		return string(buff)
	case *time.Time:
		if value == nil {
			return nil
		}
		return serializeDateTime(*value)
	case primitive.DateTime:
		return value.Time()
	default:
		return nil
	}
}

func unserializeDateTime(value interface{}) interface{} {
	switch value := value.(type) {
	case []byte:
		t := time.Time{}
		err := t.UnmarshalText(value)
		if err != nil {
			return nil
		}

		return t
	case string:
		return unserializeDateTime([]byte(value))
	case *string:
		if value == nil {
			return nil
		}
		return unserializeDateTime([]byte(*value))
	case time.Time:
		return value
	default:
		return nil
	}
}
