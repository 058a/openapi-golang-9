package location_test

import (
	"errors"
	app "openapi/internal/app/stock/location"
	mock "openapi/internal/app/stock/location/internal"
	domain "openapi/internal/domain/stock/location"
	"openapi/internal/infra/database"
	infra "openapi/internal/infra/repository/sqlboiler/stock/location"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestUpdateSuccess(t *testing.T) {
	t.Parallel()

	// Setup
	db, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := infra.NewRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	// Given
	name := "TestName"

	reqCreate, err := app.NewCreateRequest(name)
	if err != nil {
		t.Fatal(err)
	}

	resCreate, err := app.Create(reqCreate, repo, uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	// When
	newName := "NewTestName"
	reqUpdate, err := app.NewUpdateRequest(resCreate.Id, newName)
	if err != nil {
		t.Fatal(err)
	}

	err = app.Update(reqUpdate, repo)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	a, err := repo.Get(reqUpdate.Id)
	if err != nil {
		t.Fatal(err)
	}

	if a.Name.String() != newName {
		t.Errorf("%T %v, want %v", a.Name.String(), a.Name.String(), newName)
	}
}

func TestUpdateFailNameInvalid(t *testing.T) {
	t.Parallel()

	// When
	name := ""
	_, err := app.NewUpdateRequest(uuid.New(), name)
	if err != domain.ErrInvalidName {
		t.Errorf("%T %v, want %v", err, err, domain.ErrInvalidName)
	}
}

func TestUpdateFailIdNil(t *testing.T) {
	t.Parallel()

	// When
	name := "test"
	_, err := app.NewUpdateRequest(uuid.Nil, name)
	if err != domain.ErrInvalidId {
		t.Errorf("%T %v, want %v", err, err, domain.ErrInvalidId)
	}
}

func TestUpdateFailGetFail(t *testing.T) {
	t.Parallel()

	// Setup
	db, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo, err := infra.NewRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	// When
	newName := "test"
	reqUpdate, err := app.NewUpdateRequest(uuid.New(), newName)
	if err != nil {
		t.Fatal(err)
	}

	err = app.Update(reqUpdate, repo)

	// Then
	if err == nil {
		t.Errorf("expected not nil, got nil")
	}
}

func TestUpdateFailSaveFail(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := domain.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	a := domain.NewAggregate(id, name)

	repo := mock.NewMockIRepository(ctrl)
	repo.EXPECT().Get(gomock.Any()).Return(a, nil)
	repo.EXPECT().Save(gomock.Any()).Return(errors.New("test error"))

	// When
	reqUpdate, err := app.NewUpdateRequest(id.UUID(), name.String())
	if err != nil {
		t.Fatal(err)
	}

	err = app.Update(reqUpdate, repo)

	// Then
	if err == nil {
		t.Errorf("expected not nil, got nil")
	}
}
