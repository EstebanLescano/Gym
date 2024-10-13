package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EstebanLescano/Gym/internal/user"
	"github.com/EstebanLescano/Gym/pkg/transport"
	"github.com/EstebanLescano/go-fundamentals-response/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func NewUserHttpServer(endpoint user.Endpoint) http.Handler {
	r := gin.Default()

	r.POST("/users", transport.GinServer(
		transport.Endpoint(endpoint.Create),
		decodeCreateUser,
		encodeResponse,
		encodeError))

	r.GET("/users", transport.GinServer(
		transport.Endpoint(endpoint.GetAll),
		decodeGetAllUser,
		encodeResponse,
		encodeError))

	r.GET("/users/:id", transport.GinServer(
		transport.Endpoint(endpoint.Get),
		decodeGetUser,
		encodeResponse,
		encodeError))

	r.PATCH("/users/:id", transport.GinServer(
		transport.Endpoint(endpoint.Update),
		decodeUpdateUser,
		encodeResponse,
		encodeError))
	return r
}

func decodeGetUser(c *gin.Context) (interface{}, error) {
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}
	return user.GetReq{
		ID: id,
	}, nil
}

func decodeUpdateUser(c *gin.Context) (interface{}, error) {
	var req user.UpdateReq
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(err.Error())
	}
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}
	req.ID = id
	return req, nil
}

func decodeGetAllUser(c *gin.Context) (interface{}, error) {
	if err := tokenVerify(c.GetHeader("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	return nil, nil
}

func decodeCreateUser(c *gin.Context) (interface{}, error) {

	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	var req user.CreateReq
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
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

func encodeResponse(c *gin.Context, resp interface{}) error {
	r := resp.(response.Response)
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(r.StatusCode(), resp)
	return nil
}

func encodeError(c *gin.Context, err error) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	c.JSON(resp.StatusCode(), resp)
}
