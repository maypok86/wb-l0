package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/wb-l0/internal/entity"
)

type OrderUsecase interface {
	GetOrderByID(context.Context, string) (*entity.Order, error)
}

type Handler struct {
	orderUsecase OrderUsecase
}

func NewHandler(orderUsecase OrderUsecase) Handler {
	return Handler{
		orderUsecase: orderUsecase,
	}
}

func (h Handler) Register(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.newOrderRoutes(v1)
	}
}
