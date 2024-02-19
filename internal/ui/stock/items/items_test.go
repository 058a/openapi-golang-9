package items_test

import (
	"testing"

	"github.com/labstack/echo/v4"

	"openapi/internal/ui/stock/items"
)

func TestRegisterHandlers(t *testing.T) {
	t.Parallel()

	type args struct {
		e *echo.Echo
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				e: echo.New(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When
			items.RegisterHandlers(tt.args.e)

			// Then
			tt.args.e.Start(":1323")
			tt.args.e.Close()
		})
	}
}
