package item

import (
	"github.com/friendsofgo/errors"
	"github.com/google/uuid"
)

type (
	Id struct {
		value uuid.UUID
	}
)

var (
	ErrIdNil = errors.New("Id: cannot be nil")
)

func NewId(v uuid.UUID) (Id, error) {
	if v == uuid.Nil {
		return Id{}, ErrIdNil
	}
	return Id{v}, nil
}

func (v Id) UUID() uuid.UUID {
	return v.value
}

func (v Id) String() string {
	return v.value.String()
}
