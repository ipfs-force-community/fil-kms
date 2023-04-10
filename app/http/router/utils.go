package router

import (
	"fil-kms/app/http"
	"github.com/gin-gonic/gin"
)

func Handle(hf http.HandlerFactory) gin.HandlerFunc {
	return func(context *gin.Context) {
		handler := hf()
		handler.CheckParams(context)
		if context.IsAborted() {
			return
		}
		handler.Handler(context)
	}
}
