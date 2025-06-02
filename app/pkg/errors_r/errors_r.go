package errors_r

import "errors"

var (
	ErrWrongNumberOfArguments = errors.New("wrong number of arguments")
	ErrInvalidInteger         = errors.New("invalid integer")
	ErrSyntaxError            = errors.New("syntax error")
	ErrInvalidMessage         = errors.New("invalid message format")
	ErrSlaveClosedConn        = errors.New("slave closed conn")

	ErrKeyNotFound      = errors.New("key not found")
	ErrKeyNotFoundInMap = errors.New("key not found in map")
	ErrKeyNotFoundInRDB = errors.New("key not found in rdb")
	ErrInvalidRequest   = errors.New("invalid request")
)
