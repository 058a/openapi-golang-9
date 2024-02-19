package locations_test

import (
	oapicodegen "openapi/internal/infra/oapicodegen/stock/location"
	"openapi/internal/ui/stock/locations"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/google/uuid"

	"net/http"
)

func TestPostCreated(t *testing.T) {
	t.Parallel()

	// Setup
	h := &locations.Handler{}

	// When
	postReqBody := &oapicodegen.PostStockLocationJSONRequestBody{
		Name: "test",
	}
	req := NewRequest(http.MethodPost, "/stock/locations", postReqBody)
	err := h.PostStockLocation(req.context)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	defer req.recorder.Result().Body.Close()

	if req.recorder.Code != http.StatusCreated {
		t.Errorf("%T %d want %d", err, req.recorder.Code, http.StatusCreated)
	}

	postResBody, err := Response[oapicodegen.Created](req.recorder.Result())
	if err != nil {
		t.Fatal(err)
	}

	if postResBody.Id == uuid.Nil {
		t.Errorf("expected not empty, actual empty")
	}
}

func TestPostBadRequestNameEmpty(t *testing.T) {
	t.Parallel()

	// Setup
	h, err := NewHandler()
	if err != nil {
		t.Fatal(err)
	}

	// When
	postReqBody := &oapicodegen.PostStockLocationJSONRequestBody{
		Name: "",
	}
	req := NewRequest(http.MethodPost, "/stock/locations", postReqBody)
	err = h.PostStockLocation(req.context)

	// Then
	if err == nil {
		t.Fatalf("expected not nil, actual nil")
	}

	if err.(*echo.HTTPError).Code != http.StatusBadRequest {
		t.Errorf("%T %d want %d", err.(*echo.HTTPError).Code, err.(*echo.HTTPError).Code, http.StatusBadRequest)
	}
}

func TestPostBadRequestNameMaxLengthOver(t *testing.T) {
	t.Parallel()

	// Setup
	h, err := NewHandler()
	if err != nil {
		t.Fatal(err)
	}

	// When
	reqBody := &oapicodegen.PostStockLocationJSONRequestBody{
		Name: strings.Repeat("a", 101),
	}
	req := NewRequest(http.MethodPost, "/stock/locations", reqBody)

	// Then
	err = h.PostStockLocation(req.context)
	if err == nil {
		t.Fatalf("expected not nil, actual nil")
	}

	// Then
	// if err == nil {
	// 	t.Fatalf("expected not nil, actual nil")
	// }

	if err.(*echo.HTTPError).Code != http.StatusBadRequest {
		t.Errorf("%T %d want %d", err.(*echo.HTTPError).Code, err.(*echo.HTTPError).Code, http.StatusBadRequest)
	}
}
