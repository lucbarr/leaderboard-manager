package handler

type StatusCode string

const (
	CodeUnauthorized  StatusCode = "LB-001"
	CodeInternalError StatusCode = "LB-002"
	CodeBadRequest    StatusCode = "LB-003"
)
