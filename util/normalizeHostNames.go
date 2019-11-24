package util

import "strings"

func NormalizeHostNames(inputs []string) []string {
	output := make([]string, 0, len(inputs))

	for _, str := range inputs {
		if len(str) > 0 {
			name := strings.ToLower(str)
			output = append(output, name)
		}
	}

	return output
}
