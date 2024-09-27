package user

import (
	"context"
	"errors"
	"fmt"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoint struct {
		Create Controller
		GetAll Controller
		Get    Controller
		Update Controller
	}

	GetReq struct {
		ID uint64
	}

	CreateReq struct {
		Name      string `json:"name"`
		Last_name string `json:"last_name"`
		Email     string `json:"email"`
	}

	UpdateReq struct {
		ID        uint64  `json:"id"`
		Name      *string `json:"name"`
		Last_name *string `json:"last_name"`
		Email     *string `json:"email"`
	}
)

func MakeEndpoint(ctx context.Context, s Service) Endpoint {
	return Endpoint{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Get:    makeGetEndpoint(s),
		Update: makeUpdateEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)
		if req.Name == "" {
			return nil, errors.New("name is required")
		}
		if req.Last_name == "" {
			return nil, errors.New("last_name is required")
		}
		if req.Email == "" {
			return nil, errors.New("email is required")
		}

		user, err := s.Create(ctx, req.Name, req.Last_name, req.Email)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		return users, nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)
		user, err := s.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		fmt.Println(req)
		return user, nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)
		if err := s.Update(ctx, req.ID, req.Name, req.Last_name, req.Email); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
