package oanda

import "fmt"

type HTTPError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("%d %s: %v", e.StatusCode, e.Message, e.Err)
}

func (e HTTPError) Unwrap() error {
	return e.Err
}

type BadRequest struct{ HTTPError }

type Unauthorized struct{ HTTPError }

type Forbidden struct{ HTTPError }

type NotFoundError struct{ HTTPError }

type MethodNotAllowed struct{ HTTPError }
