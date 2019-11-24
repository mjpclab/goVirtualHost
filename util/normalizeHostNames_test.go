package util

import "testing"

func TestNormalizeHostNames(t *testing.T) {
	inputs := []string{"aA", "", "Bb"}
	results := NormalizeHostNames(inputs)
	if len(results) != 2 || results[0] != "aa" || results[1] != "bb" {
		t.Error(results)
	}
}
