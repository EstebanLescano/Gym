package user

import (
	"Gym/internal/domain"
	"context"
	"log"
)

// se trabaja con una interface y con una struct para que despues sea mas facil poder usar esta capa en los test
// forma de trabajarlo para tener todx testeado
type (
	Service interface {
		Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error)
		GetAll(ctx context.Context) ([]domain.User, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error) {

	user := &domain.User{
		Name:     firstName,
		LastName: lastName,
		Email:    email,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s service) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	s.log.Println("obtension de users")
	return users, err
}
