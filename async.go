package zap_mate

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

//Note: func setAsync is must be setting on the root node, Otherwise it will cause other errors!
//Child node cannot affect parent nodes,but child node all feature of extends parent node!
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

func (zml *ZapMateLogger) Flush() error {

	for len(zml.entryChan) > 0 {
		zml.asyncWrite()
	}
	zml.wg.Wait()
	return zml.Sync()
}
