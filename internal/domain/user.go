package domain

import "image"

type User struct {
	ID       uint64 `json:"id"`
	Avatar   image.NRGBA64
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Document uint64 `json:"document"`
	Email    string `json:"email"`
}
