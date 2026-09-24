package queue

import "container/heap"

type node[T any] struct {
	id    string
	value T
	idx   int
}

type IndexedHeap[T any] struct {
	items []*node[T]
	byID  map[string]*node[T]
	less  func(a T, b T) bool
}

func NewIndexedHeap[T any](less func(a T, b T) bool) *IndexedHeap[T] {
	return &IndexedHeap[T]{
		items: make([]*node[T], 0),
		byID:  make(map[string]*node[T]),
		less:  less,
	}
}

func (h *IndexedHeap[T]) Len() int {
	return len(h.items)
}

func (h *IndexedHeap[T]) Less(i int, j int) bool {
	return h.less(h.items[i].value, h.items[j].value)
}

func (h *IndexedHeap[T]) Swap(i int, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
	h.items[i].idx = i
	h.items[j].idx = j
}

func (h *IndexedHeap[T]) Push(x any) {
	newnode := x.(*node[T])
	newnode.idx = len(h.items)
	h.items = append(h.items, newnode)
}

func (h *IndexedHeap[T]) Pop() any {
	lastIdx := len(h.items) - 1

	toBeRemoved := h.items[lastIdx]
	h.items[lastIdx] = nil
	h.items = h.items[:lastIdx]

	delete(h.byID, toBeRemoved.id)
	toBeRemoved.idx = -1
	return toBeRemoved

}

func (h *IndexedHeap[T]) Put(id string, value T) {
	if n, exists := h.byID[id]; exists {
		n.value = value
		heap.Fix(h, n.idx)
		return
	}
	n := &node[T]{
		id:    id,
		value: value,
		idx:   -1,
	}
	h.byID[id] = n
	heap.Push(h, n)

}

func (h *IndexedHeap[T]) Peek() (id string, value T, ok bool) {
	if len(h.items) == 0 {
		var zero T
		return "", zero, false
	}

	n := h.items[0]

	return n.id, n.value, true
}

func (h *IndexedHeap[T]) PopMin() (id string, value T, ok bool) {
	if len(h.items) == 0 {
		var zero T
		return "", zero, false
	}

	n := heap.Pop(h).(*node[T])

	return n.id, n.value, true
}

func (h *IndexedHeap[T]) Delete(id string) bool {
	n, exists := h.byID[id]

	if !exists {
		return false
	}

	heap.Remove(h, n.idx)
	return true
}

func (h *IndexedHeap[T]) Has(id string) bool {
	_, exists := h.byID[id]
	return exists
}

//generic indexedHeap of value<T>
//with custom sort,indexes,def ops like push.pop.peek.delete.upsert

//later can be used with
//inflight_queue,ready_queue,toBeScehduledQueue
