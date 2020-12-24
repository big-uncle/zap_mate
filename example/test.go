package main

import (
	"log"
	"sync"
	"time"

	"github.com/big-uncle/zap_mate"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

type logInfo struct {
	msg    string
	fields []zap.Field //这块可能又会导致内存逃逸
	lv     zapcore.Level
}

type ZMLogger struct { //Async ZMLogger
	lock    sync.Mutex
	isAsync bool
	msgChan chan logInfo //试过十万调数据 和 十万个大小的chan   使用结构体用栈分配耗时8~9ms,使用指针逃逸到堆上耗时12~13ms  ,而且 十万个分配在堆上,GC压力真的受不了   ;
	// 因为他这种是即时栈,所以一般不会造成栈溢出的现象，但是可能对内存要求有点高，但是速度绝对快的一批   一般来说  不可能单机在20ms内写10w条的数据的
	signalChan chan string
	ChanLen    uint
	wg         sync.WaitGroup
	*zap.Logger
}

var logger ZMLogger

func main() {

	start := time.Now()
	logger.Logger = zap_mate.NewLogger("./example/test.yaml", "default")
	logger.SetAsyncer(100000)
	var num int
	for num < 100000 {
		logger.AsyncInfo("ss")
		logger.wg.Add(1)
		//ZMLogger.Warn("I am zap_mate!")
		num++
	}
	//logger.wg.Wait()
	log.Printf("耗时【%v】", time.Since(start))
	logger.wg.Wait()
	//time.Sleep(100 * time.Second)
	// 首先定义俩个不同的变量 一一对应 logger  和 sugar   继承或者 嵌套都可以
	//	sugar:=ZMLogger.Sugar()
	//
	//	sugar.Info("")
	//	//s.log(InfoLevel, "", args, nil)
	//	sugar.Infof("")
	//	//s.log(InfoLevel, template, args, nil)
	//	sugar.Infow("")
	//sugar.
	//s.log(InfoLevel, msg, nil, keysAndValues)

}
func (log *ZMLogger) AsyncDebug(msg string, fields ...zap.Field) {
	//if ce := log.Check(zap.InfoLevel, msg); ce != nil {
	//	ce.Write(fields...)
	//}
	//log.Logger.Info(msg,fields...)  //造成了 额外的栈开销 不建议使用
	log.msgChan <- logInfo{msg, fields, zap.DebugLevel}
}
func (log *ZMLogger) AsyncInfo(msg string, fields ...zap.Field) {

	log.msgChan <- logInfo{msg, fields, zap.InfoLevel}
}

func (log *ZMLogger) AsyncWarn(msg string, fields ...zap.Field) {

	log.msgChan <- logInfo{msg, fields, zap.WarnLevel}
}

func (log *ZMLogger) AsyncError(msg string, fields ...zap.Field) {

	log.msgChan <- logInfo{msg, fields, zap.ErrorLevel}
}

func (log *ZMLogger) AsyncDPanic(msg string, fields ...zap.Field) {

	log.msgChan <- logInfo{msg, fields, zap.DPanicLevel}
}

func (log *ZMLogger) AsyncPanic(msg string, fields ...zap.Field) {

	log.msgChan <- logInfo{msg, fields, zap.PanicLevel}
}

func (log *ZMLogger) AsyncFatal(msg string, fields ...zap.Field) {
	log.msgChan <- logInfo{msg, fields, zap.FatalLevel}
}

func (log *ZMLogger) SetAsyncer(chanLen uint) {
	if log.isAsync {
		return
	}
	log.lock.Lock()
	defer log.lock.Unlock()
	log.isAsync = true
	log.ChanLen = chanLen
	log.msgChan = make(chan logInfo, log.ChanLen)
	go log.startSyncLogger()
}

func (log *ZMLogger) startSyncLogger() {
	for {
		select {
		case info := <-log.msgChan:
			logger.wg.Done()
			if ce := log.Check(info.lv, info.msg); ce != nil {
				ce.Write(info.fields...)
			}
		}
	}
}
