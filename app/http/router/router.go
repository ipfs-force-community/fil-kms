package router

import (
	"fil-kms/app/http/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	engine := gin.Default()

	engine.GET("/walletList", Handle(controller.NewOnWalletList))
	engine.POST("/sign", Handle(controller.NewOnWalletSign))
	engine.POST("/local", Handle(controller.NewOnLocalReq))
	return engine
}
