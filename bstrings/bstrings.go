package bstrings

func getEnLittters() []string {
	letters := []string{}

	for ch := 'a'; ch <= 'z'; ch++ {
    letters = append(letters, string(ch))
	}
	for ch := 'A'; ch <= 'Z'; ch++ {
    letters = append(letters, string(ch))
	}
	return letters
}

func getNumbers() []string {
	numbers := []string{}

	for ch := '0'; ch <= '9'; ch++ {
		numbers = append(numbers, string(ch))
	}
	return numbers
}

var EnLittters = getEnLittters()

var AllowedUsernameMarks = []string{"_", "-"}

var AllowedPasswordMarks = []string{
	"_", "-", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")",
	"+", "=", "{", "}", "[", "]", "|", "\\", ":", ";", "'", "\"",
	"<", ">", ",", ".", "?", "/"}

var Numbers = getNumbers()