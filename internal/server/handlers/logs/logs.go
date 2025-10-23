package logs

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"backendT/internal/database/repository"
)

type Repo interface {
	LogsGetAll(ctx context.Context) ([]repository.LogsGetAllRow, error)
	LogsGetBasicViewWithOffsetLimit(ctx context.Context, params repository.LogsGetBasicViewWithOffsetLimitParams) ([]repository.LogsGetBasicViewWithOffsetLimitRow, error)
	LogsGetBasicViewWithOffsetLimitAdvanced(ctx context.Context, params repository.LogsGetBasicViewWithOffsetLimitAdvancedParams) ([]repository.LogsGetBasicViewWithOffsetLimitAdvancedRow, error)
}

type LogsHandler struct {
	repo Repo
}

func NewLogsHandler(r *repository.Queries) *LogsHandler {
	return &LogsHandler{
		repo: r,
	}
}

// GetAllLogs handles HTTP GET requests to retrieve all logs.
// @Summary Get all logs
// @Description Returns a list of all logs from the database.
// @Tags logs
// @Produce json
// @Success 200 {array} repository.Log "List of logs"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /logs [get]
func (h *LogsHandler) GetAllLogs(c echo.Context) error {
	logs, err := h.repo.LogsGetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch logs",
		})
	}
	return c.JSON(http.StatusOK, logs)
}

// GetLogsWithPagination handles HTTP GET requests to retrieve paginated logs.
// @Summary Get paginated logs without filters
// @Description Returns a paginated list of logs with basic view
// @Tags logs
// @Produce json
// @Param offset query int false "Offset for pagination"
// @Param limit query int false "Limit for pagination"
// @Success 200 {array} repository.LogsGetBasicViewWithOffsetLimitRow
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /logs/paginated [get]
func (h *LogsHandler) GetLogsWithPagination(c echo.Context) error {
	var params repository.LogsGetBasicViewWithOffsetLimitParams

	// Parse query parameters with defaults
	offset := c.QueryParam("offset")
	if offset != "" {
		offsetInt, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid offset parameter",
			})
		}
		params.Limit = offsetInt
	}

	limit := c.QueryParam("limit")
	if limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid limit parameter",
			})
		}
		params.Begining = limitInt
	}

	logs, err := h.repo.LogsGetBasicViewWithOffsetLimit(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch logs",
		})
	}

	return c.JSON(http.StatusOK, logs)
}

// GetLogsAdvanced handles HTTP GET requests to retrieve filtered logs.
// @Summary Get filtered logs
// @Description Returns filtered logs based on method, response type and time range
// @Tags logs
// @Produce json
// @Param method query string false "HTTP method to filter by"
// @Param response query int false "Response status code to filter by"
// @Param timeRange query string false "Time range (e.g. '-1 hour', '-24 hours', '-7 days')"
// @Param offset query int false "Offset for pagination"
// @Param limit query int false "Limit for pagination"
// @Success 200 {array} repository.LogsGetBasicViewWithOffsetLimitAdvancedRow
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /logs/filtered [get]
func (h *LogsHandler) GetLogsAdvanced(c echo.Context) error {
	var params repository.LogsGetBasicViewWithOffsetLimitAdvancedParams

	// Parse method
	method := c.QueryParam("method")
	if method != "" {
		params.Method = sql.NullString{String: method, Valid: true}
	}

	// Parse response type
	responseType := c.QueryParam("response")
	if responseType != "" {
		responseInt, err := strconv.ParseInt(responseType, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid response type parameter",
			})
		}
		params.Responsetype = sql.NullInt64{Int64: responseInt, Valid: true}
	}

	// Parse time range
	timeRange := c.QueryParam("timeRange")

	if timeRange != "" {
		params.Timerange = sql.NullString{String: timeRange, Valid: true}
	}

	// Parse pagination parameters
	offset := c.QueryParam("offset")
	if offset != "" {
		offsetInt, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid offset parameter",
			})
		}
		params.Limit = offsetInt
	}

	limit := c.QueryParam("limit")
	if limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid limit parameter",
			})
		}
		params.Begining = limitInt
	}

	// if no filters are provided, return bad request
	if limit == "" || offset == "" || timeRange == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "At least limit, offset and timeRange parameters must be provided",
		})
	}

	logs, err := h.repo.LogsGetBasicViewWithOffsetLimitAdvanced(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch logs",
		})
	}

	return c.JSON(http.StatusOK, logs)
}
