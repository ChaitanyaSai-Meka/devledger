package repository

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrGroupAlreadyExists = errors.New("group already exists")
	ErrUserAlreadyInGroup = errors.New("user already in group")
)
