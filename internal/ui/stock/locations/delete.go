package locations

import (
	"net/http"

	"github.com/labstack/echo/v4"

	app "openapi/internal/app/stock/location"
	"openapi/internal/infra/database"
	infra "openapi/internal/infra/repository/sqlboiler/stock/location"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Delete is a function that handles the HTTP DELETE request for deleting an existing stock item.
func (h *Handler) DeleteStockLocation(ctx echo.Context, stockLocationId openapi_types.UUID) error {
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

	// Validation
	reqDto, err := app.NewDeleteRequest(stockLocationId)
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
	if err := app.Delete(reqDto, repo); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Postcondition
	return ctx.JSON(http.StatusOK, nil)
}
