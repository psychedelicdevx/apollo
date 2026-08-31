package label

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"apollo/internal/platform/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/labels", auth.Middleware())
	g.GET("/:address", h.get)
}

func (h *Handler) get(c *gin.Context) {
	address := c.Param("address")

	l, err := h.service.Get(address)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no label for address"})
		return
	}
	c.JSON(http.StatusOK, l)
}
