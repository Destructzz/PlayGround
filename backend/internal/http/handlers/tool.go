package handlers

import (
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

// @Summary      Get Tool Docs
// @Description  Returns the Scalar documentation page for admin users
// @Tags         Tool
// @Produce      text/html
// @Success      200  {file}    file
// @Failure      301  {string}  string  "Redirect to frontend"
// @Router       /docs [get]
func (h *ToolHandler) Docs(c *gin.Context) {
	user, _, ok := pkg.ResolveSession(c, h.q)
	if !ok || user.Role != sqlc.RoleAdmin {
		c.Redirect(http.StatusMovedPermanently, pkg.GetEnv("FRONTEND_URL", "http://localhost:3000"))
		return
	}

	c.File("./static/scalar.html")
}

// @Summary      Get OpenAPI Spec
// @Description  Returns the generated OpenAPI specification for admin users
// @Tags         Tool
// @Produce      json
// @Success      200  {file}    file
// @Failure      301  {string}  string  "Redirect to frontend"
// @Router       /openapi.json [get]
func (h *ToolHandler) GetOpenAPI(c *gin.Context) {
	user, _, ok := pkg.ResolveSession(c, h.q)
	if !ok || user.Role != sqlc.RoleAdmin {
		c.Redirect(http.StatusMovedPermanently, pkg.GetEnv("FRONTEND_URL", "http://localhost:3000"))
		return
	}

	c.File("./docs/swagger.json")
}
