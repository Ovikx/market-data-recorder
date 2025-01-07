package timewindow

import (
	"fmt"
	"time"
)

type timeWindowNode[T any] struct {
	timestamp time.Time
	next      *timeWindowNode[T]
	prev      *timeWindowNode[T]
	data      T
}

func (n *timeWindowNode[T]) String() string {
	return fmt.Sprintf("%v", n.data)
}
