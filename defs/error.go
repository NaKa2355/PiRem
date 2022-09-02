package defs

import "errors"

var (
	ErrInvaildURLPath = errors.New("invaild URL path")
	ErrInvaildMethod  = errors.New("invaild http method")
	ErrInvaildInput   = errors.New("invaild input data")
	ErrDevice         = errors.New("device error")
	ErrInternal       = errors.New("internal error")
	ErrRequestTimeout = errors.New("request timeout")
)
