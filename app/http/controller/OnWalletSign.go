package controller

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/builtin/v10/market"
	"github.com/filecoin-project/lotus/api"
	"github.com/gin-gonic/gin"

	"fil-kms/app/global/http_response"
	"fil-kms/app/global/variables"
	"fil-kms/app/http"
	"fil-kms/app/http/controller/utils"
	"fil-kms/app/service/sign_service"
	utils2 "fil-kms/app/utils"
)

type onWalletSign struct {
	Addr address.Address `json:"addr" form:"addr" binding:"required"`
	Msg  []byte          `json:"msg" form:"msg" binding:"required"`
}

func NewOnWalletSign() http.Handler {
	return &onWalletSign{}
}

func (ws *onWalletSign) CheckParams(gCtx *gin.Context) {
	remoteAddr := gCtx.Request.RemoteAddr

	if err := gCtx.BindJSON(ws); err != nil {
		log.Warnf("remote: %v,err: %v", remoteAddr, err)
		utils.Error(gCtx, http_response.ValidatorParamsCheckFail, err)
		return
	}

	authorization := gCtx.Request.Header.Get("Authorization")
	ok, sign := VerificationToken(authorization)
	if !ok {
		log.Warnf("illegal authorization info:%v", authorization)
		utils.Error(gCtx, http_response.FAIL, errors.New("authorization failed"))
		return
	}

	signB, err := hex.DecodeString(sign)
	if err != nil {
		log.Warnf("decode signature failed,sign:%v", sign)
		utils.Error(gCtx, http_response.FAIL, errors.New("decode signature failed"))
		return
	}

	requestJsonInfo, err := json.Marshal(ws)
	if err != nil {
		log.Errorf("marshal requestInfo err:%v", err)
		utils.Error(gCtx, http_response.FAIL, errors.New("marshal walletSign failed"))
		return
	}

	ok = utils2.Verify(requestJsonInfo, signB, []byte(variables.AKS))
	if !ok {
		log.Warnf("signature can't match")
		utils.Error(gCtx, http_response.FAIL, errors.New("authorization failed"))
		return
	}

	err = ws.check(gCtx)
	if err != nil {
		log.Warnf("The request is illegal. err:%v", err)
		utils.Error(gCtx, http_response.FAIL, err)
		return
	}

	return
}

func (ws *onWalletSign) check(gCtx *gin.Context) error {
	//检查订单信息
	{
		var proposal market.DealProposal
		err := (&proposal).UnmarshalCBOR(bytes.NewReader(ws.Msg))
		if err == nil {
			//TODO 在这里做一些额外的检查

			if !proposal.ClientCollateral.IsZero() {
				return errors.New("clientCollateral must be zero")
			}

			return nil
		}
	}

	return errors.New("illegal message information")
}

func (ws *onWalletSign) Handler(gCtx *gin.Context) {
	signature, err := sign_service.GlobalWalletService.WalletSign(context.TODO(), ws.Addr, ws.Msg, api.MsgMeta{})
	if err != nil {
		utils.Error(gCtx, http_response.FAIL, err)
		return
	}
	utils.Success(gCtx, *signature)
}

func VerificationToken(token string) (bool, string) {
	split := strings.Split(token, ":")
	if len(split) != 2 {
		return false, ""
	}

	sign := split[1]

	if len(sign) == 0 {
		return false, ""
	}

	return true, sign
}
