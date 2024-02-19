package item

import "github.com/friendsofgo/errors"

type IRepository interface {
	Save(a *Aggregate) error
	Get(id Id) (*Aggregate, error)
	Find(id Id) (bool, error)
}

var (
	ErrIRepositoryDbEmpty     = errors.New("IRepository: db is empty")
	ErrIRepositoryRowDeleted  = errors.New("IRepository: row deleted")
	ErrIRepositoryInvalidData = errors.New("IRepository: invalid data")
	ErrIRepositoryUnexpected  = errors.New("IRepository: unexpected error")
)
