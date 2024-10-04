package domain

import _ "image"

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Document uint64 `json:"document"`
	Email    string `json:"email"`
}
