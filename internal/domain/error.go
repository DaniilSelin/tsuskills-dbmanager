package domain

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal server error")
)

type ErrorCode string

const (
	CodeOK                 ErrorCode = ""
	CodeInternal           ErrorCode = "INTERNAL_ERROR"
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodeConflict           ErrorCode = "CONFLICT"
	CodeInvalidRequestBody ErrorCode = "INVALID_REQUEST_BODY"
	CodeInvalidLimit       ErrorCode = "INVALID_LIMIT"
)
