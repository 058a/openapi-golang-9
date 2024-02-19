package item_test

import (
	"database/sql"
	"errors"
	app "openapi/internal/app/stock/item"
	domain "openapi/internal/domain/stock/item"
	"openapi/internal/infra/database"
	infra "openapi/internal/infra/repository/sqlboiler/stock/item"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestNewCreateRequest(t *testing.T) {
	t.Parallel()

	// Setup
	validId := uuid.New()
	invalidId := uuid.Nil
	validName := "test"
	invalidName := ""

	type args struct {
		id   uuid.UUID
		name string
	}
	type want struct {
		id   uuid.UUID
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
		errType error
	}{
		{
			name: "success",
			args: args{
				id:   validId,
				name: validName,
			},
			want: want{
				id:   validId,
				name: validName,
			},
			wantErr: false,
			errType: nil,
		},
		{
			name: "fail/id invalid",
			args: args{
				id:   invalidId,
				name: validName,
			},
			want:    want{},
			wantErr: true,
			errType: app.ErrValidation,
		},
		{
			name: "fail/name invalid",
			args: args{
				id:   validId,
				name: invalidName,
			},
			want:    want{},
			wantErr: true,
			errType: app.ErrValidation,
		},
	}

	// Run
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When
			got, err := app.NewCreateRequest(tt.args.id, tt.args.name)

			// Then
			if !tt.wantErr {
				// 正常系
				if err != nil {
					t.Errorf("NewCreateRequest() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got.Id.UUID() != tt.want.id {
					t.Errorf("NewCreateRequest() = %v, want %v", got.Id.UUID(), tt.want.id)
					return
				}
				if got.Name.String() != tt.want.name {
					t.Errorf("NewCreateRequest() = %v, want %v", got.Name.String(), tt.want.name)
					return
				}
				return
			}

			// 異常系
			if err == nil {
				t.Errorf("NewCreateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !errors.Is(err, tt.errType) {
				t.Errorf("NewCreateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewCreateResponse(t *testing.T) {
	t.Parallel()

	// Setup
	id, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		id domain.Id
	}
	type want struct {
		id domain.Id
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success",
			args: args{
				id: id,
			},
			want: want{
				id: id,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When
			got := app.NewCreateResponse(tt.args.id)

			// Then
			if !reflect.DeepEqual(got.Id, tt.want.id.UUID()) {
				t.Errorf("NewCreateResponse() = %v, want %v", got.Id, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	// Setup
	validDb, err := database.Open()
	if err != nil {
		t.Fatal()
	}
	t.Cleanup(func() {
		validDb.Close()
	})

	invalidDb, err := database.Open()
	if err != nil {
		t.Fatal()
	}
	invalidDb.Close()

	validId := uuid.New()

	type args struct {
		db   *sql.DB
		id   uuid.UUID
		name string
	}		
	type want struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				db:   validDb,
				id:   validId,
				name: "test",
			},
			want: want{
				id: validId,
			},
			wantErr: false,
		},
		{
			name: "fail/db connection error",
			args: args{
				db:   invalidDb,
				id:   validId,
				name: "test",
			},
			want:    want{},
			wantErr: true,
		},
	}

	// Run
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Given
			repo, err := infra.NewRepository(tt.args.db)
			if err != nil {
				t.Fatal(err)
			}

			req, err := app.NewCreateRequest(tt.args.id, tt.args.name)
			if err != nil {
				t.Fatal(err)
			}

			// When
			got, err := app.Create(req, repo)

			// Then
			if !tt.wantErr {
				// 正常系
				if err != nil {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got.Id, tt.want.id) {
					t.Errorf("Create() = %v, want %v", got.Id, tt.want)
					return
				}
				return
			}

			// 異常系
			if err == nil {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
