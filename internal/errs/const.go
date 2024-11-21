package errs

import "net/http"

// nolint: lll
var (
	ErrInternal        = NewAppError(10001, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	ErrInvalidArgument = NewAppError(10002, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	ErrNotFound        = NewAppError(10003, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	ErrUnauthorized    = NewAppError(10004, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	ErrForbidden       = NewAppError(10005, http.StatusForbidden, http.StatusText(http.StatusForbidden))
	ErrTooManyRequests = NewAppError(10006, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
)
