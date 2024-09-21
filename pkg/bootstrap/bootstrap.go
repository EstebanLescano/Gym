package bootstrap

import (
	"Gym/internal/domain"
	"Gym/internal/user"
	"log"
	"os"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}

func NewDB() user.DB {
	return user.DB{
		Users: []domain.User{{
			ID:       1,
			Name:     "Nahuel",
			LastName: "Costamagna",
			Email:    "nahuel@test.com",
		}, {
			ID:       1,
			Name:     "Esteban",
			LastName: "Costama",
			Email:    "Esteban@test.com",
		}, {
			ID:       1,
			Name:     "Maribel",
			LastName: "Costa",
			Email:    "maribel@test.com",
		}},
		MaxUserID: 3,
	}
}
