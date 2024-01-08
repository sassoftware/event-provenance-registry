// Custom implementation of JSON type because graph-gophers/graphql-go does
// not support the custom JSON scalar type
package types

import (
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

// Time is a custom GraphQL type to represent an instant in time. It has to be added to a schema
// via "scalar Time" since it is not a predeclared GraphQL type like "ID".
type Time struct {
	datatypes.Date
}

// ImplementsGraphQLType maps this custom Go type
// to the graphql scalar type in the schema.
func (Time) ImplementsGraphQLType(name string) bool {
	return name == "Time"
}

// UnmarshalGraphQL is a custom unmarshaler for Time
//
// This function will be called whenever you use the
// time scalar as an input
func (t *Time) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case datatypes.Date:
		t.Date = input
		return nil
	default:
		return fmt.Errorf("wrong type for Time: %T", input)
	}
}

// MarshalJSON is a custom marshaler for Time
//
// This function will be called whenever you
// query for fields that use the Time type
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Date)
}
