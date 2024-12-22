// internal/adapters/controllers/order_controller.go
package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shayja/go-template-api/internal/entities"
	"github.com/shayja/go-template-api/internal/usecases"
	"github.com/shayja/go-template-api/internal/utils"
)

type OrderController struct {
	OrderInteractor usecases.OrderInteractor
}

// GetAll godoc
// @Summary      Get all orders
// @Description  Retrieve a paginated list of all orders
// @Tags         Orders
// @Param        page  query     int  true  "Page number"
// @Param        userid  query   string  true  "USer ID (uuid)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Router       /order [get]
// @Security apiKey
func (uc *OrderController) GetAll(c *gin.Context) {
	AddRequestHeader(c)
	page, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "Invalid page number"})
		return
	}

	user_id := c.Query("userid")
	if utils.IsValidUUID(user_id)==false {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "Invalid user id"})
		return
	}

	res, err := uc.OrderInteractor.GetAll(page, user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	if res != nil {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "msg": "No orders found for this page"})
	}
}

// GetById godoc
// @Summary      Get an order by ID
// @Description  Retrieve order details by order ID
// @Tags         Orders
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /order/{id} [get]
// @Security apiKey
func (uc *OrderController) GetById(c *gin.Context) {
	AddRequestHeader(c)

	var uri entities.IdRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	res, err := uc.OrderInteractor.GetById(uri.Id)
	if err != nil || !utils.IsValidUUID(res.Id) {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "msg": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
}

// Create godoc
// @Summary      Create a new order
// @Description  Add a new order
// @Tags         Orders
// @Param        order  body      entities.OrderRequest  true  "Order data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /order [post]
// @Security apiKey
func (uc *OrderController) Create(c *gin.Context) {
	AddRequestHeader(c)

	var post *entities.OrderRequest
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	insertedId, err := uc.OrderInteractor.Create(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "id": insertedId, "msg": nil})
}

// UpdateStatus godoc
// @Summary      Update order status
// @Description  Update the status of an order
// @Tags         Orders
// @Param        id      path      string  true  "Order ID"
// @Param        status  body      int     true  "New status"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Router       /order/{id}/status [put]
// @Security apiKey
func (uc *OrderController) UpdateStatus(c *gin.Context) {
	AddRequestHeader(c)

	var uri entities.IdRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	var status struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	res, err := uc.OrderInteractor.UpdateStatus(uri.Id, status.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
}
