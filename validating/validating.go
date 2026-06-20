package validating

import (
	"strings"
)



func IsFromChars(str string, chars []string) (bool, string) {
	for _, ch := range str {
		found := false
		char := string(ch)
		for _, c := range chars {
			if char == c {
				found = true
				break
			}
		}
		if !found {
			return false, char
		}
	}
	return true, ""
}



func HasNotAllowedSpace(str string) bool {
	return strings.HasPrefix(str, " ") || strings.HasSuffix(str, " ") || strings.Contains(str, "  ")
}



