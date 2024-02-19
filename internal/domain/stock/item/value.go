package item

import (
	"errors"
)

type (
	Name struct {
		value string
	}
)

var (
	ErrNameEmpty = errors.New("Name: cannot be empty")
)

func NewName(v string) (Name, error) {
	if v == "" {
		return Name{}, ErrNameEmpty
	}
	return Name{v}, nil
}

func (v Name) String() string {
	return v.value
}
