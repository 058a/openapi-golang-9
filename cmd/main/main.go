package main

import (
	"net/http"
	"openapi/internal/infra/env"
	"openapi/internal/infra/validator"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	hello "openapi/internal/ui/hello"
	items "openapi/internal/ui/stock/items"
	locations "openapi/internal/ui/stock/locations"
)

func login(c echo.Context) error {
	type jwtCustomClaims struct {
		UserId string `json:"user_id"`
		jwt.RegisteredClaims
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		uuid.New().String(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 60)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	signedToken, err := token.SignedString([]byte(env.GetJwtSecret()))
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+signedToken)
	return c.JSON(http.StatusOK, "")
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
