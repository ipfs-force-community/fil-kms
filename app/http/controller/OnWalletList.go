package controller

import (
	"context"

	"github.com/gin-gonic/gin"

	"fil-kms/app/global/http_response"
	"fil-kms/app/http"
	"fil-kms/app/http/controller/utils"
	"fil-kms/app/service/sign_service"
)

type onWalletList struct {
}

func NewOnWalletList() http.Handler {
	return &onWalletList{}
}

func (ws *onWalletList) CheckParams(gCtx *gin.Context) {

}

func (ws *onWalletList) Handler(gCtx *gin.Context) {
	addrs, err := sign_service.GlobalWalletService.WalletList(context.TODO())
	if err != nil {
		utils.Error(gCtx, http_response.FAIL, err)
		return
	}

	utils.Success(gCtx, addrs)
}
