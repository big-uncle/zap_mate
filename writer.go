package zap_mate

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (zml *ZapMateLogger) write(lvl zapcore.Level, msg string, fields ...zap.Field) {
	if ce := zml.Logger.Check(lvl, msg); ce != nil { //The time of the log is generated here!
		if zml.isAsync { //if isAsync so...
			le := logMsgPool.Get().(*logEntry)
			le.entry = ce
			le.fields = fields
			zml.wg.Add(1)
			zml.entryChan <- le
		} else { //if isSync so...
			ce.Write(fields...)
		}
	}
}

func (zml *ZapMateLogger) asyncWrite() {
	for {
		select {
		case le := <-zml.entryChan:
			le.entry.Write(le.fields...)
			logMsgPool.Put(le) //There must used Sync.Put,to avoid a lot of GC
			zml.wg.Done()
		}
	}
}

func (zml *ZapMateLogger) flushWrite() {
	for {
		select {
		case le := <-zml.entryChan:
			le.entry.Write(le.fields...)
			logMsgPool.Put(le) //There must used Sync.Put,to avoid a lot of GC
			zml.wg.Done()
		default: //if there is no data in the channel,so finsh sync flush and exit
			return
		}
	}
}
