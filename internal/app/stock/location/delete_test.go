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

func TestDeleteSuccess(t *testing.T) {
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

	newId := uuid.New()
	resCreate, err := app.Create(reqCreate, repo, newId)
	if err != nil {
		t.Fatal(err)
	}

	// When
	reqDelete, err := app.NewDeleteRequest(resCreate.Id)
	if err != nil {
		t.Fatal(err)
	}

	if err := app.Delete(reqDelete, repo); err != nil {
		t.Fatal(err)
	}

	// Then
	id, err := domain.NewId(resCreate.Id)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Get(id)
	if err == nil {
		t.Errorf("expected not nil, got nil")
	}
}

func TestDeleteFailIdNil(t *testing.T) {
	t.Parallel()

	// When
	_, err := app.NewDeleteRequest(uuid.Nil)
	if err != domain.ErrInvalidId {
		t.Errorf("%T %v, want %v", err, err, domain.ErrInvalidId)
	}
}

func TestDeleteFailGetFail(t *testing.T) {
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
	reqDelete, err := app.NewDeleteRequest(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	err = app.Delete(reqDelete, repo)

	// Then
	if err == nil {
		t.Errorf("expected not nil, got nil")
	}
}

func TestDeleteFailSaveFail(t *testing.T) {
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
	reqDelete, err := app.NewDeleteRequest(id.UUID())
	if err != nil {
		t.Fatal(err)
	}

	err = app.Delete(reqDelete, repo)

	// Then
	if err == nil {
		t.Errorf("expected not nil, got nil")
	}
}
