package error

import (
	"fmt"
	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"jsin/logger"
	"net/http"
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

func (e domainError) toErrorResponse() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: e.httpStatusCode,
		Error:          e.code,
		Message:        e.message,
		ErrorText:      toErrorText(e.err),
	}
}

func toError(httpStatusCode int, code codes.Code, message string) domainError {
	return domainError{
		err:            nil,
		httpStatusCode: httpStatusCode,
		code:           code,
		message:        message,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

//--
// Error response payloads & renderers
//--

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

func (e ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func ToErrorResponse(err error) render.Renderer {
	errRes, ok := err.(domainError)
	if !ok {
		logger.Errorf("Unknown error: %v", err)
		return unknownError(err)
	}

	return errRes.toErrorResponse()
}

func unknownError(err error) *ErrResponse {
	return &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		Error:          codes.Unknown,
		Message:        http.StatusText(http.StatusInternalServerError),
		ErrorText:      toErrorText(err),
	}
}

func toErrorText(err error) string {
	var errorText string
	if err != nil {
		errorText = err.Error()
	}
	return errorText
}

// ------------------
