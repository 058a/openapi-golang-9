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

func TestPutOk(t *testing.T) {
	t.Parallel()

	// Setup
	h := &locations.Handler{}

	// Given
	postReqBody := &oapicodegen.PostStockLocationJSONRequestBody{
		Name: "test",
	}
	postReq := NewRequest(http.MethodPost, "/stock/locations", postReqBody)
	err := h.PostStockLocation(postReq.context)
	if err != nil {
		t.Fatal(err)
	}
	defer postReq.recorder.Result().Body.Close()

	postResBody, err := Response[oapicodegen.Created](postReq.recorder.Result())
	if err != nil {
		t.Fatal(err)
	}

	// When
	putReqBody := &oapicodegen.PutStockLocationJSONRequestBody{
		Name: "newTest",
	}
	putReq := NewRequest(http.MethodPut, "/stock/locations", putReqBody)
	err = h.PutStockLocation(putReq.context, postResBody.Id)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	defer putReq.recorder.Result().Body.Close()

	if putReq.recorder.Code != http.StatusOK {
		t.Errorf("%T %d want %d", putReq.recorder.Code, putReq.recorder.Code, http.StatusOK)
	}
}

func TestPutNotFound(t *testing.T) {
	t.Parallel()

	// Setup
	h := &locations.Handler{}

	// When
	putReqBody := &oapicodegen.PostStockLocationJSONRequestBody{
		Name: "newTest",
	}
	putReq := NewRequest(http.MethodPut, "/stock/locations", putReqBody)

	err := h.PutStockLocation(putReq.context, uuid.New())

	// Then
	if err == nil {
		t.Fatalf("expected not nil, actual nil")
	} else if err.(*echo.HTTPError).Code != http.StatusNotFound {
		t.Errorf("%T %d want %d", err.(*echo.HTTPError).Code, err.(*echo.HTTPError).Code, http.StatusNotFound)
	}
	defer putReq.recorder.Result().Body.Close()
}

func TestPutBadRequestNameEmpty(t *testing.T) {
	t.Parallel()

	// Setup
	h := &locations.Handler{}

	// Given
	postReqBody := &oapicodegen.PostStockLocationJSONRequestBody{
		Name: "test",
	}
	postReq := NewRequest(http.MethodPost, "/stock/locations", postReqBody)

	if err := h.PostStockLocation(postReq.context); err != nil {
		t.Fatal(err)
	}
	defer postReq.recorder.Result().Body.Close()

	postResBody, err := Response[oapicodegen.Created](postReq.recorder.Result())
	if err != nil {
		t.Fatal(err)
	}

	// When
	putReqBody := &oapicodegen.PutStockLocationJSONRequestBody{
		Name: "",
	}
	req := NewRequest(http.MethodPut, "/stock/locations", putReqBody)

	err = h.PutStockLocation(req.context, postResBody.Id)

	// Then
	if err == nil {
		t.Fatalf("expected not nil, actual nil")
	} else if err.(*echo.HTTPError).Code != http.StatusBadRequest {
		t.Errorf("%T %d want %d", err.(*echo.HTTPError).Code, err.(*echo.HTTPError).Code, http.StatusBadRequest)
	}

}

func TestPutBadRequestNameMaxLengthOver(t *testing.T) {
	t.Parallel()

	// Setup
	h := &locations.Handler{}

	// Given
	postReqBody := &oapicodegen.PostStockLocationJSONRequestBody{
		Name: "test",
	}
	postReq := NewRequest(http.MethodPost, "/stock/locations", postReqBody)

	if err := h.PostStockLocation(postReq.context); err != nil {
		t.Fatal(err)
	}
	defer postReq.recorder.Result().Body.Close()

	postResBody, err := Response[oapicodegen.Created](postReq.recorder.Result())
	if err != nil {
		t.Fatal(err)
	}

	// When
	putReqBody := &oapicodegen.PutStockLocationJSONRequestBody{
		Name: strings.Repeat("a", 101),
	}
	req := NewRequest(http.MethodPut, "/stock/locations", putReqBody)

	err = h.PutStockLocation(req.context, postResBody.Id)

	// Then
	if err == nil {
		t.Fatalf("expected not nil, actual nil")
	} else if err.(*echo.HTTPError).Code != http.StatusBadRequest {
		t.Errorf("%T %d want %d", err.(*echo.HTTPError).Code, err.(*echo.HTTPError).Code, http.StatusBadRequest)
	}
}
