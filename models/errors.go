package models

import "errors"

// These error variables are exported to other packages as these start with capital
var (
	ErrEmailTaken = errors.New("models: email address is already in use")
	ErrNotFound   = errors.New("models: resource could not be found")
)
