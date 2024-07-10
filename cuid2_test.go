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

func TestIsCuidFalse(t *testing.T) {
	tests := []string{
		"", "1", "1", "asdf98923jhf90283jh02983hjf02983fh",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			isCuid := IsCuid(tt)
			if isCuid {
				t.Errorf("The input %s should not be a cuid, but it says so.", tt)
			}
		})
	}
}

func TestIsCuidWithValidIds(t *testing.T) {
	// These ids are generated by original library
	tests := []string{
		"tra51en4jnteg5yfch4ww56u",
		"hwfvvcgfp6unxxzcbwbewil8",
		"bu37v7q6sgrgzhagarwz0lpt",
		"vkdfwossmrykise0fpd4bq8x",
		"i7ls2o9hx6rlcjc9ha9tjyu7",
		"l6l6wbqcsiwxx7jvaev1m8ra",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			isCuid := IsCuid(tt)
			if !isCuid {
				t.Errorf("The input %s should be a cuid, but it says it is not.", tt)
			}
		})
	}
}

func FuzzTestRandom(f *testing.F) {
	noOfIdsToGenerate := 1000
	for range noOfIdsToGenerate {
		f.Add(CreateId())
	}

	f.Fuzz(func(t *testing.T, id string) {
		if !IsCuid(id) {
			t.Errorf("Either id or isCuid is wrong for: %s", id)
		}
	})
}
