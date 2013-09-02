package grohl

var logger = NewLogger(nil)

func Log(data map[string]interface{}) {
  logger.Log(data)
}

func NewContext(data map[string]interface{}) *IoLogger {
  return logger.NewContext(data)
}

func AddContext(key string, value interface{}) {
  logger.AddContext(key, value)
}

func DeleteContext(key string) {
  logger.DeleteContext(key)
}

func NewTimer(context map[string]interface{}) *timer {
  return logger.NewTimer(context)
}

func SetLogger(l *IoLogger) {
  logger = l
}

func SetTimeUnit(unit string) {
  logger.TimeUnit = unit
}

func TimeUnit() string {
  return logger.TimeUnit
}
