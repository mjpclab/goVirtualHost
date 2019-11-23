package util

import "testing"

func TestFilterEmptyStrings(t *testing.T) {
	inputs := []string{"aa", "", "bb"}
	results := FilterEmptyStrings(inputs)
	if len(results) != 2 || results[0] != "aa" || results[1] != "bb" {
		t.Error(results)
	}
}
