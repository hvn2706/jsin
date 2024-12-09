package error

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

// ---------------------------------------------------------------------------------------------------------------------

type domainError struct {
	err            error
	httpStatusCode int
	code           codes.Code
	message        string
}

func (e domainError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.code, e.message)
}

func toError(httpStatusCode int, code codes.Code, message string) domainError {
	return domainError{
		err:            nil,
		httpStatusCode: httpStatusCode,
		code:           code,
		message:        message,
	}
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Error   codes.Code `json:"code"`
	TraceID string     `json:"traceId"`
	Message string     `json:"message"`

	AppCode   int64  `json:"app_code,omitempty"`   // application-specific error code
	ErrorText string `json:"error_text,omitempty"` // application-level error message, for debugging
}
