package user

import (
	"context"
	"database/sql"
	_ "errors"
	"github.com/EstebanLescano/Gym/internal/domain"
	"log"
)

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		Get(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, firstName, lastName, email *string) error
	}
	repo struct {
		db  *sql.DB
		log *log.Logger
	}
)

func NewRepo(db *sql.DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(ctx context.Context, user *domain.User) error {
	sqlQ := "INSERT INTO users(name, last_name, email) VALUES(?,?,?)"
	res, err := r.db.Exec(sqlQ, user.Name, user.LastName, user.Email)
	if err != nil {
		r.log.Println(err.Error())
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		r.log.Println(err.Error())
		return err
	}
	user.ID = uint64(id)
	r.log.Printf("User created with id: %d", id)
	return nil
}

func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {
	var user []domain.User
	sqlQ := "SELECT id, name, last_name, email FROM users"
	rows, err := r.db.Query(sqlQ)
	if err != nil {
		r.log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u domain.User
		err := rows.Scan(&u.ID, &u.Name, &u.LastName, &u.Email)
		if err != nil {
			r.log.Println(err.Error())
			return nil, err
		}
		user = append(user, u)
	}
	r.log.Printf("Users retrieved with %d users", len(user))
	return nil, nil
}

func (r *repo) Get(ctx context.Context, id uint64) (*domain.User, error) {
	/*	index := slices.IndexFunc(r.db.Users, func(v domain.User) bool {
			return v.ID == id
		})
		if index < 0 {
			return nil, ErrNotFound{id}
		}*/
	return nil, nil
}

func (r *repo) Update(ctx context.Context, id uint64, firstName, lastName, email *string) error {
	/*	user, err := r.Get(ctx, id)
		if err != nil {
			return err
		}
		if firstName != nil {
			user.Name = *firstName
		}
		if lastName != nil {
			user.LastName = *lastName
		}
		if email != nil {
			user.Email = *email
		}*/
	return nil
}
