package sqldatabase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"linebot-go/common/infrastructure/consts/contextKey"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"

	"gorm.io/gorm/utils"
	commonUtils "linebot-go/common/application/utils"
)

// ErrRecordNotFound record not found error
var ErrRecordNotFound = errors.New("record not found")

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// Writer log writer interface
type Writer interface {
	Printf(string, ...interface{})
}

var (
	// Discard logger will print any log to io.Discard
	Discard = NewCustomLogger(log.New(io.Discard, "", log.LstdFlags), logger.Config{})
	// Default Default logger
	Default = NewCustomLogger(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
)

// New initialize logger
func NewCustomLogger(writer Writer, config logger.Config) logger.Interface {
	var (
		infoStr      = "actionId:%s %s\n[info] "
		warnStr      = "actionId:%s %s\n[warn] "
		errStr       = "actionId:%s %s\n[error] "
		traceStr     = "actionId:%s %s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "actionId:%s %s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "actionId:%s %s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = "actionId:%s " + Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = "actionId:%s " + BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = "actionId:%s " + Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = "actionId:%s " + Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = "actionId:%s " + Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = "actionId:%s " + RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &CustomLogger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type CustomLogger struct {
	Writer
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l *CustomLogger) Info(c context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l *CustomLogger) Warn(c context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l *CustomLogger) Error(c context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
//
//nolint:cyclop
func (l *CustomLogger) Trace(c context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	actionId := commonUtils.GetActionId(c)
	trace := "false"
	if c.Value(contextKey.SqlTrace) != nil {
		trace = "true"
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr, actionId, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceErrStr, actionId, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr, actionId, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceWarnStr, actionId, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case trace == "true" || l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr, actionId, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceStr, actionId, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// ParamsFilter filter params
func (l *CustomLogger) ParamsFilter(c context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
