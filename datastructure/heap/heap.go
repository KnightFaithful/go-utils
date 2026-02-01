package heap

import (
	"container/heap"
	"github.com/emirpasic/gods/v2/utils"
)

type innerHeap[K comparable] struct {
	list       []K
	Comparator utils.Comparator[K]
}

func (q innerHeap[K]) Len() int            { return len(q.list) }
func (q innerHeap[K]) Swap(i, j int)       { q.list[i], q.list[j] = q.list[j], q.list[i] }
func (q *innerHeap[K]) Push(x interface{}) { q.list = append(q.list, x.(K)) }
func (q innerHeap[K]) Less(i, j int) bool  { return q.Comparator(q.list[i], q.list[j]) < 0 }
func (q *innerHeap[K]) Pop() interface{} {
	old := q.list
	n := len(old)
	x := old[n-1]
	q.list = old[:n-1]
	return x
}

func newHeapWith[K comparable](comparator utils.Comparator[K], list []K) *innerHeap[K] {
	res := &innerHeap[K]{
		list:       list,
		Comparator: comparator,
	}
	heap.Init(res)
	return res
}

type Heap[K comparable] struct {
	inner *innerHeap[K]
}

func NewHeap[K comparable](comparator utils.Comparator[K], list []K) *Heap[K] {
	return &Heap[K]{
		inner: newHeapWith[K](comparator, list),
	}
}

func (h *Heap[K]) Size() int  { return h.inner.Len() }
func (h *Heap[K]) Push(val K) { heap.Push(h.inner, val) }
func (h *Heap[K]) Pop() K {
	if h.Size() == 0 {
		var zero K
		return zero
	}
	res := heap.Pop(h.inner)
	return res.(K)
}
func (h *Heap[K]) Peek() K {
	if h.Size() == 0 {
		var zero K
		return zero
	}
	return h.inner.list[0]
}
