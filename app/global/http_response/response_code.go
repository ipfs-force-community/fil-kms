package http_response

import "encoding/json"

type ResponseInfo struct {
	ResponseCode
	Data json.RawMessage `json:"data"`
}

type ResponseCode struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func newResponseCode(code int, msg string) ResponseCode {
	return ResponseCode{
		Code: code,
		Msg:  msg,
	}
}

var (
	OK    = newResponseCode(0, "成功")
	FAIL  = newResponseCode(-1, "失败")
	OK200 = newResponseCode(200, "成功")

	ValidatorParamsCheckFail = newResponseCode(-400100, "参数校验失败")
	UnkownLocalMethod        = newResponseCode(-401206, "未知本地方法名")
)

func IsOk(responseCode int) bool {
	switch responseCode {
	case OK.Code, OK200.Code:
		return true
	}
	return false
}
