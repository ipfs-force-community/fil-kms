package controller

import (
	"context"
	"errors"

	"github.com/filecoin-project/go-address"
	"github.com/gin-gonic/gin"

	"fil-kms/app/config"
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
	cfg, ok := gCtx.Get("config")
	if !ok {
		utils.Error(gCtx, http_response.FAIL, errors.New("not found config in gin context"))
		return
	}
	addrs, err := sign_service.GlobalWalletService.WalletList(context.TODO())
	if err != nil {
		utils.Error(gCtx, http_response.FAIL, err)
		return
	}
	uniqAddrs := make(map[address.Address]struct{}, len(addrs))
	for _, addr := range addrs {
		uniqAddrs[addr] = struct{}{}
	}
	var out []address.Address
	for _, f := range cfg.(*config.Config).Filters {
		_, ok := uniqAddrs[f.Client.Address()]
		if ok {
			out = append(out, f.Client.Address())
		}
	}

	utils.Success(gCtx, out)
}
