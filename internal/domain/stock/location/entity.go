package location

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
	ErrInvalidId = errors.New("invalid id")
)

func NewId(v uuid.UUID) (Id, error) {
	if v == uuid.Nil {
		return Id{}, ErrInvalidId
	}
	return Id{v}, nil
}

func (v Id) UUID() uuid.UUID {
	return v.value
}

func (v Id) String() string {
	return v.value.String()
}
