package items_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"openapi/internal/infra/database"
	"openapi/internal/infra/env"
	oapicodegen "openapi/internal/infra/oapicodegen/stock/item"
	infra "openapi/internal/infra/repository/sqlboiler/stock/item"
	"openapi/internal/infra/validator"
	"openapi/internal/ui/stock/items"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RequestHelper struct {
	client *http.Client
}

type Request struct {
	context  echo.Context
	recorder *httptest.ResponseRecorder
}

func NewHandler() (*items.Handler, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	repo, err := infra.NewRepository(db)
	if err != nil {
		return nil, err
	}

	return &items.Handler{Repository: repo}, nil
}

func NewRequest[I any](method string, path string, reqBody *I) *Request {
	e := echo.New()

	e.Validator = validator.NewCustomValidator()

	reqBodyJson, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBodyJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	type jwtCustomClaims struct {
		UserId string `json:"user_id"`
		jwt.RegisteredClaims
	}

	claims := &jwtCustomClaims{
		uuid.New().String(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 60)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	signedToken, _ := token.SignedString([]byte(env.GetJwtSecret()))

	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", signedToken))

	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	return &Request{
		context:  ctx,
		recorder: rec,
	}
}

func Response[T any](res *http.Response) (*T, error) {
	resBodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resBody := new(T)
	json.Unmarshal(resBodyByte, resBody)
	return resBody, nil
}

func (h *RequestHelper) Post(reqBody *oapicodegen.PostStockItemJSONRequestBody) (*http.Response, error) {
	reqBodyJson, _ := json.Marshal(reqBody)
	req, err := http.NewRequest(
		http.MethodPost,
		env.GetServiceUrl()+"/stock/items",
		bytes.NewBuffer(reqBodyJson),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *RequestHelper) Put(stockItemsId uuid.UUID, reqBody *oapicodegen.PutStockItemJSONRequestBody) (*http.Response, error) {
	reqBodyJson, _ := json.Marshal(reqBody)
	req, err := http.NewRequest(
		http.MethodPut,
		env.GetServiceUrl()+"/stock/items/"+stockItemsId.String(),
		bytes.NewBuffer(reqBodyJson),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *RequestHelper) Delete(stockItemsId uuid.UUID) (*http.Response, error) {
	req, err := http.NewRequest(
		http.MethodDelete,
		env.GetServiceUrl()+"/stock/items/"+stockItemsId.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type ResponseConvertHelper struct{}

func (h *ResponseConvertHelper) AsCreated(res *http.Response) (*oapicodegen.Created, error) {
	resBodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resBody := &oapicodegen.Created{}
	json.Unmarshal(resBodyByte, &resBody)
	if resBody.Id == uuid.Nil {
		return nil, fmt.Errorf("expected not empty, actual empty")
	}
	return resBody, nil
}

func (h *ResponseConvertHelper) AsBadRequest(res *http.Response) (*oapicodegen.BadRequestResponse, error) {
	resBodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resBody := &oapicodegen.BadRequest{}
	json.Unmarshal(resBodyByte, &resBody)

	return resBody, nil
}
