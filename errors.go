package oanda

import "fmt"

type BadRequest struct {
	Code    int
	Message string
	Err     error
}

func (e BadRequest) Error() string {
	return fmt.Sprintf("400 bad request: %s", e.Err.Error())
}

func (e BadRequest) Unwrap() error {
	return e.Err
}

type Unauthorized struct {
	Code    int
	Message string
	Err     error
}

func (e Unauthorized) Error() string {
	return fmt.Sprintf("401 unauthorized: %s", e.Err.Error())
}

func (e Unauthorized) Unwrap() error {
	return e.Err
}

type Forbidden struct {
	Code    int
	Message string
	Err     error
}

func (e Forbidden) Error() string {
	return fmt.Sprintf("403 forbidden: %s", e.Err.Error())
}

func (e Forbidden) Unwrap() error {
	return e.Err
}

type NotFoundError struct {
	Code    int
	Message string
	Err     error
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("404 not found: %s", e.Message)
}

func (e NotFoundError) Unwrap() error {
	return e.Err
}

type MethodNotAllowed struct {
	Code    int
	Message string
	Err     error
}

func (e MethodNotAllowed) Error() string {
	return fmt.Sprintf("405 method not allowed: %s", e.Message)
}

func (e MethodNotAllowed) Unwrap() error {
	return e.Err
}
