package transport

import (
	"context"
	"net/http"
	"strings"
)

type Transport interface {
	Server(
		endpoint Endpoint,
		decode func(ctx context.Context, r *http.Request) (interface{}, error),
		encode func(ctx context.Context, w http.ResponseWriter, resp interface{}) error,
		encodeError func(ctx context.Context, w http.ResponseWriter, err error) error,
	)
}

type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

type transport struct {
	w   http.ResponseWriter
	r   *http.Request
	ctx context.Context
}

func New(w http.ResponseWriter, r *http.Request, ctx context.Context) Transport {
	return &transport{
		w:   w,
		r:   r,
		ctx: ctx,
	}
}

func (t *transport) Server(
	endpoint Endpoint,
	decode func(ctx context.Context, r *http.Request) (interface{}, error),
	encode func(ctx context.Context, w http.ResponseWriter, resp interface{}) error,
	encodeError func(ctx context.Context, w http.ResponseWriter, err error) error,
) {
	data, err := decode(t.ctx, t.r)
	if err != nil {
		encodeError(t.ctx, t.w, err)
		return
	}

	res, err := endpoint(t.ctx, data)
	if err != nil {
		encodeError(t.ctx, t.w, err)
	}

	if err := encode(t.ctx, t.w, res); err != nil {
		encodeError(t.ctx, t.w, err)
		return
	}
}

// esta funcion es para que acomode el path y no lo deje con espacion o le agregue la /
func Clean(url string) ([]string, int) {
	if url[0] != '/' {
		url = "/" + url
	}
	if url[len(url)-1] != '/' {
		url += "/"
	}
	parts := strings.Split(url, "/")
	return parts, len(parts)
}
