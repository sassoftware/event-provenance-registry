package errors

import (
	"errors"
	"fmt"
)

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

func SanitizeError(err error) error {
	if err == nil {
		return nil
	}

	isServerErr := false
	switch err.(type) {
	case MissingObjectError:
	case InvalidInputError:
	default:
		isServerErr = true
	}

	if isServerErr {
		// don't expose server internals
		err = errors.New("internal server error")
	}
	return err
}
