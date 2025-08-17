package algorithms

import (
	"fmt"
	"testing"
)

func TestFindPeaksValleys(t *testing.T) {
	highList := []float64{1, 10, 2, 6, 4, 5, 3, 8, 5, 7, 3, 10, 5}
	lowList := []float64{0, 8, 0, 4, 2, 3, 1, 6, 3, 5, 1, 8, 3}
	fmt.Println(len(highList), len(lowList))
	fmt.Println("----------")
	peaks, valleys, err := FindPeaksValleys(highList, lowList)
	fmt.Println(peaks, valleys)
	if err != nil {
		t.Error(err)
	}
}
