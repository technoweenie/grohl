package grohl

type Context struct {
	data              Data
	Logger            Logger
	TimeUnit          string
	ExceptionReporter ExceptionReporter
	*_statter
}

func (c *Context) Log(data Data) error {
	return c.Logger.Log(c.Merge(data))
}

func (c *Context) log(data Data) error {
	return c.Logger.Log(data)
}

func (c *Context) New(data Data) *Context {
	return newContext(c.Merge(data), c.Logger, c.TimeUnit, c.ExceptionReporter)
}

func (c *Context) Add(key string, value interface{}) {
	c.data[key] = value
}

func (c *Context) Merge(data Data) Data {
	if data == nil {
		return dupeMaps(c.data)
	} else {
		return dupeMaps(c.data, data)
	}
}

func (c *Context) Delete(key string) {
	delete(c.data, key)
}

func dupeMaps(maps ...Data) Data {
	merged := make(Data)
	for _, orig := range maps {
		for key, value := range orig {
			merged[key] = value
		}
	}
	return merged
}

func newContext(data Data, logger Logger, timeunit string, reporter ExceptionReporter) *Context {
	return &Context{
		data:              data,
		Logger:            logger,
		TimeUnit:          timeunit,
		ExceptionReporter: reporter,
		_statter:          &_statter{},
	}
}
