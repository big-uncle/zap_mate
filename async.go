package zap_mate

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//type logInfo struct { //Async msg
//	msg    string
//	fields []zap.Field //这块可能又会导致内存逃逸//判断是否有值 若无值，则不传递该值，大部分不应该有值的
//	lv     zapcore.Level
//}

type logEntry struct { //Async msg
	fields []zap.Field
	entry  *zapcore.CheckedEntry
}

var logMsgPool *sync.Pool

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
	zml.entryChan = make(chan *logEntry, zml.chanLen)
	logMsgPool = &sync.Pool{
		New: func() interface{} {
			return new(logEntry)
		},
	}
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
	for len(zml.entryChan) > 0 {
		zml.asyncWrite()
	}
	zml.wg.Wait()
}
