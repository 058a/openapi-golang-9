package location

import (
	"errors"
)

type (
	Name struct {
		string
	}
)

var (
	ErrInvalidName = errors.New("invalid name")
)

func NewName(v string) (Name, error) {
	if v == "" {
		return Name{}, ErrInvalidName
	}
	return Name{v}, nil
}

func (v Name) String() string {
	return v.string
}
