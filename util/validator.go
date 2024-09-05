package util

import "regexp"

func ValidateEmail(email string) bool {
	// Regular expression pattern for email validation
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regular expression
	re := regexp.MustCompile(emailRegexPattern)

	// Return whether the email matches the regex pattern
	return re.MatchString(email)
}
