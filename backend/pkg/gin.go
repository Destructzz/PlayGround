package pkg

import (
	"net/http"
	"strconv"

	"backend/internal/http/response"
	"github.com/gin-gonic/gin"
)

func SetGinMode(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
		return
	}

	gin.SetMode(gin.DebugMode)
}

func ParseIDParam(c *gin.Context) (int64, bool) {
	rawID := c.Param("id")
	if rawID == "" {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "can't take url parametrs", nil),
		).JSON(c)
		return 0, false
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("failed_param", "invalid id", nil),
		).JSON(c)
		return 0, false
	}

	return id, true
}
