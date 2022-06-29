package v1

import "github.com/gin-gonic/gin"

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Register(api *gin.RouterGroup) {
}
