package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handler) newOrderRoutes(v1 *gin.RouterGroup) {
	v1.GET("/order", h.getOrderByID)
}

func (h Handler) getOrderByID(c *gin.Context) {
	orderUID := c.Query("order_uid")
	if orderUID == "" {
		newErrorResponse(c, http.StatusBadRequest, "empty order_uid query param")
		return
	}
	order, err := h.orderUsecase.GetOrderByID(c.Request.Context(), orderUID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, order)
}
