package v2

import (
	"github.com/pkg/errors"
	"gitlab.oneitfarm.com/bifrost/cilog"
	"gitlab.oneitfarm.com/bifrost/cilog/redis_hook"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ciCore zapcore.Core

var (
	std          *Logger
	stdCallerFix *Logger
)

// Logger 实例
type Logger struct {
	*zap.SugaredLogger
	conf *Conf
}

// Conf 配置
type Conf struct {
	Caller     bool
	Debug      bool
	Level      zapcore.Level
	Encoding   string                 // json, console
	AppInfo    *cilog.ConfigAppData   // fixed fields
	HookConfig *redis_hook.HookConfig // set to nil if disabled
}

// Clone ...
func Clone(l *Logger) *Logger {
	c := *l.conf
	return &Logger{
		SugaredLogger: l.SugaredLogger,
		conf:          &c,
	}
}

// S 获取单例
func S() *Logger {
	return std
}

// GlobalConfig init
func GlobalConfig(conf Conf) error {
	c := conf
	l, err := newLogger(&c)
	if err != nil {
		return err
	}
	std = &Logger{
		SugaredLogger: l.Sugar(),
		conf:          &c,
	}
	stdCallerFix = &Logger{
		SugaredLogger: l.WithOptions(zap.AddCallerSkip(1)).Sugar(),
		conf:          &c,
	}
	return nil
}

func init() {
	l, _ := newLogger(&Conf{
		Level: zapcore.InfoLevel,
	})
	std = &Logger{
		SugaredLogger: l.Sugar(),
		conf:          &Conf{},
	}
	stdCallerFix = &Logger{
		SugaredLogger: l.WithOptions(zap.AddCallerSkip(1)).Sugar(),
		conf:          &Conf{},
	}
}

// NewZapLogger 创建自定义 Logger
func NewZapLogger(c *Conf) (l *zap.Logger, err error) {
	return newLogger(c)
}

func newLogger(c *Conf) (l *zap.Logger, err error) {
	conf := zap.NewProductionConfig()
	if c.Debug {
		conf = zap.NewDevelopmentConfig()
		conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	if c.Encoding != "" {
		conf.Encoding = c.Encoding
	} else {
		conf.Encoding = "console"
	}
	conf.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	conf.Level = zap.NewAtomicLevelAt(c.Level)
	if c.HookConfig != nil {
		hook, err := redis_hook.NewHook(*c.HookConfig)
		if err != nil {
			return nil, errors.Wrap(err, "hook init error")
		}
		_ciCore = NewCiCore(hook)
		fixedFields := getFixedFields(c.AppInfo)
		for k, v := range fixedFields {
			_ciCore = _ciCore.With([]zapcore.Field{zap.String(k, v)})
		}
		l, err = conf.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, _ciCore)
		}))
		if err != nil {
			return nil, errors.Wrap(err, "zap core init error")
		}
	} else {
		l, err = conf.Build()
	}
	if err != nil {
		return nil, errors.Wrap(err, "zap core init error")
	}
	l = l.WithOptions(zap.WithCaller(c.Caller), zap.AddStacktrace(zapcore.ErrorLevel))
	return
}
