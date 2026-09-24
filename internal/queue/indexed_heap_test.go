package queue

import "testing"

func TestIndexedHeap_PutAndPeek(t *testing.T) {
	h := NewIndexedHeap(func(a int, b int) bool {
		return a < b
	})
	h.Put("job-1", 50)
	h.Put("job-2", 10)
	h.Put("job-3", 20)
	h.Put("job-4", 30)

	id, value, ok := h.Peek()
	if !ok {
		t.Fatal("Expected Heap to contain an item")
	}
	if id != "job-2" {
		t.Fatalf("got id %q, want %q", id, "job-2")
	}
	if value != 10 {
		t.Fatalf("got value %d, want 10", value)
	}
}

func TestIndexedHeap_PopMin(t *testing.T) {
	h := NewIndexedHeap(func(a int, b int) bool {
		return a < b
	})
	h.Put("job-1", 50)
	h.Put("job-2", 10)
	h.Put("job-3", 20)
	h.Put("job-4", 30)

	expected := []int{10, 20, 30, 50}

	for _, want := range expected {
		_, got, ok := h.PopMin()

		if !ok {
			t.Fatalf("expected value %d, but heap was empty", want)
		}

		if got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
	}

	if h.Len() != 0 {
		t.Fatalf("got length %d, want 0", h.Len())
	}
}

func TestIndexedHeap_Has(t *testing.T) {
	h := NewIndexedHeap(func(a, b int) bool {
		return a < b
	})

	h.Put("job-1", 100)

	if !h.Has("job-1") {
		t.Fatal("expected job-1 to exist")
	}

	if h.Has("job-2") {
		t.Fatal("did not expect job-2 to exist")
	}
}

func TestIndexedHeap_Delete(t *testing.T) {
	h := NewIndexedHeap(func(a int, b int) bool {
		return a < b
	})

	h.Put("job-1", 10)
	h.Put("job-2", 20)
	h.Put("job-3", 30)

	ok := h.Delete("job-2")

	if !ok {
		t.Fatal("expected Delete to return true")
	}

	if h.Has("job-2") {
		t.Fatal("job-2 should have been deleted")
	}

	if h.Len() != 2 {
		t.Fatalf("got length %d, want 2", h.Len())
	}

	_, value, _ := h.PopMin()

	if value != 10 {
		t.Fatalf("got %d, want 10", value)
	}

	_, value, _ = h.PopMin()

	if value != 30 {
		t.Fatalf("got %d, want 30", value)
	}

}

func TestIndexedHeap_Update(t *testing.T) {
	h := NewIndexedHeap(func(a int, b int) bool {
		return a < b
	})

	h.Put("job-1", 50)
	h.Put("job-2", 20)
	h.Put("job-3", 30)

	h.Put("job-1", 10) //update

	id, value, ok := h.Peek()

	if !ok {
		t.Fatal("expected heap to contain an item")
	}

	if id != "job-1" {
		t.Fatalf("got id %q, want job-1", id)
	}

	if value != 10 {
		t.Fatalf("got value %d, want 10", value)
	}
}

func TestIndexedHeap_DeleteMissing(t *testing.T) {
	h := NewIndexedHeap(func(a int, b int) bool {
		return a < b
	})

	if h.Delete("does-not-exist") {
		t.Fatal("expected Delete to return false")
	}
}

func TestIndexedHeap_Empty(t *testing.T) {
	h := NewIndexedHeap(func(a, b int) bool {
		return a < b
	})

	_, _, ok := h.Peek()
	if ok {
		t.Fatal("Peek should return false on empty heap")
	}

	_, _, ok = h.PopMin()
	if ok {
		t.Fatal("PopMin should return false on empty heap")
	}

	if h.Len() != 0 {
		t.Fatalf("got length %d, want 0", h.Len())
	}
}

func assertHeapInvariant[T any](t *testing.T, h *IndexedHeap[T]) {
	t.Helper() 
	for i, n := range h.items {
		if n == nil {
			t.Fatalf("items[%d] is nil", i)
		}
		if n.idx != i {
			t.Fatalf(
				"node %q has idx %d, but is at items[%d]",
				n.id,
				n.idx,
				i,
			)
		}
		if h.byID[n.id] != n {
			t.Fatalf(
				"byID[%q] does not point to the same node",
				n.id,
			)
		}
	}

}
//t.Helper() //This function is a test helper.
//If something fails inside it, report the error 
// as coming from the function that called this helper, 
// rather than from this helper itself


func TestIndexedHeap_Invariant(t *testing.T) {
	h := NewIndexedHeap(func(a, b int) bool {
		return a < b
	})
	h.Put("a", 50)
	assertHeapInvariant(t, h)

	h.Put("b", 10)
	assertHeapInvariant(t, h)

	h.Put("c", 30)
	assertHeapInvariant(t, h)

	h.Put("d", 20)
	assertHeapInvariant(t, h)

	h.Delete("b")
	assertHeapInvariant(t, h)

	h.Put("a", 5)
	assertHeapInvariant(t, h)

	h.PopMin()
	assertHeapInvariant(t, h)

	h.PopMin()
	assertHeapInvariant(t, h)

}
