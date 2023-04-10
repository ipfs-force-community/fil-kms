package variables

import (
	"go.uber.org/zap"
)

var (
	/*built-in*/
	ServerIP   string
	ServerPort int

	IsDebug bool
	Log     *zap.SugaredLogger // 全局日志指针
	AKS     string
)
