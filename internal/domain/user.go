package domain

import _ "image"

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Document string `json:"document"`
	Email    string `json:"email"`
}
