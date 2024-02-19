package location_test

import (
	"openapi/internal/domain/stock/location"
	"testing"

	"github.com/google/uuid"
)

func TestNewAggregate(t *testing.T) {
	t.Parallel()

	// Given
	id, err := location.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := location.NewName("TestName")
	if err != nil {
		t.Fatal(err)
	}

	// When
	a := location.NewAggregate(id, name)

	// Then
	if a.Id != id {
		t.Errorf("%T %+v want %+v", a.Id, a.Id, id)
	}

	if a.Name != name {
		t.Errorf("%T %+v want %+v", a.Name, a.Name, name)
	}

	if a.IsDeleted() != false {
		t.Errorf("%T %+v want %+v", a.IsDeleted(), a.IsDeleted(), false)
	}
}

func TestRestoreAggregate(t *testing.T) {
	t.Parallel()

	// Given
	id, err := location.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := location.NewName("TestName")
	if err != nil {
		t.Fatal(err)
	}

	// When
	a := location.RestoreAggregate(id, name, false)

	// Then
	if a.Id != id {
		t.Errorf("%T %+v want %+v", a.Id, a.Id, id)
	}

	if a.Name != name {
		t.Errorf("%T %+v want %+v", a.Name, a.Name, name)
	}

	if a.IsDeleted() != false {
		t.Errorf("%T %+v want %+v", a.IsDeleted(), a.IsDeleted(), false)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	// Given
	id, err := location.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := location.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	a := location.NewAggregate(id, name)

	// When
	a.Delete()

	// Then
	if a.IsDeleted() != true {
		t.Errorf("%T %+v want %+v", a.IsDeleted(), a.IsDeleted(), true)
	}
}
