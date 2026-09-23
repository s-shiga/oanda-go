package oanda

import (
	"errors"
	"fmt"
)

// ErrStreamEnded is returned by streaming methods when the server ends the
// stream. OANDA routinely drops idle or slow stream connections, so callers
// should treat this as a signal to reconnect rather than a fatal error.
var ErrStreamEnded = errors.New("stream ended by server")

// ErrNilRequest is returned, before anything is sent, when a method that
// needs a request is passed nil. List methods whose parameters are all
// optional accept nil and use the defaults instead.
var ErrNilRequest = errors.New("request must not be nil")

// ErrNoAccountID is returned, before anything is sent, by methods that act on
// an Account when the client was created without [WithAccountID].
var ErrNoAccountID = errors.New("no account ID configured; use WithAccountID")

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

type NotFound struct{ HTTPError }

type MethodNotAllowed struct{ HTTPError }
