package slf4go

var loggerFactory LoggerFactory = newNativeLoggerFactory()

// Backend set new slf4go backend logger factory
func Backend(factory LoggerFactory) {
	if factory == nil {
		panic("factory can't be nil")
	}

	loggerFactory = factory
}

// Get get/create new logger by name
func Get(name string) Logger {
	return loggerFactory.GetLogger(name)
}
