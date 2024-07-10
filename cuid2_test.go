package cuid2

import (
	"strings"
	"testing"
)

func create(t *testing.T) string {
	id := CreateId()
	if id == "" {
		t.Fatalf("Cuid not generated")
	}
	return id
}

func charOfString(input string, index int) string {
	return string(input[index])
}

func TestGeneratesCuid2(t *testing.T) {
	create(t)
}

func TestFirstCharFromAtoZ(t *testing.T) {
	atoz := "abcdefghijklmnopqrstuvwxyz"
	const timesRunTest = 100
	var firstChars [timesRunTest]string
	for i := range timesRunTest {
		id := create(t)
		firstChar := charOfString(id, 0)
		if !strings.Contains(atoz, firstChar) {
			t.Errorf("First char is not from a to z. Found %q\n", firstChar)
		}
		firstChars[i] = firstChar
	}
	checkRandomness(firstChars, timesRunTest)
}

func checkRandomness(chars [100]string, timesRunTest int) {
	// TODO bestCaseCharTimes := timesRunTest / 26.0
	// check randomness
}
