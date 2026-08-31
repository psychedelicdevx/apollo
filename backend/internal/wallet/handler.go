package wallet

import (
	"errors"
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

func (h *Handler) RegisterRoutes(r *gin.Engine, quota gin.HandlerFunc) {
	g := r.Group("/wallets", auth.Middleware())
	g.GET("", h.history)
	g.GET("/:address", quota, h.lookup)
	g.GET("/:address/transactions", h.transactions)
	g.GET("/:address/tokens", h.tokens)
	g.GET("/:address/overview", h.overview)
	g.GET("/:address/graph", h.graph)
}

func (h *Handler) lookup(c *gin.Context) {
	userID := c.GetString("userID")
	address := c.Param("address")

	snap, err := h.service.Lookup(userID, address)
	if err != nil {
		if errors.Is(err, ErrInvalidAddress) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "could not fetch wallet"})
		return
	}
	c.JSON(http.StatusOK, snap)
}

func (h *Handler) history(c *gin.Context) {
	userID := c.GetString("userID")

	snaps, err := h.service.History(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load history"})
		return
	}
	c.JSON(http.StatusOK, snaps)
}

func (h *Handler) transactions(c *gin.Context) {
	address := c.Param("address")

	txs, err := h.service.Transactions(address)
	if err != nil {
		if errors.Is(err, ErrInvalidAddress) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "could not fetch transactions"})
		return
	}
	c.JSON(http.StatusOK, txs)
}

func (h *Handler) graph(c *gin.Context) {
	address := c.Param("address")

	g, err := h.service.Graph(address)
	if err != nil {
		if errors.Is(err, ErrInvalidAddress) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "could not build graph"})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *Handler) overview(c *gin.Context) {
	address := c.Param("address")

	ov, err := h.service.Overview(address)
	if err != nil {
		if errors.Is(err, ErrInvalidAddress) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "could not build overview"})
		return
	}
	c.JSON(http.StatusOK, ov)
}

func (h *Handler) tokens(c *gin.Context) {
	address := c.Param("address")
	refresh := c.Query("refresh") == "true"

	portfolio, err := h.service.TokenHoldings(address, refresh)
	if err != nil {
		if errors.Is(err, ErrInvalidAddress) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "could not fetch tokens"})
		return
	}
	c.JSON(http.StatusOK, portfolio)
}
