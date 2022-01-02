package dbbackends

import (
	"errors"
	"fmt"
)

type Error interface {
	error
	ErrorWithAdditionalInfo(string) string
}

func IsDbErr(err error) bool {
	_, ok := err.(Error)
	return ok
}

func IsAlreadyExistsErr(err error) bool {
	_, ok := err.(*ErrAlreadyExists)
	return ok
}

type baseErr struct {
	resourceName string
	errorType    string
}

func (e *baseErr) Error() string {
	return e.resourceName + " " + e.errorType
}

func (e *baseErr) ErrorWithAdditionalInfo(info string) string {
	return fmt.Sprintf("%s with %s %s", e.resourceName, info, e.errorType)
}

// func (e *baseErr) GetErrType() string {
// 	return e.errorType
// }
//
// func (e *baseErr) GetResource() string {
// 	return e.resourceName
// }

type ErrNotFound struct {
	baseErr
}

func NewErrNotFound(resourceName string) *ErrNotFound {
	return &ErrNotFound{
		baseErr: baseErr{
			resourceName: resourceName,
			errorType:    "Not found",
		},
	}
}

type ErrAlreadyExists struct {
	baseErr
}

func NewErrAlreadyExists(resourceName string) *ErrAlreadyExists {
	return &ErrAlreadyExists{
		baseErr: baseErr{
			resourceName: resourceName,
			errorType:    "Already exists",
		},
	}
}

type ErrCastObjectId struct {
	baseErr
}

func NewErrCastObjectId(resourceName string) *ErrCastObjectId {
	return &ErrCastObjectId{
		baseErr: baseErr{
			resourceName: resourceName,
			errorType:    "Invlaid Id",
		},
	}
}

var ErrUnexpected struct {
	baseErr
}

func NewErrUnexpected() *ErrCastObjectId {
	return &ErrCastObjectId{
		baseErr: baseErr{
			resourceName: "",
			errorType:    "Unexpected error",
		},
	}
}

//IsNotFound is an encapsulation for checking NotFoundError logic
//Currently using errors.As as a workaround for an issue in the errors.Is function internally
func IsNotFound(err error) bool {
	var e *ErrNotFound
	isNotFound := errors.As(err, &e)
	if !isNotFound {
		return false
	}

	return true
}
