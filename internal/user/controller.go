package user

import (
	"Gym/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoint struct {
		Create Controller
		GetAll Controller
	}

	CreateReq struct {
		Name     string `json:"name"`
		LastName string `json:"lastName"`
		Email    string `json:"email"`
	}
)

func MakeEndpoint(ctx context.Context, s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetAllUser(ctx, s, w)
		case http.MethodPost:
			decode := json.NewDecoder(r.Body)
			var user domain.User
			if err := decode.Decode(&user); err != nil {
				MsgResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			PostUser(ctx, s, w, user)

		default:
			InvalidMethod(w)
		}
	}
}
func GetAllUser(ctx context.Context, s Service, w http.ResponseWriter) {
	users, err := s.GetAll(ctx)
	if err != nil {
		MsgResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	DataResponse(w, http.StatusOK, users)
}

func PostUser(ctx context.Context, s Service, w http.ResponseWriter, data interface{}) {
	req := data.(CreateReq)
	if req.Name == "" {
		MsgResponse(w, http.StatusBadRequest, "name is required")
		return
		if req.LastName == "" {
			MsgResponse(w, http.StatusBadRequest, "lastName is required")
			return
			if req.Email == "" {
				MsgResponse(w, http.StatusBadRequest, "email is required")
				return
			}
			user, err := s.Create(ctx, req.Name, req.LastName, req.Email)
			if err != nil {
				MsgResponse(w, http.StatusInternalServerError, err.Error())
			}
			DataResponse(w, http.StatusCreated, user)
		}
	}
}
func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status":%d, "message: "method doesn't exist'}`, status)
}

func MsgResponse(w http.ResponseWriter, request int, s string) {

}

func DataResponse(w http.ResponseWriter, ok int, users []*domain.User) {

}
