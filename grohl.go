package grohl

var defaultLogger = NewLogger(nil)

func Log(data map[string]interface{}) {
	defaultLogger.Log(data)
}

func NewContext(data map[string]interface{}) *IoLogger {
	return defaultLogger.NewContext(data)
}

func NewExceptionLogger(reporter ExceptionReporter) *ExceptionLogger {
	return newExceptionLogger(defaultLogger, reporter)
}

func AddContext(key string, value interface{}) {
	defaultLogger.AddContext(key, value)
}

func DeleteContext(key string) {
	defaultLogger.DeleteContext(key)
}

func NewTimer(context map[string]interface{}) *Timer {
	return defaultLogger.NewTimer(context)
}

func SetLogger(l *IoLogger) {
	defaultLogger = l
}

func SetTimeUnit(unit string) {
	defaultLogger.TimeUnit = unit
}

func TimeUnit() string {
	return defaultLogger.TimeUnit
}
