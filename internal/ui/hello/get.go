package hello

import (
	"net/http"
	oapicodegen "openapi/internal/infra/oapicodegen/hello"

	"github.com/labstack/echo/v4"
)

// GetHello is a function that returns a JSON response with a message "Hello, World!".
//
// It takes in a parameter of type echo.Context and returns an error.
func (Api) GetHello(ctx echo.Context) error {

	return ctx.JSON(http.StatusOK, &oapicodegen.Hello{
		Message: "Hello, World!",
	})
}
