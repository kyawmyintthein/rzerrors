package rzerrors

type HTTPError interface {
	StatusCode() int
}

type ErrorWithHTTPStatus struct {
	httpStatus int
}

func NewErrorWithHTTPStatus(httpStatus int) *ErrorWithHTTPStatus {
	return &ErrorWithHTTPStatus{httpStatus: httpStatus}
}

func (err *ErrorWithHTTPStatus) StatusCode() int {
	return err.httpStatus
}
