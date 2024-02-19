package main

import (
	"net/http"
	"openapi/internal/infra/validator"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	hello "openapi/internal/ui/hello"
	items "openapi/internal/ui/stock/items"
	locations "openapi/internal/ui/stock/locations"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func login(c echo.Context) error {
	// username := c.FormValue("username")
	// password := c.FormValue("password")

	// // Throws unauthorized error
	// if username != "jon" || password != "shhh!" {
	// 	return echo.ErrUnauthorized
	// }

	// Set custom claims
	claims := &jwtCustomClaims{
		"Jon Snow",
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+t)
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
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
	e.Use(echojwt.JWT([]byte("secret")))

	e.Logger.Fatal(e.Start(":1323"))
}
