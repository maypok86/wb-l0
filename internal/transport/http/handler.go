package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maypok86/wb-l0/internal/entity"
	v1 "github.com/maypok86/wb-l0/internal/transport/http/v1"
)

type OrderUsecase interface {
	CreateOrder(context.Context, *entity.Order, time.Duration) (*entity.Order, error)
	GetOrderByID(context.Context, string) (*entity.Order, error)
	LoadDBToCache(context.Context, time.Duration) error
}

type Handler struct {
	orderUsecase OrderUsecase
}

func NewHandler(orderUsecase OrderUsecase) Handler {
	return Handler{
		orderUsecase: orderUsecase,
	}
}

func (h Handler) Init() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery(), gin.Logger())

	h.registerAPI(router)

	return router
}

func (h Handler) registerAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/healthcheck", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		v1Handler := v1.NewHandler(h.orderUsecase)
		v1Handler.Register(api)
	}
}
