package rzerrors

import (
	"github.com/twitchtv/twirp"
)

type TwirpError interface {
	Code() twirp.ErrorCode
}

type ErrorWithTwirpCode struct {
	code twirp.ErrorCode
}

func NewErrorWithTwirpCode(errCode twirp.ErrorCode) *ErrorWithTwirpCode {
	return &ErrorWithTwirpCode{
		code: errCode,
	}
}

func (err *ErrorWithTwirpCode) GRPCCode() twirp.ErrorCode {
	return err.code
}
