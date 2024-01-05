package model

import (
	"encoding/json"
	"log"
)

func FromMap(m map[string]interface{}, out interface{}) {
	// Convert map to JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("Error marshalling map: %v", err)
	}

	// Unmarshal JSON data into the struct
	if err := json.Unmarshal(jsonData, out); err != nil {
		log.Fatalf("Error unmarshalling json: %v", err)
	}

}
