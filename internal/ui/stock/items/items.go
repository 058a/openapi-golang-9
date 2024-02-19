package items

import (
	oapicodegen "openapi/internal/infra/oapicodegen/stock/item"

	domain "openapi/internal/domain/stock/item"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	oapicodegen.ServerInterface
	Repository domain.IRepository
}

func RegisterHandlers(e *echo.Echo) {
	oapicodegen.RegisterHandlers(e, &Handler{})
}
