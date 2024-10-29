package errs

import (
	"fmt"
	"runtime"
)

type AppError struct {
	AppCode    int    `json:"app_code"`    // APP code
	StatusCode int    `json:"status_code"` // REST status code
	Message    string `json:"message"`     // Error message

	err        error  `json:"-"`
	stackTrace string `json:"-"`
}

func NewAppError(
	appCode int,
	statusCode int,
	message string,
) AppError {
	return AppError{
		AppCode:    appCode,
		StatusCode: statusCode,
		Message:    message,
		stackTrace: stackTrace(),
	}
}

func (e AppError) New() AppError {
	e.stackTrace = stackTrace()
	return e
}

func (e AppError) Wrap(err error) AppError {
	e.err = err
	return e
}

func (e AppError) Rewrap(err error) AppError {
	e.err = err
	e.stackTrace = stackTrace()
	return e
}

func (e AppError) Msg(message string) AppError {
	e.Message = message
	return e
}

func (e AppError) Msgf(message string, args ...any) AppError {
	e.Message = fmt.Sprintf(message, args...)
	return e
}

func (e AppError) ReMsg(message string) AppError {
	e.Message = message
	e.stackTrace = stackTrace()
	return e
}

func (e AppError) ReMsgf(message string, args ...any) AppError {
	e.Message = fmt.Sprintf(message, args...)
	e.stackTrace = stackTrace()
	return e
}

func (e AppError) Is(in error) bool {
	ee, ok := in.(AppError)
	return ok && ee.AppCode == e.AppCode
}

func (e AppError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

func (e AppError) GetStackTrace() string {
	return e.stackTrace
}

func (e AppError) GetErr() error {
	return e.err
}

func stackTrace() string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", file, line)
}
