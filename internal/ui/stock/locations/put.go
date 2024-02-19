package locations

import (
	"net/http"

	"github.com/labstack/echo/v4"

	app "openapi/internal/app/stock/location"
	"openapi/internal/infra/database"
	oapicodegen "openapi/internal/infra/oapicodegen/stock/location"
	infra "openapi/internal/infra/repository/sqlboiler/stock/location"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Put is a function that handles the HTTP PUT request for updating an existing stock location.
func (h *Handler) PutStockLocation(ctx echo.Context, stockLocationId openapi_types.UUID) error {
	// Precondition
	db, err := database.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer db.Close()

	repo, err := infra.NewRepository(db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Binding
	req := &oapicodegen.PutStockLocationJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validation
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	reqDto, err := app.NewUpdateRequest(stockLocationId, req.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	found, err := repo.Find(reqDto.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "stock location not found")
	}

	// Main
	err = app.Update(reqDto, repo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Postcondition
	return ctx.JSON(http.StatusOK, nil)
}
