package items

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	app "openapi/internal/app/stock/item"
	"openapi/internal/infra/database"
	oapicodegen "openapi/internal/infra/oapicodegen/stock/item"
	infra "openapi/internal/infra/repository/sqlboiler/stock/item"
)

// PostStockItem is a function that handles the HTTP POST request for creating a new stock item.
func (h *Handler) PostStockItem(ctx echo.Context) error {

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
	req := &oapicodegen.PostStockItemJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validation
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	reqDto, err := app.NewCreateRequest(uuid.New(), req.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Main Process
	resDto, err := app.Create(reqDto, repo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	res := &oapicodegen.Created{Id: resDto.Id}

	// Postcondition
	if err := ctx.Validate(res); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, res)
}
