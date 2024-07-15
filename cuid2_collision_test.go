package cuid2

import (
	"fmt"
	"math"
	"math/big"
	"slices"
	"sync"
	"testing"
)

func TestCollision(t *testing.T) {
	n := int(math.Pow(float64(7), float64(8))) * 2
	var ids []string
	fmt.Printf("Generating %d unique ids\n", n)

	numPools := 7
	idPoolResponseChan := make(chan *IdPoolResponse, numPools)
	var wg sync.WaitGroup
	for i := 0; i < numPools; i++ {
		wg.Add(1)
		go createIdPool(t, n/numPools, i, idPoolResponseChan, &wg)
	}
	go func() {
		wg.Wait()
		close(idPoolResponseChan)
	}()
	m := sync.Mutex{}
	for v := range idPoolResponseChan {
		m.Lock()
		ids = slices.Concat(ids, v.Ids)
		checkHistogram(t, n/numPools, v.Histogram)
		m.Unlock()
	}
	set := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	if len(set) < len(ids) {
		t.Errorf("Collision detected. len(set) %d, len(ids) %d", len(set), len(ids))
	}
	fmt.Printf("Sample ids %v\n", ids[:10])
}

func checkHistogram(t *testing.T, numberOfIds int, histogram []int64) {
	expectedBinSize := math.Ceil(float64(numberOfIds) / float64(len(histogram)))
	tolerance := 0.05
	minBinSize := math.Round(expectedBinSize * (1 - tolerance))
	maxBinSize := math.Round(expectedBinSize * (1 + tolerance))
	fmt.Printf("For histogram %v minBinSize %f, and maxBinSize %f\n", histogram, minBinSize, maxBinSize)
	for _, n := range histogram {
		if minBinSize > float64(n) || float64(n) > maxBinSize {
			t.Errorf("Histogram outside distribution tolerance")
			break
		}
	}
}

type IdPoolResponse struct {
	Ids       []string
	Numbers   []*big.Int
	Histogram []int64
}

func createIdPool(t *testing.T, max int, poolId int, idPoolResponseChan chan *IdPoolResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	set := make(map[string]struct{}, max)
	for i := 0; i < max; i++ {
		id := CreateId()
		if !IsCuid(id) {
			t.Errorf("The id %s is not a CUID", id)
			break
		}
		set[id] = struct{}{}
		if len(set) < i {
			t.Errorf("Collision at %d. With value %s\n", i, id)
			break
		}
	}
	fmt.Printf("No collisions detected for pool %d\n", poolId)
	ids := make([]string, len(set))
	i := 0
	for key := range set {
		ids[i] = key
		i++
	}
	numbers := make([]*big.Int, len(set))
	for i, id := range ids {
		idWithoutFirstChar := id[1:] // because first char is less random (I guess)
		n := new(big.Int)
		n.SetString(idWithoutFirstChar, 36)
		numbers[i] = n
	}
	bucketCount := 20
	histogram := buildHistogram(numbers, bucketCount)
	fmt.Printf("Histogram created for pool %d\n", poolId)
	idPoolResponseChan <- &IdPoolResponse{
		Ids:       ids,
		Numbers:   numbers,
		Histogram: histogram,
	}
}

func buildHistogram(numbers []*big.Int, bucketCount int) []int64 {
	buckets := make([]int64, bucketCount)
	var counter int64 = 1
	b, _ := big.NewFloat(math.Pow(36, 23)).Int(nil)
	c := big.NewInt(int64(bucketCount))
	a := new(big.Int).Div(b, c)
	f, _ := a.Float64()
	bucketLen, _ := big.NewFloat(math.Round(f)).Int(nil)

	for _, number := range numbers {
		if counter%bucketLen.Int64() == 0 {
			fmt.Printf("Number %d\n", number)
		}
		bucketBigFloat, _ := new(big.Int).Div(number, bucketLen).Float64()
		bucket := int64(math.Floor(bucketBigFloat))
		buckets[bucket] += 1
		counter++
	}
	return buckets
}
