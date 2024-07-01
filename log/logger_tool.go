package log

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 通过一个core， 生成logger
func NewLoggerWithCore(core zapcore.Core) *zap.Logger {
	// 传入 zap.AddCaller() 显示打日志点的文件名和行数
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))
	return logger
}

var consoleLogger, otherLogger *zap.Logger

// 根据系统core， 获取分开的console core|Logger、其他 core|Logger
// console core 只考虑一个
// 1. console core, 如果不存在， 则是nil
// 2. 是剩余的其他的core
func getCachedLogCores(core zapcore.Core) (*zap.Logger, *zap.Logger) {
	if consoleLogger != nil || otherLogger != nil {
		return consoleLogger, otherLogger
	} else {
		v := reflect.ValueOf(core)

		if v.Kind() == reflect.Slice {
			length := v.Len()
			for i := 0; i < length; i++ {
				element := v.Index(i)
				if element.Kind() == reflect.Interface {
					if !element.IsNil() {
						actualValue := element.Elem()
						inf := element.Interface()
						// if o, ok := inf.(*zapcore.Core); ok {
						// 	fmt.Println("OK1", o)
						// }
						// outField := element.FieldByName("out") //reflect: call of reflect.Value.FieldByName on interface Value
						outField := actualValue.Elem().FieldByName("out")
						// outField := actualValue.FieldByName("out") //reflect: call of reflect.Value.FieldByName on ptr Value
						outFieldActual := outField.Elem()
						fmt.Println("OK3", outFieldActual.Type().String())
						if outFieldActual.Type().String() == "*os.File" {
							fmt.Println("OK2", "console 日志")

							if o, ok := inf.(zapcore.Core); ok {
								fmt.Println("OK2", o)
								consoleLogger = NewLoggerWithCore(o)
							}
							// console 日志
						} else {
							// file日志
							fmt.Println("OK2", "file日志")

							if o, ok := inf.(zapcore.Core); ok {
								fmt.Println("OK2", o)
								otherLogger = NewLoggerWithCore(o)
							}
						}

						// 不可以
						// zapcore.Core(inf)
						// zapcore.Core(actualValue)
						if actualValue.Type().String() == "*zapcore.ioCore" {
							fmt.Println(actualValue.Kind(), actualValue.Type())
						}
						fmt.Println(actualValue.Kind(), actualValue.Type())
					}
				}
			}
		}
	}
	return consoleLogger, otherLogger
}
