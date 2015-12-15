package jshapi

// Logger describes logger interface matching stdlib logger
// copied from https://godoc.org/github.com/Sirupsen/logrus#StdLogger
type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}
