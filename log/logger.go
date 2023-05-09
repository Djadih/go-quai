package log

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dominant-strategies/go-quai/common"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

type Logger = *logrus.Logger

var Log Logger = logrus.New()

type logger struct {
	ctx []interface{}
}

func init() {
	Log.SetLevel(logrus.DebugLevel)
}

func ConfigureLogger(ctx *cli.Context) {

	log_filename := "nodelogs/"

	if common.NodeLocation.Context() == common.PRIME_CTX {
		log_filename += "prime"
	} else {
		regionNum := strconv.Itoa(common.NodeLocation.Region())

		if common.NodeLocation.Context() == common.REGION_CTX {
			log_filename += "region-" + regionNum
		} else if common.NodeLocation.Context() == common.ZONE_CTX {
			zoneNum := strconv.Itoa(common.NodeLocation.Zone())
			log_filename += "zone-" + regionNum + "-" + zoneNum
		} else {
			panic("Invalid node location context")
		}
	}
	log_filename += ".logNew"

	fmt.Println("Log file name: ", log_filename)

	Log.SetOutput(&lumberjack.Logger{
		Filename:   log_filename,
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})
}

func New(out_path string) Logger {
	logger := logrus.New()
	logger.SetOutput(&lumberjack.Logger{
		Filename:   out_path,
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})
	return logger
}

func Trace(msg string, args ...interface{}) {
	Log.Trace(constructLogMessage(msg, args...))
}

func Debug(msg string, args ...interface{}) {
	Log.Debug(constructLogMessage(msg, args...))
}

func Info(msg string, args ...interface{}) {
	Log.Info(constructLogMessage(msg, args...))
}

func Warn(msg string, args ...interface{}) {
	Log.Warn(constructLogMessage(msg, args...))
}

func Error(msg string, args ...interface{}) {
	Log.Error(constructLogMessage(msg, args...))
}

func Fatal(msg string, args ...interface{}) {
	Log.Fatal(constructLogMessage(msg, args...))
}

func Panic(msg string, args ...interface{}) {
	Log.Panic(constructLogMessage(msg, args...))
}

func constructLogMessage(msg string, args ...interface{}) string {
	sb := strings.Builder{}
	sb.WriteString(msg)

	for _, arg := range args {
		sb.WriteString(fmt.Sprintf(" %v", arg))
	}

	return sb.String()
}
