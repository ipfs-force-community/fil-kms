package utils

import (
	"errors"

	"fil-kms/app/global/http_response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReturnJson(Context *gin.Context, httpCode int, code http_response.ResponseCode, data interface{}) {
	Context.JSON(httpCode, gin.H{
		"code": code.Code,
		"msg":  code.Msg,
		"data": data,
	})
	Context.Abort()
}

func Success(c *gin.Context, data interface{}) {
	ReturnJson(c, http.StatusOK, http_response.OK, data)
}

func Error(c *gin.Context, codeMsg http_response.ResponseCode, err error) {
	var data string
	if err != nil {
		data = err.Error()
	}
	ReturnJson(c, http.StatusBadRequest, codeMsg, data)
}

func Error2(c *gin.Context, err *http_response.ResponseErr) {
	Error(c, err.ResponseCode, errors.New(err.Data))
}
