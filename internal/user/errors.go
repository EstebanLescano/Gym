package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrFirstNameRequired = errors.New("first name required")
var ErrLastNameRequired = errors.New("last name required")
var ErrDocumentRequired = errors.New("document required")
var ErrEmailRequired = errors.New("email required")
var ErrThereArentFields = errors.New("there aren't any fields")

type ErrNotFound struct {
	ID uint64
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user with ID %d not found", e.ID)
}

func (e ErrNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e ErrNotFound) GetBody() ([]byte, error) {
	return json.Marshal(map[string]string{"error": e.Error()})
}

func (e ErrNotFound) GetData() interface{} {
	return nil // No hay datos adicionales que devolver
}
