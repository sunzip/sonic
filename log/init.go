package log

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"

	"github.com/go-sonic/sonic/config"
)

func NewLogger(conf *config.Config) *zap.Logger {
	_, err := os.Stat(conf.Sonic.LogDir)
	if err != nil {
		if os.IsNotExist(err) && config.LogTo(config.File) {
			err := os.MkdirAll(conf.Sonic.LogDir, os.ModePerm)
			if err != nil {
				panic("mkdir failed![%v]")
			}
		}
	}

	var core zapcore.Core
	var cores []zapcore.Core

	if config.LogTo(config.Console) {
		core := zapcore.NewCore(getDevEncoder(), os.Stdout, getLogLevel(conf.Log.Levels.App))
		cores = append(cores, core)
	}
	if config.LogTo(config.File) {
		core = zapcore.NewCore(getProdEncoder(), getWriter(conf), zap.DebugLevel)
		cores = append(cores, core)
	}
	core = zapcore.NewTee(cores...)

	// 传入 zap.AddCaller() 显示打日志点的文件名和行数
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))

	exportUseLogger = logger.WithOptions(zap.AddCallerSkip(1))
	exportUseSugarLogger = exportUseLogger.Sugar()
	return logger
}

// getWriter 自定义Writer,分割日志
func getWriter(conf *config.Config) zapcore.WriteSyncer {
	rotatingLogger := &lumberjack.Logger{
		Filename: filepath.Join(conf.Sonic.LogDir, conf.Log.FileName),
		MaxSize:  conf.Log.MaxSize,
		MaxAge:   conf.Log.MaxAge,
		Compress: conf.Log.Compress,
	}
	l := FileLoggerBaselumberjack{rotatingLogger: rotatingLogger}
	l.RegisteBeforeWrite(RemoveColor)
	return zapcore.AddSync(l)
}

// 写文件，移除打印文字样式
func RemoveColor(p []byte) []byte {
	str := string(p)
	str = strings.ReplaceAll(str, logger.Reset, "")
	str = strings.ReplaceAll(str, logger.Red, "")
	str = strings.ReplaceAll(str, logger.Green, "")
	str = strings.ReplaceAll(str, logger.Yellow, "")
	str = strings.ReplaceAll(str, logger.Blue, "")
	str = strings.ReplaceAll(str, logger.Magenta, "")
	str = strings.ReplaceAll(str, logger.Cyan, "")
	str = strings.ReplaceAll(str, logger.White, "")
	str = strings.ReplaceAll(str, logger.BlueBold, "")
	str = strings.ReplaceAll(str, logger.MagentaBold, "")
	str = strings.ReplaceAll(str, logger.RedBold, "")
	str = strings.ReplaceAll(str, logger.YellowBold, "")

	return []byte(str)
}

type FileLoggerBaselumberjack struct {
	rotatingLogger *lumberjack.Logger
	// RegisteBeforeWrite
	f func([]byte) []byte
}

func (t FileLoggerBaselumberjack) Write(p []byte) (n int, err error) {
	if t.f != nil {
		p = t.f(p)
	}
	return t.rotatingLogger.Write(p)
}

func (t *FileLoggerBaselumberjack) RegisteBeforeWrite(f func([]byte) []byte) {
	t.f = f
}

// getProdEncoder 自定义日志编码器
func getProdEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getDevEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		panic("log level error")
	}
}
