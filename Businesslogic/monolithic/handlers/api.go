package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	 "karix.com/monolith/schemas"
	 "karix.com/monolith/service"
	

)

type HTTPHandler struct {
	svc *service.UserService
}

func NewHTTPHandler(svc *service.UserService) *HTTPHandler { return &HTTPHandler{svc: svc} }

func (h *HTTPHandler) CreateUser(c *gin.Context) {
	var req schemas.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.svc.CreateUser(c.Request.Context(), req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *HTTPHandler) GetUser(c *gin.Context) {
	idS := c.Param("id")
	id, err := strconv.ParseInt(idS, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	u, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if u == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}
