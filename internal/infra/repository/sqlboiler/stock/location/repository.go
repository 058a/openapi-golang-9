package location

import (
	"context"
	"database/sql"
	"errors"
	"openapi/internal/infra/sqlboiler"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"openapi/internal/domain/stock/location"
)

type (
	Repository struct {
		location.IRepository
		db *sql.DB
	}
)

var (
	ErrDbEmpty    = errors.New("db is empty")
	ErrRowDeleted = errors.New("row deleted")
)

func NewRepository(db *sql.DB) (*Repository, error) {
	if db == nil {
		return nil, ErrDbEmpty
	}
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Save(a *location.Aggregate) error {
	data := &sqlboiler.StockLocation{
		ID:      a.Id.String(),
		Name:    a.Name.String(),
		Deleted: a.IsDeleted(),
	}

	err := data.Upsert(
		context.Background(),
		r.db,
		true,
		[]string{"id"},
		boil.Whitelist("name", "deleted"),
		boil.Infer(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Get(id location.Id) (*location.Aggregate, error) {
	data, err := sqlboiler.FindStockLocation(context.Background(), r.db, id.UUID().String())
	if err != nil {
		return nil, err
	}

	if data.Deleted {
		return nil, ErrRowDeleted
	}

	name, err := location.NewName(data.Name)
	if err != nil {
		return nil, err
	}

	a := location.RestoreAggregate(id, name, data.Deleted)

	return a, nil
}

func (r *Repository) Find(id location.Id) (bool, error) {
	data, err := sqlboiler.FindStockLocation(context.Background(), r.db, id.UUID().String())
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	if data.Deleted {
		return false, nil
	}

	return true, nil
}
