package main

// An PriorityQueue is a min-heap of ints.
type PriorityQueue []int

func (h PriorityQueue) Len() int           { return len(h) }
func (h PriorityQueue) Less(i, j int) bool { return h[i] < h[j] }
func (h PriorityQueue) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PriorityQueue) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *PriorityQueue) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Peek returns the minimum element without removing it from the heap.
func (h PriorityQueue) Peek() (int, bool) {
	if len(h) == 0 {
		return 0, false
	}
	return h[0], true
}

// FindIndex returns the index of the first occurrence of val in the heap.
// It panics if val is not found.
func (h PriorityQueue) FindIndex(val int) int {
	for i, v := range h {
		if v == val {
			return i
		}
	}
	panic("Value not found")
}
