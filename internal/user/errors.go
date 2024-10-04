package user

import (
	"errors"
	"fmt"
)

var ErrFirstNameRequired = errors.New("first name required")
var ErrLastNameRequired = errors.New("last name required")

type ErrNotFound struct {
	ID uint64
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user with ID %d not found", e.ID)
}
