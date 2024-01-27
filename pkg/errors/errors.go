package errors

import "fmt"

type MissingObjectError struct {
	Msg string
}

func (e MissingObjectError) Error() string {
	return fmt.Sprintf("object not found: %s", e.Msg)
}

type InvalidInputError struct {
	Msg string
}

func (e InvalidInputError) Error() string {
	return fmt.Sprintf("invalid input: %s", e.Msg)
}
