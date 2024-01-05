package schema_test

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"queue/config"
	"queue/schema"
	usecase "queue/usecase/mocks"
	"testing"
)

func loadConfig() *config.Conf {
	return config.New(".././config/local")
}
func TestQueueQuery(t *testing.T) {
	conf := loadConfig()

	tests := []struct {
		name          string
		queryString   string
		expected      []map[string]interface{}
		expectedError interface{}
	}{
		{
			name: "date queue",
			queryString: `
				query {
				  queue(
					input: {
					  date:"20231226"
					}
				  ) {
					date
					no
				  }
				}`,
			expected: []map[string]interface{}{{
				"date": "20231226",
				"no":   123,
			}},
		},
		{
			name: "date and user name",
			queryString: `
				query {
				  queue(
					input: {
					  date:"20231226"
					}
				  ) {
					user {
					  name
					}
				  }
				}`,
			expected: []map[string]interface{}{{
				"user": map[string]interface{}{
					"name": "User Name",
				},
			}},
		},
		{
			name: "no field in input",
			queryString: `
				query {
				  queue(
					input: {
					  date:"20231226"
                      noField: "data"
					}
				  ) {
					user {
					  name
					}
				  }
				}`,
			expected: nil,
			expectedError: `Argument "input" has invalid value {date: "20231226", noField: "data"}.
In field "noField": Unknown field.`,
		},
		{
			name: "no field in output",
			queryString: `
				query {
				  queue(
					input: {
					  date:"20231226"
					}
				  ) {
					user {
					  name
					  noField
					}
				  }
				}`,
			expected:      nil,
			expectedError: `Cannot query field "noField" on type "user".`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := []map[string]interface{}{{
				"_id":    "someId",
				"userId": "someUserId",
				"no":     int(123),
				"date":   "20231226",
				"user": map[string]interface{}{
					"_id":  "someUserId",
					"name": "User Name",
				},
			},
			}
			mockUsc := usecase.NewMockUsecaseInterface(t)
			if tt.expectedError == nil {
				mockUsc.EXPECT().GetQueue(mock.Anything, mock.Anything).Return(output, nil)
			}

			schemaHandler := schema.NewSchemaHandler(mockUsc, &schema.Config{Model: conf.Schema.Model})

			req, _ := http.NewRequest("GET", fmt.Sprintf("/graphql?query=%s", url.QueryEscape(tt.queryString)), nil)

			h := handler.New(&handler.Config{
				Schema: schemaHandler.Schema,
				Pretty: false,
			})

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			result := rr.Result()

			assert.Equal(t, http.StatusOK, result.StatusCode)

			bodyBytes, err := io.ReadAll(result.Body)
			if err != nil {
				t.Fatalf("failed reading response body: %v", err)
			}
			var actual map[string]interface{}
			err = json.Unmarshal(bodyBytes, &actual)
			if err != nil {
				t.Fatalf("failed unmarshalling response body: %v", err)
			}

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, actual["errors"].([]interface{})[0].(map[string]interface{})["message"])
				return
			}

			if tt.expected != nil {
				if (actual["data"] == nil) != (tt.expected == nil) {
					t.Fatalf("expect %+v, actual %+v", tt.expected, actual)
					return
				}
				for i, v := range actual["data"].(map[string]interface{})["queue"].([]interface{}) {
					//assert.Equal(t, tt.expected[i], v.(map[string]interface{}))
					act, _ := json.Marshal(v)
					exp, _ := json.Marshal(tt.expected[i])
					assert.Equal(t, string(exp), string(act))
				}
			}

		})
	}
}

