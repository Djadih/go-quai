package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

type Logger struct {
	*logrus.Logger
}

var Log Logger = Logger{logrus.New()}

func ConfigureLogger(ctx *cli.Context) {
	logLevel := logrus.Level(ctx.GlobalInt("verbosity"))
	Log.SetLevel(logLevel)

	if logLevel >= logrus.DebugLevel {
		Log.SetReportCaller(true)
	}

	log_filename := "nodelogs"
	regionNum := ctx.GlobalString("region")

	if ctx.GlobalIsSet("zone") {
		zoneNum := ctx.GlobalString("zone")
		log_filename = filepath.Join(log_filename, "zone-"+regionNum+"-"+zoneNum)
	} else if ctx.GlobalIsSet("region") {
		log_filename = filepath.Join(log_filename, "region-"+regionNum)
	} else {
		log_filename = filepath.Join(log_filename, "prime")
	}
	log_filename += ".log"

	Log.Formatter = &logrus.TextFormatter{
		ForceColors:      ctx.GlobalBool("showcolors"),
		PadLevelText:     true,
		FullTimestamp:    true,
		TimestampFormat:  "01-02|15:04:05",
		CallerPrettyfier: callerPrettyfier,
	}

	Log.SetOutput(&lumberjack.Logger{
		Filename:   log_filename,
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})
}

func SetLevelInt(level int) {
	Log.SetLevel(logrus.Level(level))
}

func SetLevelString(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		Log.Error("Invalid log level: ", level)
		return
	}
	Log.SetLevel(logLevel)
}

func New(out_path string) Logger {
	logger := logrus.New()

	Log.Formatter = &logrus.TextFormatter{
		ForceColors:      Log.Logger.Formatter.(*logrus.TextFormatter).ForceColors,
		PadLevelText:     true,
		FullTimestamp:    true,
		TimestampFormat:  "01-02|15:04:05",
		CallerPrettyfier: callerPrettyfier,
	}

	logger.SetOutput(&lumberjack.Logger{
		Filename:   out_path,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})
	return Logger{logger}
}

// Uses of the global logger will use the following static method.
func Trace(msg string, args ...interface{}) {
	Log.Trace(constructLogMessage(msg, args...))
}

// Individual logging instances will use the following method.
func (l Logger) Trace(msg string, args ...interface{}) {
	l.Logger.Trace(constructLogMessage(msg, args...))
}

func Debug(msg string, args ...interface{}) {
	Log.Debug(constructLogMessage(msg, args...))
}
func (l Logger) Debug(msg string, args ...interface{}) {
	l.Logger.Debug(constructLogMessage(msg, args...))
}

func Info(msg string, args ...interface{}) {
	Log.Info(constructLogMessage(msg, args...))
}
func (l Logger) Info(msg string, args ...interface{}) {
	l.Logger.Info(constructLogMessage(msg, args...))
}

func Warn(msg string, args ...interface{}) {
	Log.Warn(constructLogMessage(msg, args...))
}
func (l Logger) Warn(msg string, args ...interface{}) {
	l.Logger.Warn(constructLogMessage(msg, args...))
}

func Error(msg string, args ...interface{}) {
	Log.Error(constructLogMessage(msg, args...))
}
func (l Logger) Error(msg string, args ...interface{}) {
	l.Logger.Error(constructLogMessage(msg, args...))
}

func Fatal(msg string, args ...interface{}) {
	Log.Fatal(constructLogMessage(msg, args...))
}
func (l Logger) Fatal(msg string, args ...interface{}) {
	l.Logger.Fatal(constructLogMessage(msg, args...))
}

func Panic(msg string, args ...interface{}) {
	Log.Panic(constructLogMessage(msg, args...))
}
func (l Logger) Panic(msg string, args ...interface{}) {
	l.Logger.Panic(constructLogMessage(msg, args...))
}

func Lazy(fn func() string, logLevel string) {
	level, err := logrus.ParseLevel(logLevel)
	if err == nil && Log.IsLevelEnabled(level) {
		callCorrectLevel(level, fn())
	}
}

func callCorrectLevel(level logrus.Level, msg string, args ...interface{}) {
	switch level {
	case logrus.TraceLevel:
		Trace(msg, args...)
	case logrus.DebugLevel:
		Debug(msg, args...)
	case logrus.InfoLevel:
		Info(msg, args...)
	case logrus.WarnLevel:
		Warn(msg, args...)
	case logrus.ErrorLevel:
		Error(msg, args...)
	case logrus.FatalLevel:
		Fatal(msg, args...)
	case logrus.PanicLevel:
		Panic(msg, args...)
	default:
		Error("Unknown log level: %v", level)
	}
}

func constructLogMessage(msg string, fields ...interface{}) string {
	var pairs []string

	if len(fields)%2 != 0 {
		fields = append(fields, "MISSING VALUE")
	}

	for i := 0; i < len(fields); i += 2 {
		key := fields[i]
		value := fields[i+1]
		pairs = append(pairs, fmt.Sprintf("%v=%v", key, value))
	}

	return fmt.Sprintf("%-40s %s", msg, strings.Join(pairs, " "))
}

func callerPrettyfier(f *runtime.Frame) (string, string) {
	// pc := make([]uintptr, 1)
	// // runtime.Callers(2, f.PC)
	// pc[0] = f.PC
	// fmt.Printf("f pc: %v\n", f.PC)
	// callerFrames := runtime.CallersFrames(pc)
	// tempFrame, _ := callerFrames.Next()
	// fmt.Printf("tempFrame: %v\n", tempFrame)
	// caller, _ := callerFrames.Next()
	// fmt.Printf("caller: %v\n", caller)
	// f = &caller

	// filename := path.Base(f.File)
	// dir := path.Base(path.Dir(f.File))

	// filepath := path.Join(dir, filename)
	// return "", fmt.Sprintf("%s:%d", filepath, f.Line)

	fmt.Printf("f: %v\n", f)

	pc := make([]uintptr, 1)
	pc[0] = f.PC

	func_info := runtime.FuncForPC(f.PC)
	fmt.Printf("func_info: %v\n", func_info)
	return "", ""

	// working ish but inconsistent
	// pc := make([]uintptr, 1)
	// runtime.Callers(10, pc)
	// callerFrames := runtime.CallersFrames(pc)
	// caller, _ := callerFrames.Next()
	// return "", fmt.Sprintf("%s:%d", caller.File, caller.Line)
}

// func callerInfo() string {
// 	pc := make([]uintptr, 1)
// 	runtime.Callers(3, pc)
// 	callerFrames := runtime.CallersFrames(pc)
// 	caller, _ := callerFrames.Next()
// 	return fmt.Sprintf("%s:%d", caller.File, caller.Line)
// }
