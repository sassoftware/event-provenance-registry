// Custom implementation of JSON type because graph-gophers/graphql-go does
// not support the custom JSON scalar type
package types

import (
	"gorm.io/datatypes"
)

// datatype required to be used by gorm so we can support multiple
// database types for our JSON blobs
type JSON struct {
	datatypes.JSON
}

func (JSON) ImplementsGraphQLType(name string) bool {
	return name == "JSON"
}

func (j *JSON) UnmarshalGraphQL(input interface{}) error {
	return j.Scan(input)
}

