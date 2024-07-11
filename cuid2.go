package cuid2

import (
	"golang.org/x/crypto/sha3"
	"math"
	"math/big"
	"math/rand/v2"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	alphabet           [26]string
	DefaultRandom      func() float64
	DefaultCounter     func() int64
	DefaultFingerprint string
	defaultInit        func() string
)

const (
	// ~22k hosts before 50% chance of initial counter collision
	// with a remaining counter range of 9.0e+15 in JavaScript.
	initialCountMax = 476782367
	defaultLength   = 24
	bigLength       = 32
	base            = 36
)

func randomLetter(random func() float64) string {
	return alphabet[int(math.Floor(random()*float64(len(alphabet))))]
}

func init() {
	for i := 0; i < 26; i++ {
		alphabet[i] = string([]rune{rune(i + 97)})
	}
	DefaultRandom = rand.Float64
	DefaultCounter = createCounter(int64(math.Floor(DefaultRandom() * initialCountMax)))
	DefaultFingerprint = createFingerprint(DefaultRandom)
	defaultInit = Init(DefaultRandom, DefaultCounter, defaultLength, DefaultFingerprint)
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

func createFingerprint(random func() float64) string {
	host, _ := os.Hostname()
	userHome, _ := os.UserHomeDir()
	pid := os.Getpid()
	globals := host + userHome + string(rune(pid))
	sourceString := globals + createEntropy(bigLength, random)

	return hash(sourceString)[:bigLength]
}

func createCounter(count int64) func() int64 {
	return func() int64 {
		count = count + 1
		return count
	}
}

func createEntropy(length int, random func() float64) string {
	var entropy strings.Builder
	for entropy.Len() < length {
		entropy.WriteString(strconv.FormatInt(int64(math.Floor(random()*base)), base))
	}
	return entropy.String()
}

type Options interface {
	random() float64
}

func Init(random func() float64, counter func() int64, length int, fingerprint string) func() string {
	return func() string {
		firstLetter := randomLetter(random)

		// If we're lucky, the base 36 conversion calls may reduce hashing rounds
		// by shortening the input to the hash function a little.
		timeString := strconv.FormatInt(time.Now().UnixMilli(), base)
		count := strconv.FormatInt(counter(), base)

		// The salt should be long enough to be globally unique across the full
		// length of the hash. For simplicity, we use the same length as the
		// intended id output.
		salt := createEntropy(length, random)
		hashInput := timeString + salt + count + fingerprint
		hash := hash(hashInput)

		cuid2 := firstLetter + hash[1:length]
		return cuid2
	}
}

func CreateId() string {
	return defaultInit()
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
