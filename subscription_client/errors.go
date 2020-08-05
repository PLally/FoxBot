package subscription_client

import "errors"

var (
	ErrAlreadyExists = errors.New("The resource already exists")
	ErrNoPermissions = errors.New("No permissions to access that endpoint")
)

type SubError struct {
	s   string // a readable message to show to the user
	err error
}

func (e SubError) Error() string {
	return e.s
}

func (e SubError) Unwrap() error {
	return e.err
}
