package zap_mate

import (
	"sync"

	"go.uber.org/zap"
)

type ZapMateLogger struct { //Async ZMLogger
	lock      *sync.Mutex
	isAsync   bool
	entryChan chan *logEntry
	//试过十万调数据 和 十万个大小的chan   使用结构体用栈分配耗时8~9ms,使用指针逃逸到堆上耗时12~13ms  ,而且 十万个分配在堆上,GC压力真的受不了   ;
	// 因为他这种是即时栈,所以一般不会造成栈溢出的现象，但是可能对内存要求有点高，但是速度绝对快的一批   一般来说  不可能单机在20ms内写10w条的数据的
	signalChan chan string
	chanLen    uint
	wg         *sync.WaitGroup
	*zap.Logger
}

func NewZapMateLogger(filename, section string) *ZapMateLogger {

	return &ZapMateLogger{
		lock:       new(sync.Mutex),
		Logger:     NewZapLogger(filename, section),
		entryChan:  make(chan *logEntry, 1),
		signalChan: make(chan string, 1),
		wg:         new(sync.WaitGroup),
	}
}

func (log *ZapMateLogger) With(fields ...zap.Field) *ZapMateLogger {
	copy := log.clone()
	copy.Logger = log.Logger.With(fields...)
	return copy
}

func (log *ZapMateLogger) Named(s string) *ZapMateLogger {
	copy := log.clone()
	copy.Logger = log.Logger.Named(s)
	return copy
}

func (log *ZapMateLogger) clone() *ZapMateLogger {
	copy := *log
	return &copy
}

func (log *ZapMateLogger) Sugar() *MateSugaredLogger {
	return &MateSugaredLogger{
		log.clone(), //要确保           核心 base  zap.Logger  和  sugar.logger 是同一个地址
		log.Logger.Sugar(),
	}
}
