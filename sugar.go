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
	_multipleErrMsg     = "Multiple errors without a key."
)

// Note: base.Logger and sugar.Logger must be same pointer
type MateSugaredLogger struct {
	base               *ZapMateLogger //base is core!
	*zap.SugaredLogger                //
}

func (s *MateSugaredLogger) Desugar() *ZapMateLogger {
	base := s.base.clone()
	base.Logger = s.SugaredLogger.Desugar()
	return base

}

// Note: func setAsync is must be setting on the root node, Otherwise it will cause other errors!
// Child node cannot affect parent nodes,but child node all feature of extends parent node!
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

func (s *MateSugaredLogger) AsyncDebug(args ...interface{}) {
	s.asynclog(zap.DebugLevel, "", args, nil)
}

func (s *MateSugaredLogger) AsyncInfo(args ...interface{}) {
	s.asynclog(zap.InfoLevel, "", args, nil)
}

func (s *MateSugaredLogger) AsyncWarn(args ...interface{}) {
	s.asynclog(zap.WarnLevel, "", args, nil)
}

func (s *MateSugaredLogger) AsyncError(args ...interface{}) {
	s.asynclog(zap.ErrorLevel, "", args, nil)
}

func (s *MateSugaredLogger) AsyncDPanic(args ...interface{}) {
	s.asynclog(zap.DPanicLevel, "", args, nil)
}

func (s *MateSugaredLogger) AsyncPanic(args ...interface{}) {
	s.asynclog(zap.PanicLevel, "", args, nil)
}

func (s *MateSugaredLogger) AsyncFatal(args ...interface{}) {
	s.asynclog(zap.FatalLevel, "", args, nil)
}

func (s *MateSugaredLogger) AsyncDebugf(template string, args ...interface{}) {
	s.asynclog(zap.DebugLevel, template, args, nil)
}

func (s *MateSugaredLogger) AsyncInfof(template string, args ...interface{}) {
	s.asynclog(zap.InfoLevel, template, args, nil)
}

func (s *MateSugaredLogger) AsyncWarnf(template string, args ...interface{}) {
	s.asynclog(zap.WarnLevel, template, args, nil)
}

func (s *MateSugaredLogger) AsyncErrorf(template string, args ...interface{}) {
	s.asynclog(zap.ErrorLevel, template, args, nil)
}

func (s *MateSugaredLogger) AsyncDPanicf(template string, args ...interface{}) {
	s.asynclog(zap.DPanicLevel, template, args, nil)
}

func (s *MateSugaredLogger) AsyncPanicf(template string, args ...interface{}) {
	s.asynclog(zap.PanicLevel, template, args, nil)
}

func (s *MateSugaredLogger) AsyncFatalf(template string, args ...interface{}) {
	s.asynclog(zap.FatalLevel, template, args, nil)
}

func (s *MateSugaredLogger) AsyncDebugw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.DebugLevel, msg, nil, keysAndValues)
}

func (s *MateSugaredLogger) AsyncInfow(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.InfoLevel, msg, nil, keysAndValues)
}

func (s *MateSugaredLogger) AsyncWarnw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.WarnLevel, msg, nil, keysAndValues)
}

func (s *MateSugaredLogger) AsyncErrorw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.ErrorLevel, msg, nil, keysAndValues)
}

func (s *MateSugaredLogger) AsyncDPanicw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.DPanicLevel, msg, nil, keysAndValues)
}

func (s *MateSugaredLogger) AsyncPanicw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.PanicLevel, msg, nil, keysAndValues)
}

func (s *MateSugaredLogger) AsyncFatalw(msg string, keysAndValues ...interface{}) {
	s.asynclog(zap.FatalLevel, msg, nil, keysAndValues)
}

func (s *MateSugaredLogger) Flush() error {
	return s.base.Flush()
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

	s.base.write(lvl, msg, s.sweetenFields(context)...)
}

func (s *MateSugaredLogger) sweetenFields(args []interface{}) []zap.Field {
	if len(args) == 0 {
		return nil
	}

	var (
		// Allocate enough space for the worst case; if users pass only structured
		// fields, we shouldn't penalize them with extra allocations.
		fields    = make([]zap.Field, 0, len(args))
		invalid   invalidPairs
		seenError bool
	)

	for i := 0; i < len(args); {
		// This is a strongly-typed field. Consume it and move on.
		if f, ok := args[i].(zap.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		// If it is an error, consume it and move on.
		if err, ok := args[i].(error); ok {
			if !seenError {
				seenError = true
				fields = append(fields, zap.Error(err))
			} else {
				s.base.Error(_multipleErrMsg, zap.Error(err))
			}
			i++
			continue
		}

		// Make sure this element isn't a dangling key.
		if i == len(args)-1 {
			s.base.Error(_oddNumberErrMsg, zap.Any("ignored", args[i]))
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string, add this pair to the slice of invalid pairs.
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			// Subsequent errors are likely, so allocate once up front.
			if cap(invalid) == 0 {
				invalid = make(invalidPairs, 0, len(args)/2)
			}
			invalid = append(invalid, invalidPair{i, key, val})
		} else {
			fields = append(fields, zap.Any(keyStr, val))
		}
		i += 2
	}

	// If we encountered any invalid key-value pairs, log an error.
	if len(invalid) > 0 {
		s.base.Error(_nonStringKeyErrMsg, zap.Array("invalid", invalid))
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
