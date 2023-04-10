package log

import (
	"go.uber.org/zap"
	"log"
)

func NewZapLog(isDebug bool) *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	if isDebug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	if logger, err := config.Build(zap.AddStacktrace(zap.DPanicLevel)); err == nil {
		return logger.Sugar()
	} else {
		log.Fatal("创建zap日志包失败，详情：" + err.Error())
	}
	return nil
}
