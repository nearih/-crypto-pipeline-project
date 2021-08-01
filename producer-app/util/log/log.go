package log

import (
	"fmt"
	"os"
	"producer-app/config"
	"runtime"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/sirupsen/logrus"
)

type ILogType interface {
	ToJson()
}

type LogType string

// ToJson This method is created to prevent from assign string to type LogType
func (l LogType) ToJson() {
	fmt.Println("This method is created to prevent from assign string to type LogType")
}

const (
	TransactionLog LogType = "TransactionLog" // log what need to be document
	EventLog       LogType = "EventLog"       // log that need specific trigger point/ notification to some party
	AuditLog       LogType = "AuditLog"       // log activity of admin in admin service
	ActivityLog    LogType = "ActivityLog"    // log customer activity
	AppLog         LogType = "AppLog"         // any other log
)

type Logger struct {
	*logrus.Logger
}

type thaiTimeFormatter struct {
	logrus.Formatter
}

func (t thaiTimeFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.In(logTimezoneLocation)
	return t.Formatter.Format(e)
}

var logTimezoneLocation *time.Location

// NewLogger logger instance
func NewLogger(c *config.RootConfig) *Logger {
	log := logrus.New()

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	logTimezoneLocation = loc
	logFormat := c.Log.Format
	logLevel := c.Log.Level

	switch strings.ToLower(logFormat) {
	case "json":
		log.SetFormatter(thaiTimeFormatter{&logrus.JSONFormatter{
			PrettyPrint: true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg: "message",
			},
		}})
	default:
		log.SetFormatter(thaiTimeFormatter{&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				return "", ""
			},
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg: "message",
			}}})
	}

	log.SetOutput(os.Stdout)

	switch strings.ToLower(logLevel) {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
		log.SetReportCaller(true)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	}

	return &Logger{log}
}

// func (l *Logger) ErrorWithID(ctx context.Context, logType ILogType, args ...interface{}) {
// 	ctxResponse, _ := context_wrapper.ExtractGrpcContext(ctx)

// 	var idMessage string
// 	idMessage = fmt.Sprintf("user:%v ", ctxResponse.UserId)
// 	if ctxResponse.UserId == "" {
// 		idMessage = fmt.Sprintf("deviceId:%v ", ctxResponse.DeviceID)
// 	}

// 	if l.Level >= logrus.InfoLevel {
// 		entry := l.WithFields(logrus.Fields{})
// 		entry.Data["txId"] = ctxResponse.TxId
// 		entry.Data["logType"] = logType
// 		entry.Data["fields.file"] = fileInfo(2)
// 		entry.Data["identity"] = idMessage
// 		entry.Data["timestamp"] = time.Now().In(logTimezoneLocation).Format(time.RFC3339)
// 		entry.Error(args...)
// 	}
// }

// func (l *Logger) DebugWithID(ctx context.Context, logType ILogType, args ...interface{}) {
// 	ctxResponse, _ := context_wrapper.ExtractGrpcContext(ctx)

// 	var idMessage string
// 	idMessage = fmt.Sprintf("user:%v ", ctxResponse.UserId)
// 	if ctxResponse.UserId == "" {
// 		idMessage = fmt.Sprintf("deviceId:%v ", ctxResponse.DeviceID)
// 	}

// 	if l.Level >= logrus.InfoLevel {
// 		entry := l.WithFields(logrus.Fields{})
// 		entry.Data["txId"] = ctxResponse.TxId
// 		entry.Data["logType"] = logType
// 		entry.Data["fields.file"] = fileInfo(2)
// 		entry.Data["identity"] = idMessage
// 		entry.Data["timestamp"] = time.Now().In(logTimezoneLocation).Format(time.RFC3339)
// 		entry.Debug(args...)
// 	}
// }

// func (l *Logger) InfoWithID(ctx context.Context, logType ILogType, args ...interface{}) {
// 	ctxResponse, _ := context_wrapper.ExtractGrpcContext(ctx)

// 	var idMessage string
// 	idMessage = fmt.Sprintf("user:%v ", ctxResponse.UserId)
// 	if ctxResponse.UserId == "" {
// 		idMessage = fmt.Sprintf("deviceId:%v ", ctxResponse.DeviceID)
// 	}

// 	if l.Level >= logrus.InfoLevel {
// 		entry := l.WithFields(logrus.Fields{})
// 		entry.Data["txId"] = ctxResponse.TxId
// 		entry.Data["logType"] = logType
// 		entry.Data["fields.file"] = fileInfo(2)
// 		entry.Data["identity"] = idMessage
// 		entry.Data["timestamp"] = time.Now().In(logTimezoneLocation).Format(time.RFC3339)
// 		entry.Info(args...)
// 	}
// }

// func fileInfo(skip int) string {
// 	_, file, line, ok := runtime.Caller(skip)
// 	if !ok {
// 		file = "Cannot get file"
// 		line = 1
// 	}
// 	return fmt.Sprintf("%s:%d", file, line)
// }
