package cuid2

import (
	"golang.org/x/crypto/sha3"
	"math"
	"math/big"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var alphabet [26]string
var counter func() int64
var fingerprint string

// ~22k hosts before 50% chance of initial counter collision
// with a remaining counter range of 9.0e+15 in JavaScript.
const initialCountMax = 476782367
const defaultLength = 24
const bigLength = 32
const base = 36

func randomLetter() string {
	return alphabet[rand.IntN(len(alphabet))]
}

func init() {
	for i := 0; i < 26; i++ {
		alphabet[i] = string([]rune{rune(i + 97)})
	}
	counter = createCounter(int64(math.Floor(rand.Float64() * initialCountMax)))
	fingerprint = createFingerprint()
}

func bufToBigInt(buf [64]byte) string {
	value := new(big.Int)
	value.SetBytes(buf[:])
	return value.Text(base)
}

func hash(input string) string {
	sha3Val := sha3.Sum512([]byte(input))
	hash := bufToBigInt(sha3Val)
	return hash[1:]
}

func createFingerprint() string {
	globals := "some global a sdfas dfa sdfa sdf asdfa sdfasdf"
	sourceString := globals + createEntropy(bigLength)

	return hash(sourceString)[:bigLength]
}

func createCounter(count int64) func() int64 {
	return func() int64 {
		count = count + 1
		return count
	}
}

func createEntropy(length int) string {
	var entropy strings.Builder
	for entropy.Len() < length {
		entropy.WriteString(strconv.FormatInt(int64(math.Floor(rand.Float64()*36)), base))
	}
	return entropy.String()
}

func CreateId() string {
	firstLetter := randomLetter()

	// If we're lucky, the base 36 conversion calls may reduce hashing rounds
	// by shortening the input to the hash function a little.
	timeString := strconv.FormatInt(time.Now().UnixMilli(), base)
	count := strconv.FormatInt(counter(), base)

	// The salt should be long enough to be globally unique across the full
	// length of the hash. For simplicity, we use the same length as the
	// intended id output.
	salt := createEntropy(defaultLength)
	hashInput := timeString + salt + count + fingerprint
	hash := hash(hashInput)

	cuid2 := firstLetter + hash[1:defaultLength]
	return cuid2
}

func IsCuid(id string) bool {
	minLength := 2
	maxLength := bigLength

	length := len(id)
	matched, err := regexp.MatchString("^[0-9a-z]+$", id)
	if err != nil {
		return false
	}
	return length >= minLength && length <= maxLength && matched
}
