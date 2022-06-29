package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/maypok86/wb-l0/internal/transport/http/v1"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) GetHTTPHandler() *gin.Engine {
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
		v1Handler := v1.NewHandler()
		v1Handler.Register(api)
	}
}
