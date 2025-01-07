// Package timewindow implements a list that removes any nodes that are older than the head node (which has the most recent timestamp) by the provided range.
package timewindow

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// This collection is implemented as a doubly-linked list. This collection is a collection buffered not by size, but by time.
// Nodes are sorted by timestamp (head is the most recent, tail is the oldest).
// This list expects that inserted nodes are roughly new.
type timeWindow[T any] struct {
	mutex     *sync.Mutex
	length    int
	head      *timeWindowNode[T]
	tail      *timeWindowNode[T]
	timeRange time.Duration
}

func New[T any](timeRange time.Duration) *timeWindow[T] {
	return &timeWindow[T]{timeRange: timeRange, mutex: &sync.Mutex{}}
}

// Insert adds the input data and its associated timestamp to the time window. This function also returns a slice containing the purged data and whether the input data was added.
func (l *timeWindow[T]) Insert(data T, timestamp time.Time) ([]T, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	// Create the node
	node := &timeWindowNode[T]{data: data, timestamp: timestamp}

	// Insert the node into the correct location
	if l.length == 0 { // List is empty
		l.head = node
		l.tail = node
	} else {
		// Don't insert if it's already invalid
		if l.head.timestamp.Sub(timestamp) > l.timeRange {
			return []T{}, false
		}

		curr := l.head
		for curr != nil && node.timestamp.Before(curr.timestamp) {
			curr = curr.next
		}

		if curr == nil { // Went to end of list
			l.tail.next = node
			node.prev = l.tail
			l.tail = node
		} else if curr == l.head { // Still at head
			l.head.prev = node
			node.next = l.head
			l.head = node
		} else { // Somewhere in the middle
			prev := curr.prev
			node.prev = prev
			node.next = curr
			curr.prev = node
			prev.next = node
		}
	}
	l.length++
	return l.purge(), true
}

// Removes any tailing nodes that are too old. Returns the nodes that were removed.
func (l *timeWindow[T]) purge() []T {
	removed := make([]T, 0)
	for curr := l.tail; curr != nil && l.head.timestamp.Sub(curr.timestamp) > l.timeRange; curr = curr.prev {
		prev := l.tail.prev
		prev.next = nil
		l.tail.prev = nil
		l.tail = prev
		l.length--
		removed = append(removed, curr.data)
	}

	return removed
}

func (l *timeWindow[T]) String() string {
	nodeStrs := make([]string, 0)
	for curr := l.head; curr != nil; curr = curr.next {
		nodeStrs = append(nodeStrs, fmt.Sprintf("%v", curr))
	}
	return "[" + strings.Join(nodeStrs, " ") + "]"
}

func (l *timeWindow[T]) Length() int {
	return l.length
}
