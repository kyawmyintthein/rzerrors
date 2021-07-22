package rzerrors

type ErrorFormatter interface {
	GetArgs() []interface{}
	GetMessage() string
	GetFormattedMessage() string
}
