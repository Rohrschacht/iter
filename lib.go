package iter

import "fmt"

// Iterator can be used to process data in a pipeline pattern.
type Iterator[T any] chan T

// Pair is used as a helper when an Iterator has to hold multiple values.
type Pair[T, K any] struct {
	X T
	Y K
}

// FromChan creates an Iterator from a channel.
func FromChan[T any](c chan T) Iterator[T] {
	return c
}

// FromSlice creates an Iterator over the given slice.
func FromSlice[T any](slice []T) Iterator[T] {
	it := make(chan T)
	go func() {
		defer close(it)
		for _, v := range slice {
			it <- v
		}
	}()
	return it
}

// FromMap creates an Iterator of Pairs that contain key and value of the given map.
func FromMap[T comparable, K any](m map[T]K) Iterator[Pair[T, K]] {
	it := make(chan Pair[T, K])
	go func() {
		defer close(it)
		for key, v := range m {
			it <- Pair[T, K]{X: key, Y: v}
		}
	}()
	return it
}

// FromMapKeys creates an Iterator over the keys of the given map.
func FromMapKeys[T comparable, K any](m map[T]K) Iterator[T] {
	it := make(chan T)
	go func() {
		defer close(it)
		for key := range m {
			it <- key
		}
	}()
	return it
}

// FromMapValues creates an Iterator over the values of the given map.
func FromMapValues[K comparable, T any](m map[K]T) Iterator[T] {
	it := make(chan T)
	go func() {
		defer close(it)
		for _, v := range m {
			it <- v
		}
	}()
	return it
}

// Collect consumes the Iterator, returning a slice of all its elements.
func (it Iterator[T]) Collect() []T {
	var slice []T
	for v := range it {
		slice = append(slice, v)
	}
	return slice
}

// Filter uses the given function to determine whether elements should continue through the pipeline.
func (it Iterator[T]) Filter(f func(T) bool) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		for v := range it {
			if f(v) {
				newIter <- v
			}
		}
	}()
	return newIter
}

// Map applies the given function to all elements going through the pipeline.
func (it Iterator[T]) Map(f func(T) T) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		for v := range it {
			newIter <- f(v)
		}
	}()
	return newIter
}

// MapInto applies the given function to all elements and allows for the type to change.
func MapInto[T, K any](it Iterator[T], f func(T) K) Iterator[K] {
	newIter := make(chan K)
	go func() {
		defer close(newIter)
		for v := range it {
			newIter <- f(v)
		}
	}()
	return newIter
}

// Skip skips the first n elements of the Iterator.
//
// n can be larger than the number of elements in the Iterator, which will empty it.
func (it Iterator[T]) Skip(n uint) Iterator[T] {
	for i := uint(0); i < n; i++ {
		<-it
	}
	return it
}

// Take takes the first n elements of the Iterator.
//
// All elements after the first n elements will be discarded.
func (it Iterator[T]) Take(n uint) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		for i := uint(0); i < n; i++ {
			v, more := <-it
			if !more {
				return
			}
			newIter <- v
		}
	}()
	return newIter
}

// Nth returns a pointer to the element at position n.
//
// If there are fewer than n elements in the Iterator, nil is returned.
func (it Iterator[T]) Nth(n uint) *T {
	for i := uint(0); i < n-1; i++ {
		_, ok := <-it
		if !ok {
			return nil
		}
	}
	v, ok := <-it
	if !ok {
		return nil
	}
	return &v
}

// Count consumes the Iterator and returns its number of elements.
func (it Iterator[T]) Count() uint {
	c := uint(0)
	for range it {
		c++
	}
	return c
}

// Last returns the last element of the Iterator, consuming it in the process.
func (it Iterator[T]) Last() T {
	var l T
	for v := range it {
		l = v
	}
	return l
}

// StepBy advances the Iterator by n elements every time something is taken.
func (it Iterator[T]) StepBy(n uint) Iterator[T] {
	if n == 0 {
		return nil
	}

	newIter := make(chan T)
	go func() {
		defer close(newIter)
		for {
			v, more := <-it
			if !more {
				return
			}
			newIter <- v

			for i := uint(0); i < n-1; i++ {
				_, more := <-it
				if !more {
					return
				}
			}
		}
	}()
	return newIter
}

// Chain creates a new Iterator which returns the elements of both Iterators.
func (it Iterator[T]) Chain(other Iterator[T]) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		for v := range it {
			newIter <- v
		}
		for v := range other {
			newIter <- v
		}
	}()
	return newIter
}

