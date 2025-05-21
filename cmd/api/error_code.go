package main

type ErrorCode string

const (
	ErrValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrNotFound         ErrorCode = "NOT_FOUND"
	ErrConflict         ErrorCode = "CONFLICT"
	ErrInternal         ErrorCode = "INTERNAL"
)