func TestCreateQueueMutation(t *testing.T) {
	conf := loadConfig()

	tests := []struct {
		name          string
		queryString   string
		expected      map[string]interface{}
		expectedError interface{}
	}{
		{
			name: "create queue",
			queryString: `
				mutation {
				  createQueue(
					idCard: "12345678910"
					mobileNo: "0801234567"
					input: {
					  note: "This is my note"
					  testBoolean:true
					  testFloat:8.88
					  user: { name: "this is my name" }
					}
				  ) {
					date
					no
					testFloat
					testBoolean
				  }
				}`,
			expected: map[string]interface{}{
				"date":        "20231226",
				"no":          123,
				"testFloat":   8.88,
				"testBoolean": true,
			},
		},
		{
			name: "create queue return queue and user data",
			queryString: `
				mutation {
				  createQueue(
					idCard: "12345678910"
					mobileNo: "0801234567"
					input: {
					  note: "This is my note"
					  user: { name: "this is my name" }
					}
				  ) {
					date
					no
					user {
					  name
					}
				  }
				}`,
			expected: map[string]interface{}{
				"date": "20231226",
				"no":   123,
				"user": map[string]interface{}{
					"name": "this is my name",
				},
			},
		},
		{
			name: "no field in input",
			queryString: `
				mutation {
				  createQueue(
					idCard: "12345678910"
					mobileNo: "0801234567"
					input: {
					  note: "This is my note"
					  user: { name: "this is my name" }
					  noField: "abc"
					}
				  ) {
					date
					no
				  }
				}`,
			expected: nil,
			expectedError: `Argument "input" has invalid value {note: "This is my note", user: {name: "this is my name"}, noField: "abc"}.
In field "noField": Unknown field.`,
		},
		{
			name: "no field in output",
			queryString: `
				mutation {
				  createQueue(
					idCard: "12345678910"
					mobileNo: "0801234567"
					input: {
					  note: "This is my note"
					  user: { name: "this is my name" }
					}
				  ) {
					date
					no
					noField
				  }
				}`,
			expected:      nil,
			expectedError: `Cannot query field "noField" on type "queue".`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := map[string]interface{}{
				"_id":         "someId",
				"userId":      "someUserId",
				"no":          int(123),
				"date":        "20231226",
				"testBoolean": true,
				"testFloat":   8.88,
				"user": map[string]interface{}{
					"_id":  "someUserId",
					"name": "this is my name",
				},
			}
			mockUsc := usecase.NewMockUsecaseInterface(t)
			if tt.expectedError == nil {
				mockUsc.EXPECT().CreateQueue(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(output, nil)
			}

			schemaHandler := schema.NewSchemaHandler(mockUsc, &schema.Config{Model: conf.Schema.Model})

			req, _ := http.NewRequest("GET", fmt.Sprintf("/graphql?query=%s", url.QueryEscape(tt.queryString)), nil)

			h := handler.New(&handler.Config{
				Schema: schemaHandler.Schema,
				Pretty: false,
			})

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			result := rr.Result()

			assert.Equal(t, http.StatusOK, result.StatusCode)

			bodyBytes, err := io.ReadAll(result.Body)
			if err != nil {
				t.Fatalf("failed reading response body: %v", err)
			}
			var actual map[string]interface{}
			err = json.Unmarshal(bodyBytes, &actual)
			if err != nil {
				t.Fatalf("failed unmarshalling response body: %v", err)
			}

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, actual["errors"].([]interface{})[0].(map[string]interface{})["message"])
				return
			}

			if tt.expected != nil {
				if (actual["data"] == nil) != (tt.expected == nil) {
					t.Fatalf("expect %+v, actual %+v", tt.expected, actual)
					return
				}
				v := actual["data"].(map[string]interface{})["createQueue"]
				act, _ := json.Marshal(v)
				exp, _ := json.Marshal(tt.expected)
				assert.Equal(t, string(exp), string(act))
			}

		})
	}
}

func TestUpdateQueueMutation(t *testing.T) {
	conf := loadConfig()

	tests := []struct {
		name          string
		queryString   string
		expected      map[string]interface{}
		expectedError interface{}
	}{
		{
			name: "update queue",
			queryString: `
				mutation {
				  updateQueue(
					id: "657d6f1676863b9c94c22242"
					date: "20231226"
					slot: 3
					input: {
					  note: "This is my note"
					  user: { name: "this is my name"}
					}
				  ) {
					date
					no
				  }
				}`,
			expected: map[string]interface{}{
				"date": "20231226",
				"no":   123,
			},
		},
		{
			name: "update queue return queue and user data",
			queryString: `
				mutation {
				  updateQueue(
					id: "657d6f1676863b9c94c22242"
					date: "20231226"
					slot: 3
					input: {
					  note: "This is my note"
					  user: { name: "this is my name"}
					}
				  ) {
					date
					no
					user {
					  name
					}
				  }
				}`,
			expected: map[string]interface{}{
				"date": "20231226",
				"no":   123,
				"user": map[string]interface{}{
					"name": "this is my name",
				},
			},
		},
		{
			name: "no field in input",
			queryString: `
				mutation {
				  updateQueue(
					id: "657d6f1676863b9c94c22242"
					date: "20231226"
					slot: 3
					input: {
					  note: "This is my note"
					  user: { name: "this is my name"}
					  noField: "abc"
					}
				  ) {
					date
					no
					user {
					  name
					}
				  }
				}`,

			expected: nil,
			expectedError: `Argument "input" has invalid value {note: "This is my note", user: {name: "this is my name"}, noField: "abc"}.
In field "noField": Unknown field.`,
		},
		{
			name: "no field in output",
			queryString: `
				mutation {
				  updateQueue(
					id: "657d6f1676863b9c94c22242"
					date: "20231226"
					slot: 3
					input: {
					  note: "This is my note"
					  user: { name: "this is my name"}
					}
				  ) {
					date
					no
				    noField
					user {
					  name
					}
				  }
				}`,
			expected:      nil,
			expectedError: `Cannot query field "noField" on type "queue".`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := map[string]interface{}{
				"_id":    "someId",
				"userId": "someUserId",
				"no":     int(123),
				"date":   "20231226",
				"user": map[string]interface{}{
					"_id":  "someUserId",
					"name": "this is my name",
				},
			}
			mockUsc := usecase.NewMockUsecaseInterface(t)
			if tt.expectedError == nil {
				mockUsc.EXPECT().UpdateQueue(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(output, nil)
			}

			schemaHandler := schema.NewSchemaHandler(mockUsc, &schema.Config{Model: conf.Schema.Model})

			req, _ := http.NewRequest("GET", fmt.Sprintf("/graphql?query=%s", url.QueryEscape(tt.queryString)), nil)

			h := handler.New(&handler.Config{
				Schema: schemaHandler.Schema,
				Pretty: false,
			})

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			result := rr.Result()

			assert.Equal(t, http.StatusOK, result.StatusCode)

			bodyBytes, err := io.ReadAll(result.Body)
			if err != nil {
				t.Fatalf("failed reading response body: %v", err)
			}
			var actual map[string]interface{}
			err = json.Unmarshal(bodyBytes, &actual)
			if err != nil {
				t.Fatalf("failed unmarshalling response body: %v", err)
			}

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, actual["errors"].([]interface{})[0].(map[string]interface{})["message"])
				return
			}

			if tt.expected != nil {
				if (actual["data"] == nil) != (tt.expected == nil) {
					t.Fatalf("expect %+v, actual %+v", tt.expected, actual)
					return
				}
				v := actual["data"].(map[string]interface{})["updateQueue"]
				act, _ := json.Marshal(v)
				exp, _ := json.Marshal(tt.expected)
				assert.Equal(t, string(exp), string(act))
			}

		})
	}
}

