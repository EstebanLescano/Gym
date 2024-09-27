package handler

import (
	"Gym/internal/user"
	"Gym/pkg/transport"
	"context"
	"encoding/json"
	. "fmt"
	"log"
	"net/http"
	"strconv"
)

func NewUserHttpServer(ctx context.Context, router *http.ServeMux, endpoint user.Endpoint) {
	router.HandleFunc("/users/", UserServer(ctx, endpoint))
}

func UserServer(ctx context.Context, endpoint user.Endpoint) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Path
		log.Println(r.Method, ":", url)
		path, pathSize := transport.Clean(url)

		params := make(map[string]string)
		if pathSize == 4 && path[2] != "" {
			params["userID"] = path[2]
		}
		tran := transport.New(w, r, context.WithValue(ctx, "params", params))

		var end user.Controller
		var deco func(ctx context.Context, r *http.Request) (interface{}, error)

		switch r.Method {
		case http.MethodGet:
			switch pathSize {
			case 3:
				end = endpoint.GetAll
				deco = decodeGetAllUser
			case 4:
				end = endpoint.Get
				deco = decodeGetUser
			}

		case http.MethodPost:
			switch pathSize {
			case 3:
				end = endpoint.Create
				deco = decodeCreateUser
			}
		case http.MethodPatch:
			switch pathSize {
			case 4:
				end = endpoint.Update
				deco = decoUpdateUser
			}
		}
		if end != nil && deco != nil {
			tran.Server(
				transport.Endpoint(end),
				deco,
				encodeResponse,
				encodeError)
		} else {
			InvalidMethod(w)
		}
	}
}

func decodeGetUser(ctx context.Context, r *http.Request) (interface{}, error) {
	params := ctx.Value("params").(map[string]string)
	id, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		return nil, err
	}
	return user.GetReq{
		ID: id,
	}, nil
}

func decoUpdateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	var req user.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	params := ctx.Value("params").(map[string]string)
	id, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		return nil, err
	}
	req.ID = id
	return req, nil
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, Errorf("Error getUser")
}

func decodeCreateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, Errorf("fail to decode body '%v'", err.Error())
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	status := http.StatusOK
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func encodeError(ctx context.Context, w http.ResponseWriter, err error) error {
	status := http.StatusInternalServerError
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	errorResponse := map[string]interface{}{
		"status":  status,
		"message": err.Error(),
	}
	data, jsonErr := json.Marshal(errorResponse)
	if jsonErr != nil {
		http.Error(w, `{"status":500, "message":"Error serializing error response"}`, http.StatusInternalServerError)
		return jsonErr
	}

	_, writeErr := w.Write(data)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.WriteHeader(status)
	Fprintf(w, `{"status":%d, "message: "method doesn't exist'}`, status)
}
