package errors

import (
	"fmt"
)

type ErrNotFound struct {
	What string
}

func (e ErrNotFound) String() string {
	return fmt.Sprintf("%s not found", e.What)
}

type ErrAlreadyExists struct {
	What string
}

func (e ErrAlreadyExists) String() string {
	return fmt.Sprintf("%s already exists", e.What)
}

