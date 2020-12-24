package zap_mate

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (zml *ZapMateLogger) write(lvl zapcore.Level, msg string, fields ...zap.Field) {
	if zml.isAsync { //if isAsync so...
		zml.wg.Add(1)
		zml.msgChan <- logInfo{msg, fields, lvl}
	} else { //if isSync so...
		if ce := zml.Check(lvl, msg); ce != nil {
			ce.Write(fields...)
		}
	}
}

func (zml *ZapMateLogger) asyncWrite() {
	select {
	case info := <-zml.msgChan:
		zml.wg.Done()
		if ce := zml.Check(info.lv, info.msg); ce != nil {
			ce.Write(info.fields...)
		}
	}
}
