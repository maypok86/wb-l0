package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maypok86/wb-l0/pkg/logger"
)

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c *gin.Context, code int, msg string) {
	logger.Error(msg)
	c.AbortWithStatusJSON(code, errorResponse{Error: msg})
}
