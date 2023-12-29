package handler

import (
	"CamMonitoring/src/_interface"
	"CamMonitoring/src/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler _interface.Handler[entity.User]

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.Repository.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Create(c *gin.Context) {
	var user entity.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := h.Repository.Add(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID})
}
