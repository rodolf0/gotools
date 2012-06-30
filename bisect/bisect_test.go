package bisect

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type MyInt int

func (i MyInt) Less(other Elem) bool {
	return i < other.(MyInt)
}

// helper function to check that elements are always in order
func inOrder(list []Elem) bool {
	for i := 1; i < len(list); i++ {
		if list[i].Less(list[i-1]) {
			return false
		}
	}
	return true
}

func testInsertion(t *testing.T, insert func([]Elem, Elem) []Elem) {
	var list = []Elem{}
	// insert some numbers and verify they're in order
	var n = 1000 + int(20000.0*rand.Float32())
	for i := 0; i < n; i++ {
		list = insert(list, MyInt(10000.0*rand.Float32()))
		if !inOrder(list) {
			t.FailNow()
		}
	}
	// check that the list has the expected number of elements
	if len(list) != n {
		t.FailNow()
	}
}

func TestInsort(t *testing.T) {
	testInsertion(t, Insort)
}

func TestInsortLeft(t *testing.T) {
	testInsertion(t, InsortLeft)
}

func TestBisect(t *testing.T) {
	var list = []Elem{}
	var n = 1000 + int(2000.0*rand.Float32())
	for i := 0; i < n; i++ {
		list = Insort(list, MyInt(100.0*rand.Float32()))
	}

	var m = 10000 + int(10000.0*rand.Float32())
	for i := 0; i < m; i++ {
		// choose among a small range of values to have equal elements
		var r = MyInt(100.0 * rand.Float32())
		var idx = Bisect(list, r)
		// prev elements must be smaller or equal
		if idx > 0 && r.Less(list[idx-1]) {
			t.Fatalf("Prev element (%v) not LE than %v", list[idx-1], r)
		}
		// check that elements after insertion point are bigger
		if idx < len(list)-1 && !r.Less(list[idx+1]) {
			t.Fatalf("Next element (%v) not GT than %v", list[idx+1], r)
		}
	}
}

func TestBisectLeft(t *testing.T) {
	var list = []Elem{}
	var n = 1000 + int(2000.0*rand.Float32())
	for i := 0; i < n; i++ {
		list = Insort(list, MyInt(100.0*rand.Float32()))
	}

	var m = 10000 + int(10000.0*rand.Float32())
	for i := 0; i < m; i++ {
		var r = MyInt(100.0 * rand.Float32())
		var idx = BisectLeft(list, r)
		// check that elements before insertion point are smaller
		if idx > 0 && !list[idx-1].Less(r) {
			t.Fatalf("Prev element (%v) not LT than %v", list[idx-1], r)
		}
		// check that elements after insertion point are equal or greater
		if idx < len(list)-1 && list[idx+1].Less(r) {
			t.Fatalf("Next element (%v) not GE than %v", list[idx+1], r)
		}
	}
}

func TestIndex(t *testing.T) {
	var list = []Elem{}
	if idx, found := Index(list, MyInt(23)); found {
		t.Fatalf("Found non existent element (%v)", idx)
	}
	list = Insort(list, MyInt(32))
	if _, found := Index(list, MyInt(32)); !found {
		t.Fatalf("Can't find existent element 32 in %v", list)
	}
	list = Insort(list, MyInt(-13))
	list = Insort(list, MyInt(93))
	list = Insort(list, MyInt(3))
	if idx, _ := Index(list, MyInt(32)); idx != 2 {
		t.Fatal("Element has wrong index")
	}
}


func TestRemove(t *testing.T) {
	var list = []Elem{}
	var n = 1000 + int(2000.0*rand.Float32())
	for i := 0; i < n; i++ {
		list = Insort(list, MyInt(100.0*rand.Float32()))
	}
	for i := 0; i < 1000; i++ {
		list = Remove(list, list[int(float32(len(list))*rand.Float32())])
		if len(list) != n - i - 1 {
			t.Fatalf("Expected list length %v, got %v", n-i-1, len(list))
		}
		if !inOrder(list) {
			t.FailNow()
		}
	}
}