// Intersperse inserts the separator sep between each element of the Iterator.
func (it Iterator[T]) Intersperse(sep T) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		v, more := <-it
		if more {
			newIter <- v
		} else {
			return
		}
		for {
			v, more := <-it
			if more {
				newIter <- sep
				newIter <- v
			} else {
				return
			}
		}
	}()
	return newIter
}

// ForEach executes the given function for each element of the Iterator.
func (it Iterator[T]) ForEach(f func(T)) {
	for v := range it {
		f(v)
	}
}

// Zip creates a new Iterator that contains Pairs containing the elements of both Iterators.
//
// If one of the input Iterators is shorter than the other one, the new Iterator
// will stop at that point.
func Zip[T, K any](it Iterator[T], other Iterator[K]) Iterator[Pair[T, K]] {
	newIter := make(chan Pair[T, K])
	go func() {
		defer close(newIter)
		for {
			v1, ok1 := <-it
			if !ok1 {
				return
			}
			v2, ok2 := <-other
			if !ok2 {
				return
			}
			newIter <- Pair[T, K]{X: v1, Y: v2}
		}
	}()
	return newIter
}

// SkipWhile discards all elements until the condition of the given function is met once.
func (it Iterator[T]) SkipWhile(f func(T) bool) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		for v := range it {
			if !f(v) {
				newIter <- v
				break
			}
		}
		for v := range it {
			newIter <- v
		}
	}()
	return newIter
}

// TakeWhile takes elements until the condition of the given function is false once.
func (it Iterator[T]) TakeWhile(f func(T) bool) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		for v := range it {
			if f(v) {
				newIter <- v
			} else {
				return
			}
		}
	}()
	return newIter
}

// Inspect applies the given function on each element while the Iterator is consumed.
//
// This is helpful for debugging, see the example.
func (it Iterator[T]) Inspect(f func(T)) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		for v := range it {
			f(v)
			newIter <- v
		}
	}()
	return newIter
}

// Partition splits the contents of the iterator based on the condition defined in the given function.
//
// Two slices are returned. The first slice contains all elements of the Iterator
// for which f evaluated to true. The second slice contains all elements for
// which f evaluated to false.
func (it Iterator[T]) Partition(f func(T) bool) ([]T, []T) {
	var yes []T
	var no []T
	for v := range it {
		if f(v) {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	return yes, no
}

// Fold applies the given function to all elements, folding them into the given accumulator.
func (it Iterator[T]) Fold(acc T, f func(T, T) T) T {
	for v := range it {
		acc = f(acc, v)
	}
	return acc
}

// Reduce folds the Iterator using the given function, using the first element as the initial accumulator.
//
// Reduce returns a pointer for the accumulated value. If the Iterator is empty, this will be nil.
func (it Iterator[T]) Reduce(f func(T, T) T) *T {
	acc, ok := <-it
	if !ok {
		return nil
	}
	for v := range it {
		acc = f(acc, v)
	}
	return &acc
}

// All checks whether the given condition is true for all elements.
func (it Iterator[T]) All(f func(T) bool) bool {
	for v := range it {
		if !f(v) {
			return false
		}
	}
	return true
}

// Any checks whether there exists one element for which the given condition is true.
func (it Iterator[T]) Any(f func(T) bool) bool {
	for v := range it {
		if f(v) {
			return true
		}
	}
	return false
}

// Find returns a pointer to the first element for which the given condition is true.
//
// If no such element exists, nil is returned.
func (it Iterator[T]) Find(f func(T) bool) *T {
	for v := range it {
		if f(v) {
			return &v
		}
	}
	return nil
}

// Position returns the position of the first element for which the given condition is true as a pointer.
//
// If no such element exists, nil is returned.
func (it Iterator[T]) Position(f func(T) bool) *uint {
	p := uint(0)
	for v := range it {
		p++
		if f(v) {
			return &p
		}
	}
	return nil
}

// Interleave creates a new Iterator that alternates between the two given Iterators.
func (it Iterator[T]) Interleave(other Iterator[T]) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		for {
			v1, ok1 := <-it
			if ok1 {
				newIter <- v1
			}
			v2, ok2 := <-other
			if ok2 {
				newIter <- v2
			}
			if !ok1 && !ok2 {
				return
			}
		}
	}()
	return newIter
}

