package utils

import "errors"

var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")
var ErrBadRequest = errors.New("bad request")
var ErrInternalServer = errors.New("internal server error")
