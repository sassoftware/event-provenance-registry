// Custom implementation of JSON type because graph-gophers/graphql-go does
// not support the custom JSON scalar type
package types

import (
	"gorm.io/datatypes"
)

// JSON type is a wrapper around the datatypes.JSON type in Go.
// @property  - The code snippet defines a struct named `JSON` which embeds another struct named
// `datatypes.JSON`. This means that the `JSON` struct inherits all the fields and methods of the
// `datatypes.JSON` struct.
// datatype required to be used by gorm so we can support multiple
// database types for our JSON blobs
type JSON struct {
	datatypes.JSON
}

// ImplementsGraphQLType function is a method of the `JSON` struct. It is used to determine if
// the struct implements a specific GraphQL type. In this case, the function checks if the input `name`
// is equal to "JSON". If it is, the function returns `true`, indicating that the struct implements the
// GraphQL type "JSON".
func (JSON) ImplementsGraphQLType(name string) bool {
	return name == "JSON"
}

// UnmarshalGraphQL function is a method of the `JSON` struct. It is used to convert a GraphQL
// input value into a Go value. In this case, the function takes an `input` of type `interface{}` and
// attempts to scan it into the `JSON` struct using the `Scan` method.
func (j *JSON) UnmarshalGraphQL(input interface{}) error {
	return j.Scan(input)
}