// InterleaveShortest creates a new Iterator that alternates between the two given Iterators until at least one of them runs out.
func (it Iterator[T]) InterleaveShortest(other Iterator[T]) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		for {
			v1, ok1 := <-it
			if !ok1 {
				return
			}
			newIter <- v1
			v2, ok2 := <-other
			if !ok2 {
				return
			}
			newIter <- v2
		}
	}()
	return newIter
}

// GroupBy returns a list of slices, which elements are grouped by the given condition.
func (it Iterator[T]) GroupBy(f func(T) bool) [][]T {
	var result [][]T
	var lastState *bool
	var currentChunk []T
	for v := range it {
		if lastState == nil {
			state := f(v)
			lastState = &state
			currentChunk = append(currentChunk, v)
		} else {
			state := f(v)
			if state == *lastState {
				currentChunk = append(currentChunk, v)
			} else {
				*lastState = state
				result = append(result, currentChunk)
				currentChunk = []T{v}
			}
		}
	}
	result = append(result, currentChunk)
	return result
}

// Chunks returns a list of slices containing at most n elements of the original Iterator.
func (it Iterator[T]) Chunks(n uint) [][]T {
	var result [][]T
	var currentChunk []T
Loop:
	for {
		for i := uint(0); i < n; i++ {
			v, ok := <-it
			if !ok {
				break Loop
			}
			currentChunk = append(currentChunk, v)
		}
		result = append(result, currentChunk)
		currentChunk = nil
	}
	if currentChunk != nil {
		result = append(result, currentChunk)
	}
	return result
}

// Windows returns all overlapping subslices of length n of the original Iterator.
func (it Iterator[T]) Windows(n uint) [][]T {
	var result [][]T
	var currentWindow []T
	for i := uint(0); i < n; i++ {
		v, ok := <-it
		if !ok {
			result = append(result, currentWindow)
			return result
		}
		currentWindow = append(currentWindow, v)
	}
	result = append(result, currentWindow)
	newWindow := make([]T, n)
	copy(newWindow, currentWindow)
	currentWindow = newWindow
	for {
		v, ok := <-it
		if !ok {
			return result
		}
		for i := uint(0); i < n-1; i++ {
			currentWindow[i] = currentWindow[i+1]
		}
		currentWindow[n-1] = v
		result = append(result, currentWindow)
		newWindow := make([]T, n)
		copy(newWindow, currentWindow)
		currentWindow = newWindow
	}
}

// CartesianProduct returns an Iterator over the cartesian product of both given Iterators.
func CartesianProduct[T, K any](it Iterator[T], other Iterator[K]) Iterator[Pair[T, K]] {
	newIter := make(chan Pair[T, K])
	go func() {
		defer close(newIter)
		var elementBuffer []K
		for v := range it {
			if elementBuffer == nil {
				for vo := range other {
					elementBuffer = append(elementBuffer, vo)
					newIter <- Pair[T, K]{X: v, Y: vo}
				}
			} else {
				for _, vo := range elementBuffer {
					newIter <- Pair[T, K]{X: v, Y: vo}
				}
			}
		}
	}()
	return newIter
}

// Dedup removes duplicates from sections of consecutive elements determined by the given condition.
func (it Iterator[T]) Dedup(f func(T, T) bool) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		var lastElem *T
		for v := range it {
			if lastElem == nil {
				cp := v
				lastElem = &cp
				newIter <- v
			} else {
				if !f(*lastElem, v) {
					newIter <- v
				}
				*lastElem = v
			}
		}
	}()
	return newIter
}

// Unique produces an Iterator that returns unique elements from the given Iterator determined by the given condition.
//
// Since Iterator can take any type, f has to convert the element type into a
// comparable type. If your type is already comparable, it is enough to just
// return it in the closure. See the example.
func Unique[T any, K comparable](it Iterator[T], f func(T) K) Iterator[T] {
	newIter := make(chan T)
	go func() {
		defer close(newIter)
		m := make(map[K]bool, 0)
		for v := range it {
			cmp := f(v)
			if !m[cmp] {
				m[cmp] = true
				newIter <- v
			}
		}
	}()
	return newIter
}

// Join combines all elements into a string separated by sep.
func (it Iterator[T]) Join(sep string) string {
	out := ""
	v, ok := <-it
	if !ok {
		return out
	}
	out += fmt.Sprintf("%v", v)
	for {
		v, ok := <-it
		if !ok {
			return out
		}
		out += fmt.Sprintf("%s%v", sep, v)
	}
}
