package addresses

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddressHandler interface {
	CreateAddress(c *gin.Context)
	GetAddresses(c *gin.Context)
	UpdateAddress(c *gin.Context)
	FindAddressById(c *gin.Context)
	DeleteAddress(c *gin.Context)
}

type addressHandler struct {
	svc AddressService
}

func NewAddressHandler(svc AddressService) AddressHandler {
	return &addressHandler{svc: svc}
}

func (h *addressHandler) CreateAddress(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request CreateAddressRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.svc.CreateAddress(user_id.(uint), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Address created successfully",
		"data":    response,
	})
}

func (h *addressHandler) GetAddresses(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	contact_id := c.Query("contact_id")
	if contact_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contact_id is required"})
		return
	}

	intContactId, err := strconv.Atoi(contact_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact_id"})
		return
	}

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	search := c.DefaultQuery("search", "")

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

	response, err := h.svc.GetAddresses(user_id.(uint), uint(intContactId), intPage, intLimit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Addresses retrieved successfully",
		"data":    response,
	})
}

func (h *addressHandler) FindAddressById(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	address, err := h.svc.FindAddressById(user_id.(uint), uint(intId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Address found successfully",
		"data":    address,
	})
}

func (h *addressHandler) UpdateAddress(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var request UpdateAddressRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.svc.UpdateAddress(user_id.(uint), uint(intId), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Address updated successfully",
		"data":    response,
	})
}

func (h *addressHandler) DeleteAddress(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err = h.svc.DeleteAddress(user_id.(uint), uint(intId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Address deleted successfully",
	})
}
