package zap_mate

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//var Point = make(map[string]int)

func (zml *ZapMateLogger) write(lvl zapcore.Level, msg string, fields ...zap.Field) {
	if ce := zml.Logger.Check(lvl, msg); ce != nil {
		if zml.isAsync { //if isAsync so...
			//if ce := zml.Logger.Check(lvl, msg); ce != nil { //异步读取info可能会导致时间对不上， 在日志量非常大的情况下可能会产生微秒级别的误差
			//ce.Write(info.fields...) //因为他显示时间是根绝这个ce来获取时间，但是问题不大  都是微妙级别的，  //如果要求准确性非常高的话可以考虑传递 ce,但是那样开销比传递info大
			//看来这里不能使用check了这个应该放在  应用端，write只放写操作，不然配置了info 多余的debug都会过来，而且时间可能对不上，更会造成一个bug,就是named会不起作用
			le := logMsgPool.Get().(*logEntry)
			le.entry = ce
			le.fields = fields
			//Point[fmt.Sprintf("%p", le)] = Point[fmt.Sprintf("%p", le)] + 1  //print sync.pool create point count!
			zml.entryChan <- le
			zml.wg.Add(1)
			//}
		} else { //if isSync so...
			//if ce := zml.Logger.Check(lvl, msg); ce != nil {
			ce.Write(fields...)
			//}
		}
	}
}

func (zml *ZapMateLogger) asyncWrite() {
	select {
	case le := <-zml.entryChan:
		le.entry.Write(le.fields...) //因为他显示时间是根绝这个ce来获取时间，但是问题不大  都是微妙级别的，  //如果要求准确性非常高的话可以考虑传递 ce,但是那样开销比传递info大
		//看来这里不能使用check了这个应该放在  应用端，write只放写操作，不然配置了info 多余的debug都会过来，而且时间可能对不上，更会造成一个bug,就是named会不起作用
		logMsgPool.Put(le) //获取之后需要put   不塞回去，那么就跟没使用 sync.Pool一样   //不设置put  传递指针  和  不使用pool一样都是 17次，而 设置put 只进行了13次 GC
		zml.wg.Done()
	}
}

//
//func (zml *ZapMateLogger) asyncWrite() {
//	select {
//	case info := <-zml.msgChan:
//		zml.wg.Done()
//		log.Printf("----%p", zml.asynclogger)
//		if ce := zml.asynclogger.Check(info.lv, info.msg); ce != nil { //异步读取info可能会导致时间对不上， 在日志量非常大的情况下可能会产生微秒级别的误差
//			ce.Write(info.fields...) //因为他显示时间是根绝这个ce来获取时间，但是问题不大  都是微妙级别的，  //如果要求准确性非常高的话可以考虑传递 ce,但是那样开销比传递info大
//			//看来这里不能使用check了这个应该放在  应用端，write只放写操作，不然配置了info 多余的debug都会过来，而且时间可能对不上，更会造成一个bug,就是named会不起作用
//		}
//		//logMsgPool.Put(info) //获取之后需要put   不塞回去，那么就跟没使用 sync.Pool一样   //不设置put  传递指针  和  不使用pool一样都是 17次，而 设置put 只进行了13次 GC
//	}
//}
