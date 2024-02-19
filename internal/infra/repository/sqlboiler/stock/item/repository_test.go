package item_test

import (
	"context"
	"database/sql"
	"errors"
	domain "openapi/internal/domain/stock/item"
	"openapi/internal/infra/database"
	infra "openapi/internal/infra/repository/sqlboiler/stock/item"
	"openapi/internal/infra/sqlboiler"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestNewRepository(t *testing.T) {
	t.Parallel()

	//Setup
	db, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	type repository struct {
		domain.IRepository
		db *sql.DB
	}
	repository2 := &repository{
		db: db,
	}

	//Given
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    domain.IRepository
		wantErr bool
		errType error
	}{
		{
			name: "success",
			args: args{
				db: db,
			},
			want:    repository2,
			wantErr: false,
			errType: nil,
		},
		{
			name: "fail",
			args: args{
				db: nil,
			},
			want:    nil,
			wantErr: true,
			errType: domain.ErrIRepositoryDbEmpty,
		},
	}

	//When & Then
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := infra.NewRepository(tt.args.db)
			if !tt.wantErr {
				if got == nil {
					t.Errorf("NewRepository() = %v, want not nil", got)
					return
				}
				return
			}
			if !errors.Is(err, tt.errType) {
				t.Errorf("NewRepository() error = %v, wantErr %v", err, tt.errType)
				return
			}
		})
	}
}

func TestRepository_Save(t *testing.T) {
	t.Parallel()

	//Setup
	db, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	closeDb, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	closeDb.Close()

	id, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := domain.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	a := domain.NewAggregate(id, name)

	// Given
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		errKeyword string
	}{
		{
			name: "success",
			args: args{
				db: db,
			},
			wantErr:    false,
			errKeyword: "",
		},
		{
			name: "fail",
			args: args{
				db: closeDb,
			},
			wantErr:    true,
			errKeyword: "sql: database is closed",
		},
	}

	// When & Then
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := infra.NewRepository(tt.args.db)
			if err != nil {
				t.Fatal(err)
			}

			if err := repo.Save(a); err == nil && tt.wantErr && !strings.Contains(err.Error(), tt.errKeyword) {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repository_Get(t *testing.T) {
	t.Parallel()

	//Setup
	//正常なDBインスタンス
	db, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	//正常なリポジトリ
	repo, err := infra.NewRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	//接続を切ったDBインスタンス
	closedDb, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	closedDb.Close()

	//異常なリポジトリ
	errRepo, err := infra.NewRepository(closedDb)
	if err != nil {
		t.Fatal(err)
	}

	//Id
	id, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	//SaveしないId
	noSaveId, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	//Name
	name, err := domain.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	//Saveする集約
	a := domain.NewAggregate(id, name)
	if err := repo.Save(a); err != nil {
		t.Fatal(err)
	}

	//削除した状態でSaveする集約
	deleteId, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	b := domain.NewAggregate(deleteId, name)
	b.Delete()
	if err := repo.Save(b); err != nil {
		t.Fatal(err)
	}

	//不正なデータのId
	invalidId, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	data := &sqlboiler.StockItem{
		ID:      invalidId.String(),
		Name:    "",
		Deleted: false,
	}
	if err := data.Upsert(
		context.Background(),
		db,
		true,
		[]string{"id"},
		boil.Whitelist("name", "deleted"),
		boil.Infer(),
	); err != nil {
		t.Fatal(err)
	}

	//Given
	type args struct {
		repo domain.IRepository
		id   domain.Id
	}
	tests := []struct {
		name       string
		args       args
		want       *domain.Aggregate
		wantErr    bool
		errKeyword string
	}{
		//成功する場合
		{
			name: "success",
			args: args{
				repo: repo,
				id:   a.Id,
			},
			want:       a,
			wantErr:    false,
			errKeyword: "",
		},
		//DB未接続で失敗する場合
		{
			name: "fail/db not connected",
			args: args{
				repo: errRepo,
				id:   a.Id,
			},
			want:       nil,
			wantErr:    true,
			errKeyword: "sql: database is closed",
		},
		// 見つからない場合
		{
			name: "fail/not found",
			args: args{
				repo: repo,
				id:   noSaveId,
			},
			want:       nil,
			wantErr:    true,
			errKeyword: "sql: no rows in result set",
		},
		//削除されていた場合
		{
			name: "fail/deleted",
			args: args{
				repo: repo,
				id:   deleteId,
			},
			want:       b,
			wantErr:    true,
			errKeyword: domain.ErrIRepositoryRowDeleted.Error(),
		},
		//不正な状態で保存されていた場合
		{
			name: "fail/invalid data",
			args: args{
				repo: repo,
				id:   invalidId,
			},
			want:       nil,
			wantErr:    true,
			errKeyword: domain.ErrIRepositoryInvalidData.Error(),
		},
	}

	//When & Then
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.args.repo.Get(tt.args.id)
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("repository.Get() = %v, want %v", got, tt.want)
				}
				return
			}

			if err == nil {
				t.Errorf("repository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(err.Error(), tt.errKeyword) {
				t.Errorf("repository.Get() error = %v, wantErr %v", err, tt.errKeyword)
				return
			}
		})
	}
}

func Test_repository_Find(t *testing.T) {
	t.Parallel()

	//Setup
	//正常なDBインスタンス
	db, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	//正常なリポジトリ
	repo, err := infra.NewRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	//接続を切ったDBインスタンス
	closedDb, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	closedDb.Close()

	//異常なリポジトリ
	errRepo, err := infra.NewRepository(closedDb)
	if err != nil {
		t.Fatal(err)
	}

	//Id
	id, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	//SaveしないId
	noSaveId, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	//Name
	name, err := domain.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	//Saveする集約
	a := domain.NewAggregate(id, name)
	if err := repo.Save(a); err != nil {
		t.Fatal(err)
	}

	//削除した状態でSaveする集約
	deleteId, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	b := domain.NewAggregate(deleteId, name)
	b.Delete()
	if err := repo.Save(b); err != nil {
		t.Fatal(err)
	}

	//Given
	type args struct {
		repo domain.IRepository
		id   domain.Id
	}
	tests := []struct {
		name       string
		args       args
		want       bool
		wantErr    bool
		errKeyword string
	}{
		//成功する場合
		{
			name: "success",
			args: args{
				repo: repo,
				id:   a.Id,
			},
			want:       true,
			wantErr:    false,
			errKeyword: "",
		},
		//DB未接続で失敗する場合
		{
			name: "fail/db not connected",
			args: args{
				repo: errRepo,
				id:   a.Id,
			},
			want:       false,
			wantErr:    true,
			errKeyword: domain.ErrIRepositoryUnexpected.Error(),
		},
		// 見つからない場合
		{
			name: "fail/not found",
			args: args{
				repo: repo,
				id:   noSaveId,
			},
			want:       false,
			wantErr:    false,
			errKeyword: "",
		},
		//削除されていた場合
		{
			name: "fail/deleted",
			args: args{
				repo: repo,
				id:   deleteId,
			},
			want:       false,
			wantErr:    false,
			errKeyword: "",
		},
	}

	//When & Then
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.args.repo.Find(tt.args.id)
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("repository.Find() = %v, want %v", got, tt.want)
				}
				return
			}

			if err == nil {
				t.Errorf("repository.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(err.Error(), tt.errKeyword) {
				t.Errorf("repository.Find() error = %v, wantErr %v", err, tt.errKeyword)
				return
			}
		})
	}
}
