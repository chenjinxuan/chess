package log

import (
	"github.com/cihub/seelog"
)

var (
	logger    seelog.LoggerInterface
	loggerTpl = `
	<seelog minlevel="debug">
		<outputs>
			<filter formatid="standard" levels="debug,info,warn,error,critical">
				<console />
			</filter>
		</outputs>
		<formats>
			<format id="standard" format="[%Date %Time] [%Level] [%Msg]%n" />
		</formats>
	</seelog>`
)

func init() {
	logger, _ = seelog.LoggerFromConfigAsBytes([]byte(loggerTpl))
}

// Tracef formats message according to format specifier
// and writes to log with level = Trace.
func Tracef(format string, params ...interface{}) {
	logger.Tracef(format, params...)
}

// Debugf formats message according to format specifier
// and writes to log with level = Debug.
func Debugf(format string, params ...interface{}) {
	logger.Debugf(format, params...)
}

// Infof formats message according to format specifier
// and writes to log with level = Info.
func Infof(format string, params ...interface{}) {
	logger.Infof(format, params...)
}

// Warnf formats message according to format specifier
// and writes to log with level = Warn.
func Warnf(format string, params ...interface{}) error {
	return logger.Warnf(format, params...)
}

// Errorf formats message according to format specifier
// and writes to log with level = Error.
func Errorf(format string, params ...interface{}) error {
	return logger.Errorf(format, params...)
}

// Criticalf formats message according to format specifier
// and writes to log with level = Critical.
func Criticalf(format string, params ...interface{}) error {
	return logger.Criticalf(format, params...)
}

// Trace formats message using the default formats for its operands
// and writes to log with level = Trace
func Trace(v ...interface{}) {
	logger.Trace(v...)
}

// Debug formats message using the default formats for its operands
// and writes to log with level = Debug
func Debug(v ...interface{}) {
	logger.Debug(v...)
}

// Info formats message using the default formats for its operands
// and writes to log with level = Info
func Info(v ...interface{}) {
	logger.Info(v...)
}

// Warn formats message using the default formats for its operands
// and writes to log with level = Warn
func Warn(v ...interface{}) error {
	return logger.Warn(v...)
}

// Error formats message using the default formats for its operands
// and writes to log with level = Error
func Error(v ...interface{}) error {
	return logger.Error(v...)
}

// Critical formats message using the default formats for its operands
// and writes to log with level = Critical
func Critical(v ...interface{}) error {
	return logger.Critical(v...)
}

func PrintErr(err error, funcName string) error {
	return Errorf(funcName+" Error(%s)", err)
}
