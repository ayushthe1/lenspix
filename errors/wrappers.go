package errors

import "errors"

// These variables are used to give us access to existing
// functions in the std lib errors package. We can also
// wrap them in custom functionality as needed if we want,
// or mock them during testing.
var (
	As = errors.As
	Is = errors.Is
)
