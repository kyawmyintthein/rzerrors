package rzerrors

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/kyawmyintthein/rzerrors/proto/errorpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	_unknownError = "unknown_error"
)

type RZError struct {
	messageFormat string
	cause         error
	args          []interface{}
}

func NewRZError(messageFormat string, args ...interface{}) *RZError {
	err := &RZError{
		messageFormat: messageFormat,
		args:          args,
	}
	return err
}

func (e *RZError) GetArgs() []interface{} {
	return e.args
}

func (e *RZError) GetMessage() string {
	return e.messageFormat
}

func (e *RZError) GetFormattedMessage() string {

	return fmt.Sprintf(e.messageFormat, e.args...)
}

func (e *RZError) Wrap(err error) {
	e.cause = err
}

func (e *RZError) Error() string {
	return fmt.Sprintf(e.messageFormat, e.args...)
}

func (w *RZError) Cause() error { return w.cause }

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

func GetErrorMessages(e error) string {
	return extractFullErrorMessage(e, false)
}

func GetErrorMessagesWithStack(e error) string {
	return extractFullErrorMessage(e, true)
}

func extractFullErrorMessage(e error, includeStack bool) string {
	type causer interface {
		Cause() error
	}

	var ok bool
	var lastClErr error
	errMsg := bytes.NewBuffer(make([]byte, 0, 1024))
	razerErr := e
	for {
		_, ok := razerErr.(StackTracer)
		if ok {
			lastClErr = razerErr
		}

		errorWithFormat, ok := razerErr.(ErrorFormatter)
		if ok {
			errMsg.WriteString(errorWithFormat.GetFormattedMessage())
		}

		errorCauser, ok := razerErr.(causer)
		if ok {
			innerErr := errorCauser.Cause()
			if innerErr == nil {
				break
			}
			razerErr = innerErr
		} else {
			// We have reached the end and traveresed all inner errors.
			// Add last message and exit loop.
			errMsg.WriteString(razerErr.Error())
			break
		}
		errMsg.WriteString(", ")
	}

	stackError, ok := lastClErr.(StackTracer)
	if includeStack && ok {
		errMsg.WriteString("\nSTACK TRACE:\n")
		errMsg.WriteString(stackError.GetStack())
	}
	return errMsg.String()
}

func ConvertGRPCStatusError(err error) error {
	grpcCode := codes.Unknown
	errorID := _unknownError
	httpStatus := http.StatusInternalServerError
	errorDescription := err.Error()
	var args []interface{}

	errorWithID, ok := err.(ErrorID)
	if ok {
		errorID = errorWithID.ID()
	}

	errorWithFormatter, ok := err.(ErrorFormatter)
	if ok {
		errorDescription = errorWithFormatter.GetFormattedMessage()
		args = errorWithFormatter.GetArgs()
	}

	errorWithGrpcCode, ok := err.(GRPCError)
	if ok {
		grpcCode = errorWithGrpcCode.GRPCCode()
	}

	errorWithHTTPStatus, ok := err.(HTTPError)
	if ok {
		httpStatus = errorWithHTTPStatus.StatusCode()
	}

	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, fmt.Sprintf("%v", arg))
	}

	stat := status.New(grpcCode, errorID)
	stat, _ = stat.WithDetails(&errorpb.ErrorMessage{
		HttpStatusCode:     int32(httpStatus),
		MessageDescription: errorDescription,
		Message:            errorID,
		Args:               strArgs,
		DebugInfo:          GetErrorMessagesWithStack(err),
	})
	return stat.Err()
}
