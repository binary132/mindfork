package message

// Error is a Message containing an error.
type Error struct {
	Err error
}

// MakeError wraps an error in an Error.
func MakeError(err error) Error {
	return Error{Err: err}
}

// Error implements error on Error.
func (e Error) Error() string {
	if err := e.Err; err != nil {
		return e.Err.Error()
	}
	return ""
}
