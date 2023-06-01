package logger_utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

// type Logger struct {
// 	*logrus.Logger
// }

// var Log Logger = Logger{logrus.New()}
// var logger *logrus.Logger
// type Logger *logrus.Logger

func ConfigureDefaultLogger(ctx *cli.Context) {
	logLevel := logrus.Level(ctx.GlobalInt("verbosity"))
	logrus.SetLevel(logLevel)

	if logLevel >= logrus.DebugLevel {
		logrus.SetReportCaller(true)
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

	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	ForceColors:      ctx.GlobalBool("showcolors"),
	// 	PadLevelText:     true,
	// 	FullTimestamp:    true,
	// 	ForceQuote: 	 true,
	// 	QuoteEmptyFields: true,
	// 	TimestampFormat:  "01-02|15:04:05",
	// 	CallerPrettyfier: callerPrettyfier,
	// })

	logrus.SetFormatter(new(customFormatter))

	logrus.SetOutput(&lumberjack.Logger{
		Filename:   log_filename,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})
}

func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

func SetLevelInt(level int) {
	logrus.SetLevel(logrus.Level(level))
}

func SetLevelString(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Error("Invalid log level: ", level)
		return
	}
	logrus.SetLevel(logLevel)
}

func New(out_path string) *logrus.Logger {
	logger := logrus.New()

	logger.Formatter = &logrus.TextFormatter{
		ForceColors:      logger.Formatter.(*logrus.TextFormatter).ForceColors,
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
	return logger
}

func Lazy(fn func() string, logLevel string) {
	level, err := logrus.ParseLevel(logLevel)
	if err == nil && logrus.IsLevelEnabled(level) {
		callCorrectLevel(level, fn())
	}
}

func callCorrectLevel(level logrus.Level, msg string) {
	switch level {
	case logrus.TraceLevel:
		logrus.Trace(msg)
	case logrus.DebugLevel:
		logrus.Debug(msg)
	case logrus.InfoLevel:
		logrus.Info(msg)
	case logrus.WarnLevel:
		logrus.Warn(msg)
	case logrus.ErrorLevel:
		logrus.Error(msg)
	case logrus.FatalLevel:
		logrus.Fatal(msg)
	case logrus.PanicLevel:
		logrus.Panic(msg)
	default:
		logrus.Error("Unknown log level: %v", level)
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
	filename := filepath.Base(f.File)
	dir := filepath.Base(filepath.Dir(f.File))

	filepath := filepath.Join(dir, filename)
	return "", fmt.Sprintf("%s:%d", filepath, f.Line)
}

type customFormatter struct {

}

func (f *customFormatter) Format (entry *logrus.Entry) ([]byte, error) {
	// fmt.Printf("entry: %v\n", entry)
	fmt.Print(entry.Message + "\n")
	return nil, nil
}

// func callerPrettyfier(f *runtime.Frame) (string, string) {
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

// fmt.Printf("f: %v\n", f)

// pc := make([]uintptr, 1)
// pc[0] = f.PC

// func_info := runtime.FuncForPC(f.PC)
// fmt.Printf("func_info: %v\n", func_info)
// return "", ""

// working ish but inconsistent
// pc := make([]uintptr, 1)
// runtime.Callers(10, pc)
// pc[0] = f.PC
// fmt.Printf("test: %v\n", runtime.FuncForPC(pc[0]).Name())
// callerFrames := runtime.CallersFrames(pc)

// caller, _ := callerFrames.Next()
// return "", fmt.Sprintf("%s:%d", caller.File, caller.Line)

// return "", fmt.Sprintf("%s:%d", f.File, f.Line)
// }

// func callerInfo() string {
// 	pc := make([]uintptr, 1)
// 	runtime.Callers(3, pc)
// 	callerFrames := runtime.CallersFrames(pc)
// 	caller, _ := callerFrames.Next()
// 	return fmt.Sprintf("%s:%d", caller.File, caller.Line)
// }
