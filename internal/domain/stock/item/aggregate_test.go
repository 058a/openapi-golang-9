package item_test

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"openapi/internal/domain/stock/item"
)

func TestNewAggregate(t *testing.T) {
	t.Parallel()

	// Given
	id, err := item.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := item.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		id   item.Id
		name item.Name
	}

	type want struct {
		id   item.Id
		name item.Name
		deleted bool
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
				name: name,
			},
			want: want{
				id: id,
				name: name,
				deleted: false,
			},
		},
	}

	// Run
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When
			got := item.NewAggregate(tt.args.id, tt.args.name)

			// Then
			if !reflect.DeepEqual(got.Id, tt.want.id) {
				t.Errorf("NewAggregate() = %v, want %v", got.Id, tt.want.id)
			}

			if !reflect.DeepEqual(got.Name, tt.want.name) {
				t.Errorf("NewAggregate() = %v, want %v", got.Name, tt.want.name)
			}

			if !reflect.DeepEqual(got.IsDeleted(), tt.want.deleted) {
				t.Errorf("NewAggregate() = %v, want %v", got.IsDeleted(), tt.want.deleted)
			}
		})
	}
}

func TestRestoreAggregate(t *testing.T) {
	t.Parallel()

	// Given
	id, err := item.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := item.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		id      item.Id
		name    item.Name
		deleted bool
	}
	type want struct {
		Id      item.Id
		Name    item.Name
		deleted bool
	}
	tests := []struct {
		name string
		args args
		want *want
	}{
		{
			name: "active",
			args: args{
				id: id,
				name: name,
				deleted: false,
			},
			want: &want{
				Id: id,
				Name: name,
				deleted: false,
			},
		},
		{
			name: "deleted",
			args: args{
				id: id,
				name: name,
				deleted: true,
			},
			want: &want{
				Id: id,
				Name: name,
				deleted: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When
			got := item.RestoreAggregate(tt.args.id, tt.args.name, tt.args.deleted)

			// Then
			if !reflect.DeepEqual(got.Id, tt.want.Id) {
				t.Errorf("RestoreAggregate() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got.Name, tt.want.Name) {
				t.Errorf("RestoreAggregate() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got.IsDeleted(), tt.want.deleted) {
				t.Errorf("RestoreAggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregate_Delete(t *testing.T) {
	t.Parallel()

	// Setup
	id, err := item.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := item.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		Id      item.Id
		Name    item.Name
		deleted bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{
				Id: id,
				Name: name,
				deleted: false,
			},
			want: true,
		},
		{
			name: "no change",
			args: args{
				Id: id,
				Name: name,
				deleted: true,
			},
			want: true,
		},
	}

	// Run
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Given
			a := item.RestoreAggregate(tt.args.Id, tt.args.Name, tt.args.deleted)

			// When
			a.Delete()

			// Then
			if got := a.IsDeleted(); got != tt.want {
				t.Errorf("Aggregate.IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}
