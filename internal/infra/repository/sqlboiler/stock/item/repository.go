package item

import (
	"context"
	"database/sql"
	"errors"
	"openapi/internal/infra/sqlboiler"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"openapi/internal/domain/stock/item"
)

type (
	repository struct {
		item.IRepository
		db *sql.DB
	}
)

func NewRepository(db *sql.DB) (item.IRepository, error) {
	if db == nil {
		return nil, item.ErrIRepositoryDbEmpty
	}
	return &repository{
		db: db,
	}, nil
}

func (r *repository) Save(a *item.Aggregate) error {
	data := &sqlboiler.StockItem{
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

func (r *repository) Get(id item.Id) (*item.Aggregate, error) {
	data, err := sqlboiler.FindStockItem(context.Background(), r.db, id.UUID().String())
	if err != nil {
		// 想定外のエラー
		return nil, errors.Join(item.ErrIRepositoryUnexpected, err)
	}

	if data.Deleted {
		// 対象が削除されている	
		return nil, item.ErrIRepositoryRowDeleted
	}

	name, err := item.NewName(data.Name)
	if err != nil {
		// データが不正
		return nil, errors.Join(item.ErrIRepositoryInvalidData, err)
	}

	a := item.RestoreAggregate(id, name, data.Deleted)

	return a, nil
}

func (r *repository) Find(id item.Id) (bool, error) {
	data, err := sqlboiler.FindStockItem(context.Background(), r.db, id.UUID().String())
	if err != nil && err != sql.ErrNoRows {
		// 想定外のエラー
		return false, errors.Join(item.ErrIRepositoryUnexpected, err)
	}

	if err == sql.ErrNoRows {
		// 対象が見つからない
		return false, nil
	}

	if data.Deleted {
		// 対象が削除されている
		return false, nil
	}

	// 対象が見つかった
	return true, nil
}
