package tra

import "net/http"

func NewHTTPError(statusCode int, message string) HTTPError {
	return HTTPError{
		statusCode: statusCode,
		Message:    message,
	}
}

type HTTPError struct {
	statusCode int

	Message string
}

var (
	UnauthorizedError = NewHTTPError(http.StatusUnauthorized, "401 unauthorized")
	ForbiddenError    = NewHTTPError(http.StatusForbidden, "403 forbidden")
)
