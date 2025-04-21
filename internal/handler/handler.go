package handler

import (
	"medods_hire_me/internal/service"

	"github.com/gin-gonic/gin"
	
)

type Handler struct {
	services *service.Service
}

func New(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()
	authApi := router.Group("/auth")
	{
		authApi.POST("/token", h.issueToken)
		authApi.POST("/refresh", h.refreshToken)
	}

	protected := router.Group("/protected")
	{
		protected.Use(h.JWTAuth())
		protected.PATCH("/mail", h.setEmail)
	}
	unsafeApi := router.Group("/unsafe")
	{
		unsafeApi.GET("/info", nil)
	}

	return router
}