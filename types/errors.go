package types

import (
	"errors"
	"fmt"
)

type InvalidDeckError struct {
	Err error
}

func (r *InvalidDeckError) Error() string {
	return fmt.Sprintf("Invalid Deck: %s", r.Err)
}

func InvalidDeckErrorFromString(s string) *InvalidDeckError {
	return &InvalidDeckError{
		Err: errors.New(s),
	}
}
