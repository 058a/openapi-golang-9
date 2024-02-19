package items_test

import (
	"net/http"
	oapicodegen "openapi/internal/infra/oapicodegen/stock/item"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestHandler_PostStockItem(t *testing.T) {
	// Setup
	t.Parallel()

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				name: "test",
			},
			wantErr: false,
		},
		{
			name: "fail/name empty",
			args: args{
				name: "",
			},
			wantErr: true,
		},
	}
	// Run
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Given
			h, err := NewHandler()
			if err != nil {
				t.Fatal(err)
			}

			postReqBody := &oapicodegen.PostStockItemJSONRequestBody{
				Name: tt.args.name,
			}
			req := NewRequest(http.MethodPost, "/stock/locations", postReqBody)

			// When
			err = h.PostStockItem(req.context)

			// Then
			if !tt.wantErr {
				// 正常系
				if err != nil {
					t.Errorf("PostStockItem() error = %v", err)
				}
				defer req.recorder.Result().Body.Close()

				if req.recorder.Code != http.StatusCreated {
					t.Errorf("PostStockItem() code = %v want %v", req.recorder.Code, http.StatusCreated)
				}

				postResBody, err := Response[oapicodegen.Created](req.recorder.Result())
				if err != nil {
					t.Errorf("PostStockItem() error = %v", err)
				}

				if postResBody.Id == uuid.Nil {
					t.Errorf("PostStockItem() Id = %v want not empty", postResBody.Id)
				}
				return
			}

			// 異常系
			if err.(*echo.HTTPError).Code != http.StatusBadRequest {
				t.Errorf("PostStockItem() code = %v want %v", err.(*echo.HTTPError).Code, http.StatusBadRequest)
			}
		})
	}
}
