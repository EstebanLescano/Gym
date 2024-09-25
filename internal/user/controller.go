package user

import (
	"context"
	"errors"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoint struct {
		Create Controller
		GetAll Controller
		Get    Controller
	}
	GetReq struct {
		ID uint64
	}

	CreateReq struct {
		Name      string `json:"name"`
		Last_name string `json:"last_name"`
		Email     string `json:"email"`
	}
)

func MakeEndpoint(ctx context.Context, s Service) Endpoint {
	return Endpoint{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Get:    makeGetEndpoint(s),
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
		if req.ID == 0 {
			return nil, errors.New("id is required")
		}
		return nil, nil
	}
}
