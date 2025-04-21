package handler

import (
	"fmt"
	"medods_hire_me/internal/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) issueToken(c *gin.Context) {
	var request utils.TokenRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
		return
	}

	clientIP := c.GetHeader("X-Real-IP")
	if clientIP == "" {
		clientIP = c.ClientIP()
	}
	tokens, err := h.services.IssueToken(request.UserID.String(), clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to issue tokens: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) refreshToken(c *gin.Context) {
	var request utils.RefreshRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
		return
	}
	clientIP := c.GetHeader("X-Real-IP")
	if clientIP == "" {
		clientIP = c.ClientIP()
	}

	tokens, err := h.services.RefreshToken(request.RefreshToken, request.AccessToken, clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to refresh tokens: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) setEmail(c *gin.Context) {
	var request utils.MailRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %s", err.Error())})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Invalid request body")})
		return
	}
	err := h.services.SetEmail(userID.(string), request.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Fail")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok"})
}
