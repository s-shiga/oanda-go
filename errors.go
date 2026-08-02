package oanda

import (
	"errors"
	"fmt"
)

// ErrStreamEnded is returned by streaming methods when the server ends the
// stream. OANDA routinely drops idle or slow stream connections, so callers
// should treat this as a signal to reconnect rather than a fatal error.
var ErrStreamEnded = errors.New("stream ended by server")

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
