package utils

import (
	"errors"
	"strings"
	"unicode"
)


func ConvertToMobile(mobile string) (string, error) {

	mobile = strings.Trim(mobile, "")
	// remove non numeric
	mobile = RemoveNonNumbers(mobile)

	// Validate the mobile number format (basic example, expand as needed)
	if  len(mobile) < 10 || mobile[0] != '0' {
		return "", errors.New("Invalid mobile number")
	}

	return mobile, nil
}


func RemoveNonNumbers(input string) string {
	// Create a new string builder
	var builder strings.Builder

	// Iterate over the input string
	for _, r := range input {
		// If the current character is not a digit, append it to the string builder
		if unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}

	// Get the resulting string from the string builder
	output := builder.String()
	return output
}
