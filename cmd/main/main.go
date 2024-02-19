package main

import (
	"openapi/internal/infra/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	hello "openapi/internal/ui/hello"
	items "openapi/internal/ui/stock/items"
	locations "openapi/internal/ui/stock/locations"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = validator.NewCustomValidator()

	hello.RegisterHandlers(e)
	locations.RegisterHandlers(e)
	items.RegisterHandlers(e)

	e.Logger.Fatal(e.Start(":1323"))
}
