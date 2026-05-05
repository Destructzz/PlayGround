package handlers

import (
	"backend/internal/http/middleware"
	"backend/internal/repo/sqlc"
	"backend/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ToolHandler struct {
	q *sqlc.Queries
}

func NewToolHandler(q *sqlc.Queries) *ToolHandler {
	return &ToolHandler{q: q}
}

// @Summary      Get Tool
// @Description  Get tool
// @Tags         Tool
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /docs [get]
func (h *ToolHandler) Docs(c *gin.Context) {
	user, _, ok := middleware.ResolveSession(c, h.q)
	if !ok || user.Role != sqlc.RoleAdmin {
		c.Redirect(http.StatusMovedPermanently, pkg.GetEnv("FRONTEND_URL", "http://localhost:3000"))
		return
	}

	c.File("./static/scalar.html")
}

// @Summary      Get OpenAPI Spec
// @Description  Get openapi spec
// @Tags         Tool
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /openapi.json [get]
func (h *ToolHandler) GetOpenAPI(c *gin.Context) {
	user, _, ok := middleware.ResolveSession(c, h.q)
	if !ok || user.Role != sqlc.RoleAdmin {
		c.Redirect(http.StatusMovedPermanently, pkg.GetEnv("FRONTEND_URL", "http://localhost:3000"))
		return
	}

	c.File("./docs/swagger.json")
}

