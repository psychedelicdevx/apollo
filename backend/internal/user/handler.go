package user

import (
	"apollo/internal/platform/auth"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/auth")
	g.POST("/register", h.register)
	g.POST("/login", h.login)
	g.GET("/me", auth.Middleware(), h.me)
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) register(c *gin.Context) {
	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.service.Register(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email})
}

func (h *Handler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) QuotaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		if err := h.service.ConsumeSearch(userID); err != nil {
			if errors.Is(err, ErrQuotaExceeded) {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "quota check failed"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (h *Handler) me(c *gin.Context) {
	userID := c.GetString("userID")

	u, err := h.service.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            u.ID,
		"email":         u.Email,
		"role":          u.Role,
		"daily_limit":   u.DailyLimit,
		"searches_used": u.SearchesUsed,
		"created_at":    u.CreatedAt,
	})
}
