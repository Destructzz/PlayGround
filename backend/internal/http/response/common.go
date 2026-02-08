package response

import (
	"backend/internal/http/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func withRequestID(c *gin.Context, payload gin.H) {
	if rid, ok := middleware.RequestIDFromContext(c); ok {
		payload["request_id"] = rid
	}
}

func Struct(c *gin.Context, status int, nameObj string, obj any){
	payload := gin.H{
		nameObj:   obj,
		"timestamp": time.Now().UTC(),
	}

	withRequestID(c, payload)
	c.JSON(http.StatusOK, payload)
}

type ParamResponse struct {
	Param     string    `json:"param" example:"123"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}
