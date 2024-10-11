package user

import (
	"context"
	"github.com/EstebanLescano/Gym/internal/domain"
	"log"
)

// se trabaja con una interface y con una struct para que despues sea mas facil poder usar esta capa en los test
// forma de trabajarlo para tener todx testeado
type (
	Service interface {
		Create(ctx context.Context, name, lastName, document, email string) (*domain.User, error)
		GetAll(ctx context.Context) ([]domain.User, error)
		Get(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, name, lastName, document, email *string) error
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

func (s service) Create(ctx context.Context, name, lastName, document, email string) (*domain.User, error) {
	user := &domain.User{
		Name:     name,
		LastName: lastName,
		Document: document,
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
	return users, err
}

func (s service) Get(ctx context.Context, id uint64) (*domain.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Println("Error retrieving user:", err)
		return nil, err
	}
	return user, nil
}

func (s service) Update(ctx context.Context, id uint64, firstName, lastName, document, email *string) error {
	if err := s.repo.Update(ctx, id, firstName, lastName, document, email); err != nil {
		return err
	}
	return nil
}
