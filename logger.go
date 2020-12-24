package zap_mate

import (
	"sync"

	"go.uber.org/zap"
)

type ZapMateLogger struct { //Async ZMLogger
	lock    sync.Mutex
	isAsync bool
	msgChan chan logInfo //试过十万调数据 和 十万个大小的chan   使用结构体用栈分配耗时8~9ms,使用指针逃逸到堆上耗时12~13ms  ,而且 十万个分配在堆上,GC压力真的受不了   ;
	// 因为他这种是即时栈,所以一般不会造成栈溢出的现象，但是可能对内存要求有点高，但是速度绝对快的一批   一般来说  不可能单机在20ms内写10w条的数据的
	signalChan chan string
	chanLen    uint
	wg         sync.WaitGroup
	*zap.Logger
}

func NewZapMateLogger(filename, section string) *ZapMateLogger {

	return &ZapMateLogger{
		Logger:     NewZapLogger(filename, section),
		msgChan:    make(chan logInfo, 1),
		signalChan: make(chan string, 1),
	}
}
