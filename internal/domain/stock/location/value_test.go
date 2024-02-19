package location_test

import (
	"testing"

	"openapi/internal/domain/stock/location"
)

func TestNewName(t *testing.T) {
	t.Parallel()

	// When
	value := "test"
	name, err := location.NewName(value)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	if name.String() != value {
		t.Errorf("%T %+v want %+v", name, name, value)
	}
}

func TestNewNameFail(t *testing.T) {
	t.Parallel()

	// When
	value := ""
	name, err := location.NewName(value)
	if err == nil {
		t.Fatal("expected error but returned nil")
	}

	// Then
	if name.String() != value {
		t.Errorf("%T %+v want %+v", name, name, value)
	}
}
