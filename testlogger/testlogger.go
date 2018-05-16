package testlogger

import "github.com/bitrise-io/go-utils/log"

type logger interface {
	Donef(string, ...interface{})
	Printf(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Timestamp() log.DefaultLogger
}

// MyPkg ...
type MyPkg struct {
	l logger
}

var lgr logger

// SetLogger ...
func SetLogger(lg logger) {
	lgr = lg
}

// TestFn1 ...
func (m MyPkg) TestFn1() {
	lgr.Timestamp().Donef("function called")
	lgr.Printf("some data")
	lgr.Errorf("+ some data")
	lgr.Timestamp().Warnf("+ some data")
}

// TestFn2 ...
func (m MyPkg) TestFn2() {
	lgr.Printf("some data")
	lgr.Timestamp().Donef("function called")
	lgr.Errorf("+ some data")
	lgr.Timestamp().Warnf("+ some data")
}
