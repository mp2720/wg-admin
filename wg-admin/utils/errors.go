package utils

import (
	"fmt"
)

type ErrNotFound struct {
	What string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.What)
}

type ErrAlreadyExists struct {
	What string
}

func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("%s with given unique identifiers already exists", e.What)
}
