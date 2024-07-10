package cuid2

import (
	"math/rand/v2"
)

var alphabet [26]string

func randomLetter() string {
	return alphabet[rand.IntN(len(alphabet))]
}

func init() {
	for i := 0; i < 26; i++ {
		alphabet[i] = string([]rune{rune(i + 97)})
	}
}

func CreateId() string {
	firstLetter := randomLetter()
	return firstLetter
}
