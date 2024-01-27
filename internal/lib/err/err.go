package err

import "errors"

/* ---------------------------------- Auth ---------------------------------- */
// ErrNoAuthenticationFound is returned when no authentication is found in the request
var ErrNoAuthenticationFound = errors.New("no authentication found")

// ErrMalformedHeader is returned when the authentication header is malformed
var ErrMalformedHeader = errors.New("malformed auth header")
