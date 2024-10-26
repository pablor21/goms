package console

var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"
var colorBlue = "\033[34m"
var colorPurple = "\033[35m"
var colorCyan = "\033[36m"
var colorWhite = "\033[37m"
var colorReset = "\033[0m"

func Red(str string) string {
	return colorRed + str + colorReset
}

func Green(str string) string {
	return colorGreen + str + colorReset
}

func Yellow(str string) string {
	return colorYellow + str + colorReset
}

func Blue(str string) string {
	return colorBlue + str + colorReset
}

func Purple(str string) string {
	return colorPurple + str + colorReset
}

func Cyan(str string) string {
	return colorCyan + str + colorReset
}

func White(str string) string {
	return colorWhite + str + colorReset
}

func Bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func Underline(str string) string {
	return "\033[4m" + str + "\033[0m"
}

func Inverse(str string) string {
	return "\033[7m" + str + "\033[0m"
}

func Strikethrough(str string) string {
	return "\033[9m" + str + "\033[0m"
}
