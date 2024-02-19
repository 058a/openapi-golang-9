package location_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	domain "openapi/internal/domain/stock/location"
	"openapi/internal/infra/database"
	infra "openapi/internal/infra/repository/sqlboiler/stock/location"
	"openapi/internal/infra/sqlboiler"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestNewRepository(t *testing.T) {
	t.Parallel()

	// Setup
	db, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// When
	repo, err := infra.NewRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	if repo == nil {
		t.Errorf("expected not nil, got nil")
	}
}

func TestNewRepositoryFail(t *testing.T) {
	t.Parallel()

	// When
	repo, err := infra.NewRepository(nil)

	// Then
	if err != infra.ErrDbEmpty {
		t.Errorf("%T %+v want %+v", err, err, infra.ErrDbEmpty)
	}

	if repo != nil {
		t.Errorf("expected nil, got not nil")
	}
}

func TestSave(t *testing.T) {
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
	id, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := domain.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	a := domain.NewAggregate(id, name)

	// When
	beforeCreatedAt := time.Now().UTC()
	if err = repo.Save(a); err != nil {
		t.Fatal(err)
	}

	beforeData, err := sqlboiler.FindStockLocation(context.Background(), db, a.Id.String())
	if err != nil {
		t.Fatal(err)
	}

	beforeModel, err := repo.Get(a.Id)
	if err != nil {
		t.Fatalf("expected error but returned nil, %+v", err)
	}

	if err = repo.Save(a); err != nil {
		t.Fatal(err)
	}

	afterData, err := sqlboiler.FindStockLocation(context.Background(), db, a.Id.String())
	if err != nil {
		t.Fatal(err)
	}

	afterModel, err := repo.Get(a.Id)
	if err != nil {
		t.Fatalf("expected error but returned nil, %+v", err)
	}

	// Then
	if !reflect.DeepEqual(beforeModel, afterModel) {
		t.Errorf("%T %+v want %+v", afterModel, afterModel, beforeModel)
	}

	if beforeData.ID != id.String() {
		t.Errorf("%T %+v want %+v", beforeData.ID, beforeData.ID, id)
	}

	if beforeData.Name != name.String() {
		t.Errorf("%T %+v want %+v", beforeData.Name, beforeData.Name, name)
	}

	if beforeData.Deleted != false {
		t.Errorf("%T %+v want %+v", beforeData.Deleted, beforeData.Deleted, false)
	}

	beforeCreatedAtNC := beforeCreatedAt.UnixMilli()
	createdAtAtNC := beforeData.CreatedAt.UnixMilli()
	if !(createdAtAtNC >= beforeCreatedAtNC) {
		t.Errorf("%T %d want greater than %d", createdAtAtNC, createdAtAtNC, beforeCreatedAtNC)
	}

	if beforeData.UpdatedAt.Equal(beforeData.CreatedAt) != true {
		t.Errorf("%T %+v want %+v", beforeData.UpdatedAt, beforeData.UpdatedAt, beforeData.CreatedAt)
	}

	if afterData.ID != beforeData.ID {
		t.Errorf("%T %+v want %+v", afterData.ID, afterData.ID, beforeData.ID)
	}

	if afterData.Name != beforeData.Name {
		t.Errorf("%T %+v want %+v", afterData.Name, afterData.Name, beforeData.ID)
	}

	if afterData.Deleted != beforeData.Deleted {
		t.Errorf("%T %+v want %+v", afterData.Deleted, afterData.Deleted, beforeData.Deleted)
	}

	if afterData.CreatedAt != beforeData.CreatedAt {
		t.Errorf("%T %+v want %+v", afterData.CreatedAt, afterData.CreatedAt, beforeData.CreatedAt)
	}

	if afterData.UpdatedAt.Before(afterData.CreatedAt) == true {
		t.Errorf("%T %+v want greater than %+v ", afterData.UpdatedAt, afterData.UpdatedAt, afterData.CreatedAt)
	}
}

func TestGet(t *testing.T) {
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
	id, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := domain.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	a := domain.NewAggregate(id, name)

	// When
	notFound, err := repo.Get(a.Id)
	if err == nil {
		t.Fatal("expected not nil, got nil")
	}

	if err = repo.Save(a); err != nil {
		t.Fatal(err)
	}

	found, err := repo.Get(a.Id)
	if err != nil {
		t.Fatal(err)
	}

	a.Delete()

	if err = repo.Save(a); err != nil {
		t.Fatal(err)
	}

	deleted, err := repo.Get(a.Id)
	if err == nil {
		t.Fatal("expected not nil, got nil")
	}

	// Then
	if notFound != nil {
		t.Errorf("%T %+v want %+v", notFound, notFound, nil)
	}

	if found == nil {
		t.Errorf("expected not nil, got nil")
	}

	if deleted != nil {
		t.Errorf("%T %+v want %+v", deleted, deleted, nil)
	}
}

func TestFind(t *testing.T) {
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
	id, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := domain.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	a := domain.NewAggregate(id, name)

	// When
	notFound, err := repo.Find(a.Id)
	if err != nil {
		t.Fatal(err)
	}

	if err = repo.Save(a); err != nil {
		t.Fatal(err)
	}

	found, err := repo.Find(a.Id)
	if err != nil {
		t.Fatal(err)
	}

	a.Delete()

	if err = repo.Save(a); err != nil {
		t.Fatal(err)
	}

	notFoundDueToDeleted, err := repo.Find(a.Id)
	if err != nil {
		t.Fatal(err)
	}

	// Then
	if notFound != false {
		t.Errorf("%T %+v want %+v", notFound, notFound, false)
	}

	if found != true {
		t.Errorf("%T %+v want %+v", found, found, false)
	}

	if notFoundDueToDeleted != false {
		t.Errorf("%T %+v want %+v", notFoundDueToDeleted, notFoundDueToDeleted, false)
	}
}

func TestFailDbClose(t *testing.T) {
	t.Parallel()

	// Setup
	db, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	db.Close()

	repo, err := infra.NewRepository(db)
	if err != nil {
		t.Fatal(err)
	}

	// Given
	id, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	name, err := domain.NewName("test")
	if err != nil {
		t.Fatal(err)
	}

	a := domain.NewAggregate(id, name)

	// Then
	_, err = repo.Get(a.Id)
	if err == nil {
		t.Errorf("expected error but returned nil")
	}

	_, err = repo.Find(a.Id)
	if err == nil {
		t.Errorf("expected error but returned nil")
	}

	err = repo.Save(a)
	if err == nil {
		t.Errorf("expected error but returned nil")
	}
}

func TestGetFailDataInvalid(t *testing.T) {
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
	id, err := domain.NewId(uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	data := &sqlboiler.StockLocation{
		ID:      id.String(),
		Name:    "",
		Deleted: false,
	}

	err = data.Upsert(
		context.Background(),
		db,
		true,
		[]string{"id"},
		boil.Whitelist("name", "deleted"),
		boil.Infer(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// When
	_, err = repo.Get(id)

	// Then
	if err == nil {
		t.Errorf("expected not nil, got nil")
	}
}
