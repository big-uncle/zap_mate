package zap_mate

import (
	"fmt"

	"go.uber.org/multierr"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	_oddNumberErrMsg    = "Ignored key without a value."
	_nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
)

type MateSugaredLogger struct {
	base               *ZapMateLogger //他还是以这个来进行运算的
	*zap.SugaredLogger                //这个只是  让matesugar可以使用 原生的方法罢了，  但需要这个里面嵌套的  logger 和 zapmatelogger 里的logger 时刻保证一致
}

func (s *MateSugaredLogger) Desugar() *ZapMateLogger {
	base := s.base.clone()
	base.Logger = s.SugaredLogger.Desugar()
	return base

}

func (s *MateSugaredLogger) Named(name string) *MateSugaredLogger {
	return &MateSugaredLogger{
		base:          s.base.Named(name),
		SugaredLogger: s.SugaredLogger.Named(name),
	}
}

func (s *MateSugaredLogger) With(args ...interface{}) *MateSugaredLogger {

	return &MateSugaredLogger{
		base:          s.base.With(s.sweetenFields(args)...),
		SugaredLogger: s.SugaredLogger.With(args...),
	}

}

// Debug uses fmt.Sprint to construct and log a message.
func (s *MateSugaredLogger) AsyncDebug(args ...interface{}) {
	s.asynclog(zap.DebugLevel, "", args, nil)
}

// Info uses fmt.Sprint to construct and log a message.
func (s *MateSugaredLogger) AsyncInfo(args ...interface{}) {
	s.asynclog(zap.InfoLevel, "", args, nil)
}

// Warn uses fmt.Sprint to construct and log a message.
func (s *MateSugaredLogger) AsyncWarn(args ...interface{}) {
	s.asynclog(zap.WarnLevel, "", args, nil)
}

// Error uses fmt.Sprint to construct and log a message.
func (s *MateSugaredLogger) AsyncError(args ...interface{}) {
	s.asynclog(zap.ErrorLevel, "", args, nil)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (s *MateSugaredLogger) AsyncDPanic(args ...interface{}) {
	s.asynclog(zap.DPanicLevel, "", args, nil)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (s *MateSugaredLogger) AsyncPanic(args ...interface{}) {
	s.asynclog(zap.PanicLevel, "", args, nil)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (s *MateSugaredLogger) AsyncFatal(args ...interface{}) {
	s.asynclog(zap.FatalLevel, "", args, nil)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (s *MateSugaredLogger) AsyncDebugf(template string, args ...interface{}) {
	s.asynclog(zap.DebugLevel, template, args, nil)
}

// Infof uses fmt.Sprintf to log a templated message.
func (s *MateSugaredLogger) AsyncInfof(template string, args ...interface{}) {
	s.asynclog(zap.InfoLevel, template, args, nil)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (s *MateSugaredLogger) AsyncWarnf(template string, args ...interface{}) {
	s.asynclog(zap.WarnLevel, template, args, nil)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (s *MateSugaredLogger) AsyncErrorf(template string, args ...interface{}) {
	s.asynclog(zap.ErrorLevel, template, args, nil)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (s *MateSugaredLogger) AsyncDPanicf(template string, args ...interface{}) {
	s.asynclog(zap.DPanicLevel, template, args, nil)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (s *MateSugaredLogger) AsyncPanicf(template string, args ...interface{}) {
	s.asynclog(zap.PanicLevel, template, args, nil)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (s *MateSugaredLogger) AsyncFatalf(template string, args ...interface{}) {
	s.asynclog(zap.FatalLevel, template, args, nil)
}

func (s *MateSugaredLogger) AsyncDebugw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.DebugLevel, msg, nil, keysAndValues)
}

func (s *MateSugaredLogger) AsyncInfow(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.InfoLevel, msg, nil, keysAndValues)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s *MateSugaredLogger) AsyncWarnw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.WarnLevel, msg, nil, keysAndValues)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s *MateSugaredLogger) AsyncErrorw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.ErrorLevel, msg, nil, keysAndValues)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func (s *MateSugaredLogger) AsyncDPanicw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.DPanicLevel, msg, nil, keysAndValues)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (s *MateSugaredLogger) AsyncPanicw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.PanicLevel, msg, nil, keysAndValues)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func (s *MateSugaredLogger) AsyncFatalw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.FatalLevel, msg, nil, keysAndValues)
}

// Sync flushes any buffered log entries.
func (s *MateSugaredLogger) Sync() error {
	return s.base.Sync()
}

func (s *MateSugaredLogger) asynclog(lvl zapcore.Level, template string, fmtArgs []interface{}, context []interface{}) {

	if lvl < zap.DPanicLevel && !s.base.Core().Enabled(lvl) {
		return
	}

	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	if ce := s.base.Check(lvl, msg); ce != nil {
		if s.base.isAsync {
			le := logMsgPool.Get().(*logEntry)
			le.entry = ce
			le.fields = s.sweetenFields(context)
			s.base.entryChan <- le
			s.base.wg.Add(1)
		} else {
			ce.Write(s.sweetenFields(context)...)
		}
	}
}

func (s *MateSugaredLogger) sweetenFields(args []interface{}) []zap.Field {
	if len(args) == 0 {
		return nil
	}

	fields := make([]zap.Field, 0, len(args))
	var invalid invalidPairs

	for i := 0; i < len(args); {
		if f, ok := args[i].(zap.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		if i == len(args)-1 {
			s.base.DPanic(_oddNumberErrMsg, zap.Any("ignored", args[i]))
			break
		}

		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			if cap(invalid) == 0 {
				invalid = make(invalidPairs, 0, len(args)/2)
			}
			invalid = append(invalid, invalidPair{i, key, val})
		} else {
			fields = append(fields, zap.Any(keyStr, val))
		}
		i += 2
	}

	if len(invalid) > 0 {
		s.base.DPanic(_nonStringKeyErrMsg, zap.Array("invalid", invalid))
	}
	return fields
}

type invalidPair struct {
	position   int
	key, value interface{}
}

type invalidPairs []invalidPair

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("position", int64(p.position))
	zap.Any("key", p.key).AddTo(enc)
	zap.Any("value", p.value).AddTo(enc)
	return nil
}

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	for i := range ps {
		err = multierr.Append(err, enc.AppendObject(ps[i]))
	}
	return err
}
