package strategy

import "errors"

var (
	ErrUnknownStrategy = errors.New("unknown strategy kind")
	ErrNotImplemented  = errors.New("not implemented")
	ErrInvalidArg      = errors.New("invalid strategy argument")
)
