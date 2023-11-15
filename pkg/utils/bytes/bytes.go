package bytes

// HasQuotes checks if the first and last characters are either " or '.
func HasQuotes(input []byte) bool {
	return len(input) >= 2 &&
		((input[0] == '"' && input[len(input)-1] == '"') ||
			(input[0] == '\'' && input[len(input)-1] == '\''))
}

// TrimQuotes removes the first and last character if they are both either
// " or ', otherwise it is a noop.
func TrimQuotes(input []byte) []byte {
	if HasQuotes(input) {
		return input[1 : len(input)-1]
	}
	return input
}
