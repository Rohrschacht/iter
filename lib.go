package iter

type Iterator[T any] chan T

type Pair[T, K any] struct {
	x T
	y K
}

func FromChan[T any](c chan T) Iterator[T] {
	return c
}

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

func FromMap[T comparable, K any](m map[T]K) Iterator[Pair[T, K]] {
	it := make(chan Pair[T, K])
	go func() {
		defer close(it)
		for key, v := range m {
			it <- Pair[T, K]{x: key, y: v}
		}
	}()
	return it
}

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

func (it Iterator[T]) Collect() []T {
	var slice []T
	for v := range it {
		slice = append(slice, v)
	}
	return slice
}

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

func (it Iterator[T]) Skip(n uint) Iterator[T] {
	for i := uint(0); i < n; i++ {
		<-it
	}
	return it
}

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

func (it Iterator[T]) Count() uint {
	c := uint(0)
	for range it {
		c++
	}
	return c
}

func (it Iterator[T]) Last() T {
	var l T
	for v := range it {
		l = v
	}
	return l
}

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

func (it Iterator[T]) ForEach(f func(T)) {
	for v := range it {
		f(v)
	}
}

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
			newIter <- Pair[T, K]{x: v1, y: v2}
		}
	}()
	return newIter
}

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

func (it Iterator[T]) Fold(acc T, f func(T, T) T) T {
	for v := range it {
		acc = f(acc, v)
	}
	return acc
}

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

func (it Iterator[T]) All(f func(T) bool) bool {
	for v := range it {
		if !f(v) {
			return false
		}
	}
	return true
}

func (it Iterator[T]) Any(f func(T) bool) bool {
	for v := range it {
		if f(v) {
			return true
		}
	}
	return false
}

func (it Iterator[T]) Find(f func(T) bool) *T {
	for v := range it {
		if f(v) {
			return &v
		}
	}
	return nil
}

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
