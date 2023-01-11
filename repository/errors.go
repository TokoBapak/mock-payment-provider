package repository

import "errors"

var ErrDuplicate = errors.New("duplicate")
var ErrNotFound = errors.New("not found")
var ErrExpired = errors.New("expired")
