package zap_mate

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//var Point = make(map[string]int)

func (zml *ZapMateLogger) write(lvl zapcore.Level, msg string, fields ...zap.Field) {
	if ce := zml.Logger.Check(lvl, msg); ce != nil { //The time of the log is generated here!
		if zml.isAsync { //if isAsync so...
			le := logMsgPool.Get().(*logEntry)
			le.entry = ce
			le.fields = fields
			//Point[fmt.Sprintf("%p", le)] = Point[fmt.Sprintf("%p", le)] + 1  //print sync.pool create point count!
			zml.entryChan <- le
			zml.wg.Add(1)
		} else { //if isSync so...s
			ce.Write(fields...)
		}
	}
}

func (zml *ZapMateLogger) asyncWrite() {
	select {
	case le := <-zml.entryChan:
		le.entry.Write(le.fields...)
		logMsgPool.Put(le) //There must used Sync.Put,to avoid a lot of GC
		zml.wg.Done()
	}
}
