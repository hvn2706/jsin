package error

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func ErrNotFound(message string) error {
	return toError(http.StatusNotFound, codes.NotFound, message)
}

func ErrNotFoundWithError(err error) error {
	return toError(http.StatusNotFound, codes.NotFound, err.Error())
}

func ErrLocked(message string) error {
	return toError(http.StatusLocked, codes.NotFound, message)
}

func ErrLockedWithError(err error) error {
	return toError(http.StatusLocked, codes.NotFound, err.Error())
}

func ErrInternalWithMsg(message string) error {
	return toError(http.StatusInternalServerError, codes.Internal, message)
}
func ErrInternalWithError(err error) error {
	return domainError{
		err:            err,
		httpStatusCode: http.StatusInternalServerError,
		code:           codes.Internal,
		message:        http.StatusText(http.StatusInternalServerError),
	}
}
