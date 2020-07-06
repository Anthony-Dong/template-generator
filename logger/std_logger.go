package logger

import (
	"fmt"
	"os"
)

const (
	generator = "GEN"
)

// \033[33m[WARN]\033[0m
// \033[36m[DEBUG]\033[0m

func InfoF(format string, v ...interface{}) {
	str := fmt.Sprintf("\033[32m[%s-INFO]\033[0m %s\n", generator, format)
	fmt.Printf(str, v ...)
}

func ErrorF(format string, v ...interface{}) {
	str := fmt.Sprintf("\033[31m[%s-ERROR]\033[0m %s\n", generator, format)
	fmt.Printf(str, v ...)
}

func FatalF(format string, v ...interface{}) {
	str := fmt.Sprintf("\033[35m[%s-FATAL]\033[0m %s\n", generator, format)
	fmt.Printf(str, v ...)
	os.Exit(-1)
}
