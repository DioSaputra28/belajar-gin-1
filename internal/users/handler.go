package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUsers(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	FindUserById(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userHandler struct {
	svc UserService
}

func NewUserHandler(svc UserService) UserHandler {
	return &userHandler{svc: svc}
}

func (h *userHandler) GetUsers(c *gin.Context) {
	page, ok := c.GetQuery("page")
	if !ok {
		page = "1"
	}

	limit, ok := c.GetQuery("limit")
	if !ok {
		limit = "10"
	}

	search, ok := c.GetQuery("search")
	if !ok {
		search = ""
	}

	intPage, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	users, err := h.svc.GetUsers(intPage, intLimit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Users found successfully",
		"data":    users,
	})
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var user CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_db, err := h.svc.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user_db,
	})
}

func (h *userHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user UpdateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	err = h.svc.UpdateUser(uint(intId), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    user,
	})
}


func (h *userHandler) FindUserById(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	user, err := h.svc.FindUserById(uint(intId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User found successfully",
		"data":    user,
	})
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	err = h.svc.DeleteUser(uint(intId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"data":    "User deleted successfully",
	})
}