func TestDeleteQueueMutation(t *testing.T) {
	conf := loadConfig()

	tests := []struct {
		name          string
		queryString   string
		expected      map[string]interface{}
		expectedError interface{}
	}{
		{
			name: "delete queue",
			queryString: `
				mutation {
				  deleteQueue(
					id: "658aea00f9b1f054e7333519"
				  ) {
					date
					no
				  }
				}`,
			expected: map[string]interface{}{
				"date": "20231226",
				"no":   123,
			},
		},
		{
			name: "delete queue return queue and user data",
			queryString: `
				mutation {
				  deleteQueue(
					id: "658aea00f9b1f054e7333519"
				  ) {
					date
					no
					user {
					  name
					}
				  }
				}`,
			expected: map[string]interface{}{
				"date": "20231226",
				"no":   123,
				"user": map[string]interface{}{
					"name": "this is my name",
				},
			},
		},
		{
			name: "no field in output",
			queryString: `
				mutation {
				  deleteQueue(
					id: "658aea00f9b1f054e7333519"
				  ) {
					date
					no
				    noField
					user {
					  name
					}
				  }
				}`,
			expected:      nil,
			expectedError: `Cannot query field "noField" on type "queue".`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := map[string]interface{}{
				"_id":    "someId",
				"userId": "someUserId",
				"no":     int(123),
				"date":   "20231226",
				"user": map[string]interface{}{
					"_id":  "someUserId",
					"name": "this is my name",
				},
			}
			mockUsc := usecase.NewMockUsecaseInterface(t)
			if tt.expectedError == nil {
				mockUsc.EXPECT().DeleteQueue(mock.Anything, mock.Anything).Return(output, nil)
			}

			schemaHandler := schema.NewSchemaHandler(mockUsc, &schema.Config{Model: conf.Schema.Model})

			req, _ := http.NewRequest("GET", fmt.Sprintf("/graphql?query=%s", url.QueryEscape(tt.queryString)), nil)

			h := handler.New(&handler.Config{
				Schema: schemaHandler.Schema,
				Pretty: false,
			})

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			result := rr.Result()

			assert.Equal(t, http.StatusOK, result.StatusCode)

			bodyBytes, err := io.ReadAll(result.Body)
			if err != nil {
				t.Fatalf("failed reading response body: %v", err)
			}
			var actual map[string]interface{}
			err = json.Unmarshal(bodyBytes, &actual)
			if err != nil {
				t.Fatalf("failed unmarshalling response body: %v", err)
			}

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, actual["errors"].([]interface{})[0].(map[string]interface{})["message"])
				return
			}

			if tt.expected != nil {
				if (actual["data"] == nil) != (tt.expected == nil) {
					t.Fatalf("expect %+v, actual %+v", tt.expected, actual)
					return
				}
				v := actual["data"].(map[string]interface{})["deleteQueue"]
				act, _ := json.Marshal(v)
				exp, _ := json.Marshal(tt.expected)
				assert.Equal(t, string(exp), string(act))
			}

		})
	}
}
