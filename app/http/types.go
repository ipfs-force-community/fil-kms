package http

import "github.com/gin-gonic/gin"

type Handler interface {
	CheckParams(*gin.Context)
	Handler(*gin.Context)
}

type HandlerFactory func() Handler