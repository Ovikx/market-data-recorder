package timewindow_test

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"testing"
	"time"

	"github.com/Ovikx/market-data-recorder/internal/datastructure/timewindow"
)

func TestBasicInsertNoOverflow(t *testing.T) {
	t.Parallel()

	l := timewindow.New[int](10 * time.Hour)
	l.Insert(0, time.Now())
	l.Insert(1, time.Now().Add(-2*time.Hour))
	l.Insert(2, time.Now().Add(-time.Hour))
	l.Insert(3, time.Now().Add(time.Hour))

	if fmt.Sprintf("%v", l) != "[3 0 2 1]" {
		t.Errorf("expected list to be [3 0 2 1] but got %v", l)
	}
}

func TestBasicInsertOverflow(t *testing.T) {
	t.Parallel()

	l := timewindow.New[int](time.Hour)
	l.Insert(0, time.Now())
	l.Insert(1, time.Now().Add(-2*time.Hour))
	l.Insert(2, time.Now().Add(-time.Hour/2))
	l.Insert(3, time.Now().Add(-3*time.Hour))

	if fmt.Sprintf("%v", l) != "[0 2]" {
		t.Errorf("expected list to be [0 2] but got %v", l)
	}
}

type randsource struct{}

func (r randsource) Uint64() uint64 {
	return 1337
}

func TestInsertOverflow(t *testing.T) {
	t.Parallel()

	times := make([]int64, 0)
	t1 := time.Now()
	t2 := t1.Add(time.Hour * 2)

	// Create the timestamps
	for i := 0; i < 10000; i++ {
		times = append(times, int64(rand.Float64()*float64(t2.UnixMilli()-t1.UnixMilli())+float64(t1.UnixMilli())))
	}

	// Insert into time window
	l := timewindow.New[int64](time.Hour)
	for _, t := range times {
		l.Insert(t, time.UnixMilli(t))
	}

	// Sort descending
	slices.SortFunc(times, func(a, b int64) int {
		return int(b - a)
	})

	// Clean the input array
	firstTime := times[0]
	cleaned := make([]int64, 0)
	for _, t := range times {
		if time.UnixMilli(firstTime).Sub(time.UnixMilli(t)) <= time.Hour {
			cleaned = append(cleaned, t)
		}
	}

	if fmt.Sprintf("%v", cleaned) != fmt.Sprintf("%v", l) {
		t.Error("arrays didn't match..")
	}
}
