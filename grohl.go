package grohl

var DefaultLogger = NewLogger(nil)

func Log(data map[string]interface{}) {
	DefaultLogger.Log(data)
}

func NewContext(data map[string]interface{}) *IoLogger {
	return DefaultLogger.NewContext(data)
}

func NewExceptionLogger(reporter ExceptionReporter) *ExceptionLogger {
	return newExceptionLogger(DefaultLogger, reporter)
}

func AddContext(key string, value interface{}) {
	DefaultLogger.AddContext(key, value)
}

func DeleteContext(key string) {
	DefaultLogger.DeleteContext(key)
}

func NewTimer(context map[string]interface{}) *Timer {
	return DefaultLogger.NewTimer(context)
}

func SetTimeUnit(unit string) {
	DefaultLogger.TimeUnit = unit
}

func TimeUnit() string {
	return DefaultLogger.TimeUnit
}
