package rzerrors

import (
	"google.golang.org/grpc/codes"
)

type GRPCError interface {
	GRPCCode() codes.Code
}

type ErrorWithGRPCCode struct {
	code codes.Code
}

func NewErrorWithgRPCCode(grpcCode codes.Code) *ErrorWithGRPCCode {
	return &ErrorWithGRPCCode{
		code: grpcCode,
	}
}

func (err *ErrorWithGRPCCode) GRPCCode() codes.Code {
	return err.code
}
