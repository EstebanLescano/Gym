package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EstebanLescano/Gym/internal/user"
	"github.com/EstebanLescano/Gym/pkg/transport"
	"github.com/EstebanLescano/go-fundamentals-response/response"
	"log"
	"net/http"
	"os"
	"strconv"
)

func NewUserHttpServer(ctx context.Context, router *http.ServeMux, endpoint user.Endpoint) {
	router.HandleFunc("/users/", UserServer(ctx, endpoint))
}

func UserServer(ctx context.Context, endpoint user.Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Path
		log.Println(r.Method, ":", url)
		path, pathSize := transport.Clean(url)

		params := make(map[string]string)
		if pathSize == 4 && path[2] != "" {
			params["userID"] = path[2]
		}
		params["token"] = r.Header.Get("Authorization")

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
	params, ok := ctx.Value("params").(map[string]string)
	if !ok || params["userID"] == "" {
		return nil, errors.New("missing user ID")
	}
	id, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
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
	if err := tokenVerify(params["token"]); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	id, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}
	req.ID = id
	return req, nil
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeCreateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	params := ctx.Value("params").(map[string]string)
	if err := tokenVerify(params["token"]); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("fail to decode body '%v'", err.Error()))
	}
	return req, nil
}

func tokenVerify(token string) error {
	if os.Getenv("TOKEN") != token {
		return errors.New("invalid token")
	}
	return nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Verificar si la respuesta es nil
	if resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(w).Encode(map[string]string{
			"error": "No response from server",
		})
	}
	// Intentar convertir la respuesta a response.Response
	r, ok := resp.(response.Response)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid response type",
		})
	}
	// Si es una respuesta válida, devolverla
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(ctx context.Context, w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var resp response.Response
	if ok := errors.As(err, &resp); !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
	}
	w.WriteHeader(resp.StatusCode())
	return json.NewEncoder(w).Encode(resp)
}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status":%d, "message": "method doesn't exist"}`, status)
}
