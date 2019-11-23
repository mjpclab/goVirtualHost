package util

func FilterEmptyStrings(inputs []string) []string {
	output := make([]string, 0, len(inputs))

	for _, s := range inputs {
		if len(s) > 0 {
			output = append(output, s)
		}
	}

	return output
}
