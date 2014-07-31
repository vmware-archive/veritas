package say

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var EnableColor = true

const DefaultStyle = "\x1b[0m"
const BoldStyle = "\x1b[1m"
const RedColor = "\x1b[91m"
const GreenColor = "\x1b[32m"
const YellowColor = "\x1b[33m"
const CyanColor = "\x1b[36m"
const GrayColor = "\x1b[90m"
const LightGrayColor = "\x1b[37m"

func Red(format string, args ...interface{}) string {
	return Colorize(RedColor, format, args...)
}

func Green(format string, args ...interface{}) string {
	return Colorize(GreenColor, format, args...)
}

func Yellow(format string, args ...interface{}) string {
	return Colorize(YellowColor, format, args...)
}

func Cyan(format string, args ...interface{}) string {
	return Colorize(CyanColor, format, args...)
}

func Gray(format string, args ...interface{}) string {
	return Colorize(GrayColor, format, args...)
}

func LightGray(format string, args ...interface{}) string {
	return Colorize(LightGrayColor, format, args...)
}

func Colorize(colorCode string, format string, args ...interface{}) string {
	var out string

	if len(args) > 0 {
		out = fmt.Sprintf(format, args...)
	} else {
		out = format
	}

	if EnableColor {
		return fmt.Sprintf("%s%s%s", colorCode, out, DefaultStyle)
	} else {
		return out
	}
}

func Ask(text string) string {
	Print(0, text+":\n> ")
	var response string
	fmt.Scanln(&response)
	return response
}

func AskWithDefault(text string, defaultResponse string) string {
	Print(0, "%s [%s]:\n> ", text, Green(defaultResponse))
	var response string
	fmt.Scanln(&response)
	if response == "" {
		return defaultResponse
	}
	return response
}

func AskForIntegerWithDefault(text string, defaultResponse int) int {
	Print(0, "%s [%s]:\n> ", text, Green("%d", defaultResponse))
	var response string
	fmt.Scanln(&response)
	if response == "" {
		return defaultResponse
	}
	asInteger, err := strconv.Atoi(response)
	if err != nil {
		Println(0, Red("That was an invalid response..."))
		return AskForIntegerWithDefault(text, defaultResponse)
	}

	return asInteger
}

func AskForBoolWithDefault(text string, defaultResponse bool) bool {
	Print(0, "%s [%s]:\n> ", text, Green("%t", defaultResponse))
	var response string
	fmt.Scanln(&response)
	if response == "true" {
		return true
	}
	if response == "false" {
		return false
	}
	if response == "" {
		return defaultResponse
	}
	Println(0, Red("That was an invalid response... try 'true' or 'false'"))
	return AskForBoolWithDefault(text, defaultResponse)
}

func Pick(text string, options []string) string {
	Println(0, "%s:", text)
	for i, option := range options {
		Println(1, "[%s] %s", Green("%d", i), option)
	}
	Print(0, "> ")
	var response string
	fmt.Scanln(&response)
	index, err := strconv.Atoi(response)
	if err != nil {
		Println(0, Red("That was an invalid selection..."))
		return Pick(text, options)
	}
	return options[index]
}

func PrintBanner(text string, bannerCharacter string) {
	FprintBanner(os.Stdout, text, bannerCharacter)
}

func PrintDelimiter() {
	FprintDelimiter(os.Stdout)
}

func Print(indentation int, format string, args ...interface{}) {
	Fprint(os.Stdout, indentation, format, args...)
}

func Println(indentation int, format string, args ...interface{}) {
	Fprintln(os.Stdout, indentation, format, args...)
}

func FprintBanner(w io.Writer, text string, bannerCharacter string) {
	fmt.Fprintln(w, text)
	fmt.Fprintln(w, strings.Repeat(bannerCharacter, len(text)))
}

func FprintDelimiter(w io.Writer) {
	fmt.Fprintln(w, Colorize(GrayColor, "%s", strings.Repeat("-", 30)))
}

func Fprint(w io.Writer, indentation int, format string, args ...interface{}) {
	fmt.Fprint(w, Indent(indentation, format, args...))
}

func Fprintln(w io.Writer, indentation int, format string, args ...interface{}) {
	fmt.Fprintln(w, Indent(indentation, format, args...))
}

func Clear() {
	fmt.Print("\x1b[2J\x1b[;H")
}

func Indent(indentation int, format string, args ...interface{}) string {
	var text string

	if len(args) > 0 {
		text = fmt.Sprintf(format, args...)
	} else {
		text = format
	}

	stringArray := strings.Split(text, "\n")
	padding := ""
	if indentation >= 0 {
		padding = strings.Repeat("  ", indentation)
	}
	for i, s := range stringArray {
		stringArray[i] = fmt.Sprintf("%s%s", padding, s)
	}

	return strings.Join(stringArray, "\n")
}
