package text

import "fmt"

func PurpleText(text string, args ...interface{}) string {
	return fmt.Sprintf("\033[35m%s\033[0m", fmt.Sprintf(text, args...))
}

func DarkYellowText(text string, args ...interface{}) string {
	return fmt.Sprintf("\033[33m\033[2m%s\033[0m", fmt.Sprintf(text, args...))
}
