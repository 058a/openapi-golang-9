package locations

import (
	oapicodegen "openapi/internal/infra/oapicodegen/stock/location"

	domain "openapi/internal/domain/stock/location"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	oapicodegen.ServerInterface
	Repository domain.IRepository
}

func RegisterHandlers(e *echo.Echo) {
	oapicodegen.RegisterHandlers(e, &Handler{})
}
