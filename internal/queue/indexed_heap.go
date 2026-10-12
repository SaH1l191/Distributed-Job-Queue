package queue

import "container/heap"

type node[T any] struct {
	id    string
	value T
	idx   int
}


//follows contract of container/heap interface
//implementing push,pop,len,less,swap of heap
type IndexedHeap[T any] struct {
	items []*node[T]
	byID  map[string]*node[T]
	less  func(a T, b T) bool
}
//The heap.Push, heap.Pop, heap.Fix, and heap.
// Remove functions belong to Go's standard library, container/heap. 
// They implement the heap algorithms but rely on your IndexedHeap methods 
// to interact with your data structure.

// Go provides the algorithm; you provide the implementation that satisfies its contract and maintains your indexed-heap-specific state.

func NewIndexedHeap[T any](less func(a T, b T) bool) *IndexedHeap[T] {

	if(less == nil){
		panic("less function cannot be nil")
	}
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

//not a complete minimum-removal
//primitive
func (h *IndexedHeap[T]) Pop() any {
	lastIdx := len(h.items) - 1

	toBeRemoved := h.items[lastIdx]
	h.items[lastIdx] = nil
	h.items = h.items[:lastIdx]

	delete(h.byID, toBeRemoved.id)
	toBeRemoved.idx = -1
	return toBeRemoved

}

func (h *IndexedHeap[T]) Put(id string, value T) { //upsert
	if n, exists := h.byID[id]; exists {
		n.value = value
		heap.Fix(h, n.idx) //standard library uses the node's current index to restore the heap ordering
		return
	}
	n := &node[T]{
		id:    id,
		value: value,
		idx:   -1,
	}
	h.byID[id] = n
	heap.Push(h, n)
	// /internal go heap handles everything
	//Calls h.Push(n) to append the new node.
	// Compares the new node with its parent using h.Less(...).
	// If the heap property is violated, calls h.Swap(...).
	// Repeats until the heap property is restored.

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
	//internal go heap handles heap algorithm 
	// Calls Swap(0, lastIndex) to move the root to the end.
	//Restores heap order using my Less() and Swap() methods
	// Calls my Pop() method to remove that last element.

	return n.id, n.value, true
}

func (h *IndexedHeap[T]) Delete(id string) bool {
	n, exists := h.byID[id]

	if !exists {
		return false
	}

	heap.Remove(h, n.idx)
	// Uses your Swap() to move the target node to the end.
	// Restores heap order using your Less() and Swap() methods where needed.
	// Calls your Pop() to remove the target node from the slice.

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
