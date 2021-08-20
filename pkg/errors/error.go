package errors

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"strings"
)

type Error struct {
	Code    string
	Message string
	Err     error
}

const (
	ErrInvalid            = "invalid"
	ErrUnauthorized       = "unauthorized"
	ErrNotFound           = "not_found"
	ErrConflict           = "conflict"
	ErrInternal           = "internal"
	ErrServiceUnavailable = "service_unavailable"
	ErrNotImplemented     = "not_implemented"
)

//CheckError is for checking psql specific errors
func CheckError(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.SQLState() {
		case pgerrcode.UniqueViolation:
			split := strings.Split(pgErr.ConstraintName, "_")
			return Errorf(ErrInvalid, "%s already exists", split[1])
		default:
			return Errorf(ErrInternal, "Internal Error")
		}
	}

	return err
}

func (e *Error) Error() string {
	return e.Message
}

func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return ErrInternal
}

func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error"
}

func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}