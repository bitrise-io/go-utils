package log

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/colorstring"
)

// Fail ...
func Fail(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Println(colorstring.Red(message))
	os.Exit(1)
}

// Error ...
func Error(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Printf(colorstring.Red(message))
}

// Warn ...
func Warn(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Printf(colorstring.Yellow(message))
}

// Info ...
func Info(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Printf(colorstring.Blue(message))
}

// Detail ...
func Detail(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Println(message)
}

// Done ...
func Done(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Printf(colorstring.Green(message))
}
