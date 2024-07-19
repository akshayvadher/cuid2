package cuid2

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
	"testing"
)

func TestHistogram(t *testing.T) {
	n := 100000
	fmt.Printf("Testing %d unique ids\n", n)
	poolId := rand.IntN(100)
	poolResponse := CreateIdPool(t, n, poolId)
	ids := poolResponse.Ids
	sampleIds := ids[:10]
	fmt.Printf("Sample ids %v\n", sampleIds)
	t.Run("Test collision", func(t *testing.T) {
		CheckCollision(t, ids)
	})
	t.Run("Test char frequency", func(t *testing.T) {
		testCharFrequency(t, n, ids)
	})
	t.Run("Test histogram", func(t *testing.T) {
		testHistogram(t, poolResponse, n)
	})

}

func testCharFrequency(t *testing.T, n int, ids []string) {
	tolerance := 0.1
	idLength := 23
	totalLetters := idLength * n
	base := 36
	expectedBinSize := math.Ceil(float64(totalLetters) / float64(base))
	minBinSize := math.Round(expectedBinSize * (1 - tolerance))
	maxBinSize := math.Round(expectedBinSize * (1 + tolerance))

	// Drop the first character because it will always be a letter, making
	// the letter frequency skewed.
	testIds := make([]string, len(ids))
	for i, id := range ids {
		testIds[i] = id[1:]
	}
	charFrequencies := make(map[string]int)
	for _, id := range testIds {
		chars := strings.Split(id, "")
		for _, char := range chars {
			charFrequencies[char] += 1
		}
	}
	fmt.Println("Testing character frequency...")
	fmt.Printf("expectedBinSize %v\n", expectedBinSize)
	fmt.Printf("minBinSize %v\n", minBinSize)
	fmt.Printf("maxBinSize %v\n", maxBinSize)
	fmt.Printf("charFrequencies %v\n", charFrequencies)
	for k, v := range charFrequencies {
		if float64(v) < minBinSize || float64(v) > maxBinSize {
			t.Errorf("The char %v is out of the expected bin size with value %v\n", k, v)
		}
	}
	if len(charFrequencies) != base {
		t.Errorf("Not all of the chars are presention in ids. Got only %v\n", len(charFrequencies))
	}
}

func testHistogram(t *testing.T, poolResponse *IdPoolResponse, n int) {
	histogram := poolResponse.Histogram
	expectedBinSize := math.Ceil(float64(n) / float64(len(histogram)))
	tolerance := 0.1
	minBinSize := math.Round(expectedBinSize * (1 - tolerance))
	maxBinSize := math.Round(expectedBinSize * (1 + tolerance))
	fmt.Printf("Histogram %v\n", histogram)
	fmt.Printf("expectedBinSize %v\n", expectedBinSize)
	fmt.Printf("minBinSize %v\n", minBinSize)
	fmt.Printf("maxBinSize %v\n", maxBinSize)
	for _, i := range histogram {
		if float64(i) < minBinSize || float64(i) > maxBinSize {
			t.Errorf("Histogram is out of distribution tolerance")
		}
	}
}
