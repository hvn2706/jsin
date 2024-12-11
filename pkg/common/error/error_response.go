package error

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

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
