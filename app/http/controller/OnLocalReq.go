package controller

import (
	"context"
	"encoding/json"
	"fil-kms/app/global/http_response"
	"fil-kms/app/http"
	"fil-kms/app/http/controller/utils"
	"fil-kms/app/service/sign_service"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/gin-gonic/gin"
)

type onLocalReq struct {
	Method string `json:"method" form:"method" binding:"required"`
	Params []byte `json:"params" form:"params" binding:"required"`
}

func NewOnLocalReq() http.Handler {
	return &onLocalReq{}
}

func (received *onLocalReq) CheckParams(gCtx *gin.Context) {
	remoteAddr := gCtx.Request.RemoteAddr
	if err := gCtx.BindJSON(&received); err != nil {
		log.Warnf("remote: %v,err: %v", remoteAddr, err)
		utils.Error(gCtx, http_response.ValidatorParamsCheckFail, err)
		return
	}
}

func (received *onLocalReq) Handler(gCtx *gin.Context) {
	method, ok := MethodMap[received.Method]
	if !ok {
		utils.Error(gCtx, http_response.UnkownLocalMethod, nil)
		return
	}
	result, err := method(received.Params)
	if err != nil {
		utils.Error(gCtx, http_response.FAIL, err)
		return
	}
	utils.Success(gCtx, result)
	return
}

func walletImport(params []byte) ([]byte, error) {
	var ki types.KeyInfo

	err := json.Unmarshal(params, &ki)
	if err != nil {
		return nil, err
	}
	addr, err := sign_service.GlobalWalletService.WalletImport(context.TODO(), &ki)
	if err != nil {
		return nil, err
	}

	return addr.Bytes(), nil
}

func walletDelete(params []byte) ([]byte, error) {
	addr, err := address.NewFromBytes(params)
	if err != nil {
		return nil, err
	}

	err = sign_service.GlobalWalletService.WalletDelete(context.TODO(), addr)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func walletList(params []byte) ([]byte, error) {
	addrList, err := sign_service.GlobalWalletService.WalletList(context.TODO())
	if err != nil {
		return nil, err
	}

	ret, err := json.Marshal(addrList)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

var MethodMap = make(map[string]func([]byte) ([]byte, error))

func init() {
	MethodMap["WalletImport"] = walletImport
	MethodMap["WalletDelete"] = walletDelete
	MethodMap["WalletList"] = walletList
}
