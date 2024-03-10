package response

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Response[T any] struct {
	S           bool     `json:"success"`
	M           string   `json:"message"`
	E           []string `json:"errors"`
	D           *T       `json:"data"`
	code        int
	responseFor ResponseFor
}

func New[T any](responseFor ResponseFor) *Response[T] {
	return &Response[T]{responseFor: responseFor}
}

type Sender interface {
	Send(http.ResponseWriter)
}

func (r *Response[T]) Data(data T) *Response[T] {
	r.D = &data
	return r
}

func (r *Response[T]) Success(ok bool) *Response[T] {
	r.S = ok
	return r
}

func (r *Response[T]) Error(err error) *Response[T] {
	r.E = strings.Split(err.Error(), "\n")
	return r
}

func (r *Response[T]) Code(code int) Sender {
	r.code = code
	return r
}

func (r *Response[T]) Send(w http.ResponseWriter) {
	r.M = messages[r.responseFor](len(r.E))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.code)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
