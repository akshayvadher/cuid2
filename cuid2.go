package cuid2

import (
	"errors"
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
	alphabet           = [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	DefaultRandom      = rand.Float64
	DefaultCounter     func() int64
	DefaultFingerprint string
	defaultInit        func() string
	envVariableKeys    string
	cuidRegex          = regexp.MustCompile("^[a-z][0-9a-z]+$")
)

const (
	// ~22k hosts before 50% chance of initial counter collision
	// with a remaining counter range of 9.0e+15 in JavaScript.
	initialCountMax = 476782367
	defaultLength   = 24
	bigLength       = 32
	base            = 36
)

func init() {
	envVariableKeys = strings.Join(getEnvVariableKeys(), "_")
	DefaultCounter = createCounter(int64(DefaultRandom() * initialCountMax))
	DefaultFingerprint = createFingerprint(DefaultRandom)
	defaultInit = Init(DefaultRandom, DefaultCounter, defaultLength, DefaultFingerprint)
}

func getEnvVariableKeys() []string {
	var ek []string
	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			ek = append(ek, e[:i])
		}
	}
	return ek
}

func randomLetter(random func() float64) string {
	return alphabet[int(random()*float64(len(alphabet)))]
}

func bufToBigInt(buf [64]byte) string {
	value := new(big.Int)
	value.SetBytes(buf[:])
	return value.Text(base)
}

func hash(input string) string {
	sha3Val := sha3.Sum512([]byte(input))
	hash := bufToBigInt(sha3Val)
	// Drop the first character because it will bias the histogram
	// to the left.
	return hash[1:]
}

func createFingerprint(random func() float64) string {
	pid := os.Getpid()
	globals := envVariableKeys + strconv.Itoa(pid)
	sourceString := globals + createEntropy(bigLength, random)

	return hash(sourceString)[:bigLength]
}

func createCounter(start int64) func() int64 {
	count := start
	return func() int64 {
		count++
		return count
	}
}

func createEntropy(length int, random func() float64) string {
	var entropy strings.Builder
	entropy.Grow(length)
	for entropy.Len() < length {
		entropy.WriteString(strconv.FormatInt(int64(math.Floor(random()*base)), base))
	}
	return entropy.String()
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

func CreateIdOf(len int) (string, error) {
	if len < 2 || len > bigLength {
		return "", errors.New("len should be between 2 and 32")
	}
	return Init(DefaultRandom, DefaultCounter, len, DefaultFingerprint)(), nil
}

func IsCuid(id string) bool {
	minLength := 2
	maxLength := bigLength

	length := len(id)
	matched := cuidRegex.MatchString(id)
	return length >= minLength && length <= maxLength && matched
}
