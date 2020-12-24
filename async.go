package zap_mate

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logInfo struct { //Async msg
	msg    string
	fields []zap.Field //这块可能又会导致内存逃逸
	lv     zapcore.Level
}

func (zml *ZapMateLogger) AsyncDebug(msg string, fields ...zap.Field) {
	zml.write(zap.DebugLevel, msg, fields...)
}

func (zml *ZapMateLogger) AsyncInfo(msg string, fields ...zap.Field) {
	zml.write(zap.InfoLevel, msg, fields...)

}

func (zml *ZapMateLogger) AsyncWarn(msg string, fields ...zap.Field) {
	zml.write(zap.WarnLevel, msg, fields...)

}

func (zml *ZapMateLogger) AsyncError(msg string, fields ...zap.Field) {
	zml.write(zap.ErrorLevel, msg, fields...)

}

func (zml *ZapMateLogger) AsyncDPanic(msg string, fields ...zap.Field) {
	zml.write(zap.DPanicLevel, msg, fields...)

}

func (zml *ZapMateLogger) AsyncPanic(msg string, fields ...zap.Field) {
	zml.write(zap.PanicLevel, msg, fields...)

}

func (zml *ZapMateLogger) AsyncFatal(msg string, fields ...zap.Field) {
	zml.write(zap.FatalLevel, msg, fields...)

}

func (zml *ZapMateLogger) SetAsyncer(chanLen uint) *ZapMateLogger {
	zml.lock.Lock()
	defer zml.lock.Unlock()
	if zml.isAsync {
		return zml
	}
	zml.isAsync = true
	zml.chanLen = chanLen
	zml.msgChan = make(chan logInfo, zml.chanLen)
	go zml.startAsyncLogger()
	return zml
}

func (zml *ZapMateLogger) startAsyncLogger() {
	for {
		zml.asyncWrite()
	}
}

func (zml *ZapMateLogger) Flush() {
	defer zml.Sync()
	for len(zml.msgChan) > 0 {
		zml.asyncWrite()
	}
	zml.wg.Wait()
}
