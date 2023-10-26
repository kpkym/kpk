package common

import (
	"github.com/gookit/color"
	"os"
)

var pinkC = color.HEX("FF69B4")

var level = 2

var tracingC = color.FgGray
var debugC = color.FgCyan
var infoC = color.FgGreen
var warningC = color.FgYellow
var errorC = color.FgRed

func canPrint(nowLevel int) bool {
	return level <= nowLevel
}

func PinkPrintLn(msg ...any) {
	pinkC.Println(msg...)
}

func TracingPrintLn(msg ...any) {
	if canPrint(0) {
		tracingC.Println(msg...)
	}
}

func DebugPrintLn(msg ...any) {
	if canPrint(1) {
		debugC.Println(msg...)
	}
}

func InfoPrintLn(msg ...any) {
	infoC.Println(msg...)
}

func WarningPrintLn(msg ...any) {
	warningC.Println(msg...)
}

func ErrorPrintLn(msg ...any) {
	errorC.Println(msg...)
	os.Exit(1)
}
