package zap_mate

import (
	"sync"

	"go.uber.org/zap"
)

type ZapMateLogger struct { //Async ZMLogger
	lock      *sync.Mutex
	isAsync   bool
	entryChan chan *logEntry
	chanLen   uint
	wg        *sync.WaitGroup
	*zap.Logger
}

func NewZapMateLogger(filename, section string) *ZapMateLogger {

	return &ZapMateLogger{
		isAsync:   false,
		lock:      new(sync.Mutex),
		Logger:    NewZapLogger(filename, section),
		entryChan: make(chan *logEntry, 1),
		wg:        new(sync.WaitGroup),
	}
}

//Note: func setAsync is must be setting on the root node, Otherwise it will cause other errors!
//Child node cannot affect parent nodes,but child node all feature of extends parent node!
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
		log.clone(), // Note: base.Logger and sugar.Logger must be same pointer
		log.Logger.Sugar(),
	}
}
