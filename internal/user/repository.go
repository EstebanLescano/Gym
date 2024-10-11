package user

import (
	"context"
	"database/sql"
	"errors"
	_ "errors"
	"fmt"
	"github.com/EstebanLescano/Gym/internal/domain"
	"log"
	"strings"
)

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		Get(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, firstName, lastName, document, email *string) error
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
	sqlQ := "INSERT INTO users(name, last_name,document, email) VALUES(?,?,?,?)"
	res, err := r.db.Exec(sqlQ, user.Name, user.LastName, user.Document, user.Email)
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
	sqlQ := "SELECT id, name, last_name, document, email FROM users"

	// Ejecuta la consulta con el contexto
	rows, err := r.db.QueryContext(ctx, sqlQ) // Aquí es importante pasar el contexto
	if err != nil {
		r.log.Println("Error executing query:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.LastName, &u.Document, &u.Email); err != nil {
			r.log.Println("Error scanning row:", err.Error())
			return nil, err
		}
		users = append(users, u) // Agrega usuarios a la lista
	}
	// Manejar el caso donde no se encuentren usuarios
	if len(users) == 0 {
		r.log.Println("No users found")
	}
	return users, nil
}

func (r *repo) Get(ctx context.Context, id uint64) (*domain.User, error) {
	sqlQ := "SELECT id, name, last_name, document, email FROM users WHERE id=?"
	var u domain.User
	if err := r.db.QueryRow(sqlQ, id).Scan(&u.ID, &u.Name, &u.LastName, &u.Document, &u.Email); err != nil {
		r.log.Println(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound{id}
		}
		return nil, err
	}
	r.log.Println("User retrieved with id:", id)
	return &u, nil
}

func (r *repo) Update(ctx context.Context, id uint64, Name, lastName, document, email *string) error {
	var fields []string
	var values []interface{}

	if Name != nil {
		fields = append(fields, "name=?")
		values = append(values, *Name)
	}
	if lastName != nil {
		fields = append(fields, "last_name=?")
		values = append(values, *lastName)
	}
	if document != nil {
		fields = append(fields, "document=?")
		values = append(values, *document)
	}
	if email != nil {
		fields = append(fields, "email=?")
		values = append(values, *email)
	}
	if len(fields) == 0 {
		r.log.Println(ErrThereArentFields.Error())
		return ErrThereArentFields
	}
	values = append(values, id)

	sqlQ := fmt.Sprintf("UPDATE users SET %s WHERE id=?", strings.Join(fields, ","))
	res, err := r.db.Exec(sqlQ, values...)
	if err != nil {
		r.log.Println(err.Error())
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		r.log.Println(err.Error())
		return err
	}
	if row == 0 {
		err := ErrNotFound{id}
		r.log.Println(err.Error())
		return err
	}
	r.log.Println("User updated with id:", id)
	return nil
}
