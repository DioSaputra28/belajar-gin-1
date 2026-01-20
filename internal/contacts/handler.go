package contacts

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactHandler interface {
	GetContacts(c *gin.Context)
	CreateContact(c *gin.Context)
	FindContactById(c *gin.Context)
	UpdateContact(c *gin.Context)
	DeleteContact(c *gin.Context)
}

type contactHandler struct {
	svc ContactService
}

func NewContactHandler(svc ContactService) ContactHandler {
	return &contactHandler{svc: svc}
}

func (h *contactHandler) GetContacts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
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

	response, err := h.svc.GetContacts(intPage, intLimit, user_id.(int), search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contacts retrieved successfully",
		"data":    response,
	})
}


func (h *contactHandler) CreateContact(c *gin.Context) {
	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contact_db, err := h.svc.CreateContact(contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Contact created successfully",
		"data":    contact_db,
	})
}

func (h *contactHandler) FindContactById(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	contact, err := h.svc.FindContactById(uint(intId), user_id.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact found successfully",
		"data":    contact,
	})
}

func (h *contactHandler) UpdateContact(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	var contact Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.UpdateContact(uint(intId), user_id.(uint), contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact updated successfully",
		"data":    contact,
	})
}

func (h *contactHandler) DeleteContact(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	err = h.svc.DeleteContact(uint(intId), user_id.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact deleted successfully",
	})
}
