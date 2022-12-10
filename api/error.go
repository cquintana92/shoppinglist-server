package api

import "errors"

var (
	noAuthorizationError        = errors.New("NoAuthorization")
	badAuthorizationFormatError = errors.New("BadAuthorizationFormat")
	invalidAuthorizationError   = errors.New("InvalidAuthorization")
)
