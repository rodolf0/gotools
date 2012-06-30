/*
 * A container that keeps elements sorted and uses arrays as underlying storage.
 * The main usage is to find where a new element would fit within a range.
 */
package bisect

type Elem interface {
	Less(other Elem) bool
}

// return the insertion point to keep list in order, if key is already
// present, the insertion point will be before equal elements
func BisectLeft(list []Elem, key Elem) int {
	if len(list) == 0 || !list[0].Less(key) {
		return 0
	} else if list[len(list)-1].Less(key) {
		return len(list)
	}

	if mid := len(list) / 2; !list[mid].Less(key) {
		return BisectLeft(list[:mid], key)
	} else {
		return BisectLeft(list[mid:], key) + mid
	}
	panic("BisectLeft bug")
	return 0
}

// return the insertion point to keep list in order, if key is already
// present, the insertion point will be after the existing elements
func Bisect(list []Elem, key Elem) int {
	if len(list) == 0 || key.Less(list[0]) {
		return 0
	} else if !key.Less(list[len(list)-1]) {
		return len(list)
	}

	if mid := len(list) / 2; key.Less(list[mid]) {
		return Bisect(list[:mid], key)
	} else {
		return Bisect(list[mid:], key) + mid
	}
	panic("Bisect bug")
	return 0
}

// insert value x into the list in sorted order (after equal elements)
func Insort(list []Elem, x Elem) []Elem {
	var idx = Bisect(list, x)
	return append(list[:idx], append([]Elem{x}, list[idx:]...)...)
}

// insert value x into the list in sorted order (before equal elements)
func InsortLeft(list []Elem, x Elem) []Elem {
	var idx = BisectLeft(list, x)
	return append(list[:idx], append([]Elem{x}, list[idx:]...)...)
}

// find element x inside list, flag indicates if element was found
func Index(list []Elem, x Elem) (int, bool) {
	var idx = BisectLeft(list, x)
	if idx >= len(list) || list[idx].Less(x) || x.Less(list[idx]) {
		return 0, false
	}
	return idx, true
}

// find and remove x
func Remove(list []Elem, x Elem) []Elem {
	var idx, found = Index(list, x)
	if !found {
		panic("Removing non existent element")
	}
	return append(list[:idx], list[idx+1:]...)
}
