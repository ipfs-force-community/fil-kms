package router

import (
	"fil-kms/app/config"
	"fil-kms/app/http/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(cfg *config.Config) *gin.Engine {
	engine := gin.Default()

	engine.Use(func(ctx *gin.Context) {
		ctx.Set("config", cfg)
		ctx.Next()
	})

	engine.GET("/walletList", Handle(controller.NewOnWalletList))
	engine.POST("/sign", Handle(controller.NewOnWalletSign))
	engine.POST("/local", Handle(controller.NewOnLocalReq))
	return engine
}
