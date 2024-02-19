package main

import (
	"net/http"
	"openapi/internal/infra/env"
	"openapi/internal/infra/validator"

	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	hello "openapi/internal/ui/hello"
	items "openapi/internal/ui/stock/items"
	locations "openapi/internal/ui/stock/locations"

	"openapi/internal/infra/auth"
)

func login(ctx echo.Context) error {
	token, err := auth.EncodeToken(uuid.New())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ctx.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+token)
	return ctx.JSON(http.StatusOK, "")
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = validator.NewCustomValidator()

	hello.RegisterHandlers(e)
	locations.RegisterHandlers(e)
	items.RegisterHandlers(e)

	e.POST("/login", login)

	r := e.Group("/stock")
	r.Use(echojwt.JWT([]byte(env.GetJwtSecret())))

	e.Logger.Fatal(e.Start(":1323"))
}
