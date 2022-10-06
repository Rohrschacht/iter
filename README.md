<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# iter

```go
import "github.com/rohrschacht/iter"
```

Package iter implements a generic iterator type.

## About

This package implements Rust\-inspired iterators using Go 1.18 generics. Internally, the iterators are implemented using Goroutines and channels. This, using the provided methods on the iterators, one can define a pipeline that automatically uses multiple threads.

## Examples

```
it := iter.FromSlice([]int{1, 2, 3, 4, 5, 6}).
	Filter(func(i int) bool { return i%2 == 0 }).
	Map(func(i int) int { return i * i }).
	Collect()
expected := []int{4, 16, 36}
```

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
    Filter(func(i int) bool { return i%2 == 0 }).
    Map(func(i int) int { return i * i }).
    Collect()
fmt.Println(it)
```

#### Output

```
[4 16 36]
```

</p>
</details>

## Index

- [type Iterator](<#type-iterator>)
  - [func CartesianProduct[T, K any](it Iterator[T], other Iterator[K]) Iterator[Pair[T, K]]](<#func-cartesianproduct>)
  - [func FromChan[T any](c chan T) Iterator[T]](<#func-fromchan>)
  - [func FromMap[T comparable, K any](m map[T]K) Iterator[Pair[T, K]]](<#func-frommap>)
  - [func FromMapKeys[T comparable, K any](m map[T]K) Iterator[T]](<#func-frommapkeys>)
  - [func FromMapValues[K comparable, T any](m map[K]T) Iterator[T]](<#func-frommapvalues>)
  - [func FromSlice[T any](slice []T) Iterator[T]](<#func-fromslice>)
  - [func MapInto[T, K any](it Iterator[T], f func(T) K) Iterator[K]](<#func-mapinto>)
  - [func Unique[T any, K comparable](it Iterator[T], f func(T) K) Iterator[T]](<#func-unique>)
  - [func Zip[T, K any](it Iterator[T], other Iterator[K]) Iterator[Pair[T, K]]](<#func-zip>)
  - [func (it Iterator[T]) All(f func(T) bool) bool](<#func-iteratort-all>)
  - [func (it Iterator[T]) Any(f func(T) bool) bool](<#func-iteratort-any>)
  - [func (it Iterator[T]) Chain(other Iterator[T]) Iterator[T]](<#func-iteratort-chain>)
  - [func (it Iterator[T]) Chunks(n uint) [][]T](<#func-iteratort-chunks>)
  - [func (it Iterator[T]) Collect() []T](<#func-iteratort-collect>)
  - [func (it Iterator[T]) Count() uint](<#func-iteratort-count>)
  - [func (it Iterator[T]) Dedup(f func(T, T) bool) Iterator[T]](<#func-iteratort-dedup>)
  - [func (it Iterator[T]) Filter(f func(T) bool) Iterator[T]](<#func-iteratort-filter>)
  - [func (it Iterator[T]) Find(f func(T) bool) *T](<#func-iteratort-find>)
  - [func (it Iterator[T]) Fold(acc T, f func(T, T) T) T](<#func-iteratort-fold>)
  - [func (it Iterator[T]) ForEach(f func(T))](<#func-iteratort-foreach>)
  - [func (it Iterator[T]) GroupBy(f func(T) bool) [][]T](<#func-iteratort-groupby>)
  - [func (it Iterator[T]) Inspect(f func(T)) Iterator[T]](<#func-iteratort-inspect>)
  - [func (it Iterator[T]) Interleave(other Iterator[T]) Iterator[T]](<#func-iteratort-interleave>)
  - [func (it Iterator[T]) InterleaveShortest(other Iterator[T]) Iterator[T]](<#func-iteratort-interleaveshortest>)
  - [func (it Iterator[T]) Intersperse(sep T) Iterator[T]](<#func-iteratort-intersperse>)
  - [func (it Iterator[T]) Join(sep string) string](<#func-iteratort-join>)
  - [func (it Iterator[T]) Last() T](<#func-iteratort-last>)
  - [func (it Iterator[T]) Map(f func(T) T) Iterator[T]](<#func-iteratort-map>)
  - [func (it Iterator[T]) Nth(n uint) *T](<#func-iteratort-nth>)
  - [func (it Iterator[T]) Partition(f func(T) bool) ([]T, []T)](<#func-iteratort-partition>)
  - [func (it Iterator[T]) Position(f func(T) bool) *uint](<#func-iteratort-position>)
  - [func (it Iterator[T]) Reduce(f func(T, T) T) *T](<#func-iteratort-reduce>)
  - [func (it Iterator[T]) Skip(n uint) Iterator[T]](<#func-iteratort-skip>)
  - [func (it Iterator[T]) SkipWhile(f func(T) bool) Iterator[T]](<#func-iteratort-skipwhile>)
  - [func (it Iterator[T]) StepBy(n uint) Iterator[T]](<#func-iteratort-stepby>)
  - [func (it Iterator[T]) Take(n uint) Iterator[T]](<#func-iteratort-take>)
  - [func (it Iterator[T]) TakeWhile(f func(T) bool) Iterator[T]](<#func-iteratort-takewhile>)
  - [func (it Iterator[T]) Windows(n uint) [][]T](<#func-iteratort-windows>)
- [type Pair](<#type-pair>)


## type Iterator

Iterator can be used to process data in a pipeline pattern.

```go
type Iterator[T any] chan T
```

### func CartesianProduct

```go
func CartesianProduct[T, K any](it Iterator[T], other Iterator[K]) Iterator[Pair[T, K]]
```

CartesianProduct returns an Iterator over the cartesian product of both given Iterators.

<details><summary>Example</summary>
<p>

```go
it1 := FromSlice([]int{1, 2, 3})
it2 := FromSlice([]int{4, 5, 6, 7, 8})
cp := CartesianProduct(it1, it2)
fmt.Println(cp.Collect())
```

#### Output

```
[{1 4} {1 5} {1 6} {1 7} {1 8} {2 4} {2 5} {2 6} {2 7} {2 8} {3 4} {3 5} {3 6} {3 7} {3 8}]
```

</p>
</details>

### func FromChan

```go
func FromChan[T any](c chan T) Iterator[T]
```

FromChan creates an Iterator from a channel.

<details><summary>Example</summary>
<p>

```go
c := make(chan int)
go func() {
    defer close(c)
    c <- 1
    c <- 2
    c <- 3
    c <- 4
}()
it := FromChan(c)
s := it.Collect()
fmt.Println(s)
```

#### Output

```
[1 2 3 4]
```

</p>
</details>

### func FromMap

```go
func FromMap[T comparable, K any](m map[T]K) Iterator[Pair[T, K]]
```

FromMap creates an Iterator of Pairs that contain key and value of the given map.

<details><summary>Example</summary>
<p>

```go
m := map[int]string{1: "1", 2: "2", 3: "3"}
it := FromMap(m)
fmt.Println(it.Collect())
```

</p>
</details>

### func FromMapKeys

```go
func FromMapKeys[T comparable, K any](m map[T]K) Iterator[T]
```

FromMapKeys creates an Iterator over the keys of the given map.

<details><summary>Example</summary>
<p>

```go
m := map[int]string{1: "1", 2: "2", 3: "3"}
it := FromMapKeys(m)
fmt.Println(it.Collect())
```

</p>
</details>

### func FromMapValues

```go
func FromMapValues[K comparable, T any](m map[K]T) Iterator[T]
```

FromMapValues creates an Iterator over the values of the given map.

<details><summary>Example</summary>
<p>

```go
m := map[int]string{1: "1", 2: "2", 3: "3"}
it := FromMapValues(m)
fmt.Println(it.Collect())
```

</p>
</details>

### func FromSlice

```go
func FromSlice[T any](slice []T) Iterator[T]
```

FromSlice creates an Iterator over the given slice.

<details><summary>Example</summary>
<p>

```go
s := []int{1, 2, 3, 4, 5}
it := FromSlice(s)
fmt.Println(it.Collect())

s2 := []string{"foo", "bar"}
it2 := FromSlice(s2)
fmt.Println(it2.Collect())
```

#### Output

```
[1 2 3 4 5]
[foo bar]
```

</p>
</details>

### func MapInto

```go
func MapInto[T, K any](it Iterator[T], f func(T) K) Iterator[K]
```

MapInto applies the given function to all elements and allows for the type to change.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6})
mappedIter := MapInto(it, func(i int) string { return fmt.Sprintf("%d", i) })
fmt.Println(mappedIter.Collect())
```

#### Output

```
[1 2 3 4 5 6]
```

</p>
</details>

### func Unique

```go
func Unique[T any, K comparable](it Iterator[T], f func(T) K) Iterator[T]
```

Unique produces an Iterator that returns unique elements from the given Iterator determined by the given condition.

Since Iterator can take any type, f has to convert the element type into a comparable type. If your type is already comparable, it is enough to just return it in the closure. See the example.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 1, 2, 2, 1, 3, 4, 5, 6, 1, 6})
uniq := Unique(it, func(x int) int { return x })
fmt.Println(uniq.Collect())
```

#### Output

```
[1 2 3 4 5 6]
```

</p>
</details>

### func Zip

```go
func Zip[T, K any](it Iterator[T], other Iterator[K]) Iterator[Pair[T, K]]
```

Zip creates a new Iterator that contains Pairs containing the elements of both Iterators.

If one of the input Iterators is shorter than the other one, the new Iterator will stop at that point.

<details><summary>Example</summary>
<p>

```go
it1 := FromSlice([]int{1, 2, 3})
it2 := FromSlice([]int{4, 5, 6})
it3 := Zip(it1, it2)
fmt.Println(it3.Collect())
```

#### Output

```
[{1 4} {2 5} {3 6}]
```

</p>
</details>

### func \(Iterator\[T\]\) All

```go
func (it Iterator[T]) All(f func(T) bool) bool
```

All checks whether the given condition is true for all elements.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3})
b := it.All(func(x int) bool { return x < 100 })
fmt.Println(b)
```

#### Output

```
true
```

</p>
</details>

### func \(Iterator\[T\]\) Any

```go
func (it Iterator[T]) Any(f func(T) bool) bool
```

Any checks whether there exists one element for which the given condition is true.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3})
b := it.Any(func(x int) bool { return x%2 == 0 })
fmt.Println(b)
```

#### Output

```
true
```

</p>
</details>

### func \(Iterator\[T\]\) Chain

```go
func (it Iterator[T]) Chain(other Iterator[T]) Iterator[T]
```

Chain creates a new Iterator which returns the elements of both Iterators.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6})
other := FromSlice([]int{7, 8, 9})
result := it.Chain(other)
fmt.Println(result.Collect())
```

#### Output

```
[1 2 3 4 5 6 7 8 9]
```

</p>
</details>

### func \(Iterator\[T\]\) Chunks

```go
func (it Iterator[T]) Chunks(n uint) [][]T
```

Chunks returns a list of slices containing at most n elements of the original Iterator.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{-2, -1, 1, 2, 3, -4, -5, 7, 8})
chunks := it.Chunks(3)
fmt.Println(chunks)
```

#### Output

```
[[-2 -1 1] [2 3 -4] [-5 7 8]]
```

</p>
</details>

### func \(Iterator\[T\]\) Collect

```go
func (it Iterator[T]) Collect() []T
```

Collect consumes the Iterator, returning a slice of all its elements.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{3, 4, 5})
fmt.Println(it.Collect())
```

#### Output

```
[3 4 5]
```

</p>
</details>

### func \(Iterator\[T\]\) Count

```go
func (it Iterator[T]) Count() uint
```

Count consumes the Iterator and returns its number of elements.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6})
fmt.Println(it.Count())
```

#### Output

```
6
```

</p>
</details>

### func \(Iterator\[T\]\) Dedup

```go
func (it Iterator[T]) Dedup(f func(T, T) bool) Iterator[T]
```

Dedup removes duplicates from sections of consecutive elements determined by the given condition.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 1, 1, 2, 2, 3, 4, 5, 6, 6}).
    Dedup(func(x, y int) bool { return x == y })
fmt.Println(it.Collect())
```

#### Output

```
[1 2 3 4 5 6]
```

</p>
</details>

### func \(Iterator\[T\]\) Filter

```go
func (it Iterator[T]) Filter(f func(T) bool) Iterator[T]
```

Filter uses the given function to determine whether elements should continue through the pipeline.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6})
filteredIter := it.Filter(func(i int) bool { return i%2 == 0 })
fmt.Println(filteredIter.Collect())
```

#### Output

```
[2 4 6]
```

</p>
</details>

### func \(Iterator\[T\]\) Find

```go
func (it Iterator[T]) Find(f func(T) bool) *T
```

Find returns a pointer to the first element for which the given condition is true.

If no such element exists, nil is returned.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{3, 4, 5})
x := it.Find(func(x int) bool { return x%2 == 0 })
fmt.Println(*x)
```

#### Output

```
4
```

</p>
</details>

### func \(Iterator\[T\]\) Fold

```go
func (it Iterator[T]) Fold(acc T, f func(T, T) T) T
```

Fold applies the given function to all elements, folding them into the given accumulator.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3})
n := it.Fold(0, func(acc, x int) int { return acc + x })
fmt.Println(n)
```

#### Output

```
6
```

</p>
</details>

### func \(Iterator\[T\]\) ForEach

```go
func (it Iterator[T]) ForEach(f func(T))
```

ForEach executes the given function for each element of the Iterator.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3})
it.ForEach(func(x int) { fmt.Println(x) })
```

#### Output

```
1
2
3
```

</p>
</details>

### func \(Iterator\[T\]\) GroupBy

```go
func (it Iterator[T]) GroupBy(f func(T) bool) [][]T
```

GroupBy returns a list of slices, which elements are grouped by the given condition.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{-2, -1, 1, 2, 3, -4, -5, 7, 8})
grouped := it.GroupBy(func(x int) bool { return x > 0 })
fmt.Println(grouped)
```

#### Output

```
[[-2 -1] [1 2 3] [-4 -5] [7 8]]
```

</p>
</details>

### func \(Iterator\[T\]\) Inspect

```go
func (it Iterator[T]) Inspect(f func(T)) Iterator[T]
```

Inspect applies the given function on each element while the Iterator is consumed.

This is helpful for debugging, see the example.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
    Filter(func(x int) bool { return x%2 == 0 }).
    Inspect(func(x int) { fmt.Printf("got through filter: %d\n", x) }).
    Map(func(x int) int { return x * x })
fmt.Println(it.Collect())
```

#### Output

```
got through filter: 2
got through filter: 4
got through filter: 6
[4 16 36]
```

</p>
</details>

### func \(Iterator\[T\]\) Interleave

```go
func (it Iterator[T]) Interleave(other Iterator[T]) Iterator[T]
```

Interleave creates a new Iterator that alternates between the two given Iterators.

<details><summary>Example</summary>
<p>

```go
it1 := FromSlice([]int{1, 2, 3})
it2 := FromSlice([]int{4, 5, 6, 7, 8})
interleaved := it1.Interleave(it2)
fmt.Println(interleaved.Collect())
```

#### Output

```
[1 4 2 5 3 6 7 8]
```

</p>
</details>

### func \(Iterator\[T\]\) InterleaveShortest

```go
func (it Iterator[T]) InterleaveShortest(other Iterator[T]) Iterator[T]
```

InterleaveShortest creates a new Iterator that alternates between the two given Iterators until at least one of them runs out.

<details><summary>Example</summary>
<p>

```go
it1 := FromSlice([]int{1, 2, 3})
it2 := FromSlice([]int{4, 5, 6, 7, 8})
interleaved := it1.InterleaveShortest(it2)
fmt.Println(interleaved.Collect())
```

#### Output

```
[1 4 2 5 3 6]
```

</p>
</details>

### func \(Iterator\[T\]\) Intersperse

```go
func (it Iterator[T]) Intersperse(sep T) Iterator[T]
```

Intersperse inserts the separator sep between each element of the Iterator.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3}).
    Intersperse(5)
fmt.Println(it.Collect())
```

#### Output

```
[1 5 2 5 3]
```

</p>
</details>

### func \(Iterator\[T\]\) Join

```go
func (it Iterator[T]) Join(sep string) string
```

Join combines all elements into a string separated by sep.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4})
fmt.Println(it.Join(","))
```

#### Output

```
1,2,3,4
```

</p>
</details>

### func \(Iterator\[T\]\) Last

```go
func (it Iterator[T]) Last() T
```

Last returns the last element of the Iterator, consuming it in the process.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6})
fmt.Println(it.Last())
```

#### Output

```
6
```

</p>
</details>

### func \(Iterator\[T\]\) Map

```go
func (it Iterator[T]) Map(f func(T) T) Iterator[T]
```

Map applies the given function to all elements going through the pipeline.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6})
squaredIter := it.Map(func(i int) int { return i * i })
fmt.Println(squaredIter.Collect())
```

#### Output

```
[1 4 9 16 25 36]
```

</p>
</details>

### func \(Iterator\[T\]\) Nth

```go
func (it Iterator[T]) Nth(n uint) *T
```

Nth returns a pointer to the element at position n.

If there are fewer than n elements in the Iterator, nil is returned.

<details><summary>Example</summary>
<p>

```go
n := FromSlice([]int{1, 2, 3, 4, 5, 6}).Nth(3)
fmt.Println(*n)
```

#### Output

```
3
```

</p>
</details>

### func \(Iterator\[T\]\) Partition

```go
func (it Iterator[T]) Partition(f func(T) bool) ([]T, []T)
```

Partition splits the contents of the iterator based on the condition defined in the given function.

Two slices are returned. The first slice contains all elements of the Iterator for which f evaluated to true. The second slice contains all elements for which f evaluated to false.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6})
even, odd := it.Partition(func(x int) bool { return x%2 == 0 })
fmt.Println(even)
fmt.Println(odd)
```

#### Output

```
[2 4 6]
[1 3 5]
```

</p>
</details>

### func \(Iterator\[T\]\) Position

```go
func (it Iterator[T]) Position(f func(T) bool) *uint
```

Position returns the position of the first element for which the given condition is true as a pointer.

If no such element exists, nil is returned.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{3, 4, 5})
x := it.Position(func(x int) bool { return x%2 == 0 })
fmt.Println(*x)
```

#### Output

```
2
```

</p>
</details>

### func \(Iterator\[T\]\) Reduce

```go
func (it Iterator[T]) Reduce(f func(T, T) T) *T
```

Reduce folds the Iterator using the given function, using the first element as the initial accumulator.

Reduce returns a pointer for the accumulated value. If the Iterator is empty, this will be nil.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3})
n := it.Reduce(func(acc, x int) int { return acc + x })
fmt.Println(*n)
```

#### Output

```
6
```

</p>
</details>

### func \(Iterator\[T\]\) Skip

```go
func (it Iterator[T]) Skip(n uint) Iterator[T]
```

Skip skips the first n elements of the Iterator.

n can be larger than the number of elements in the Iterator, which will empty it.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
    Skip(3)
fmt.Println(it.Collect())
```

#### Output

```
[4 5 6]
```

</p>
</details>

### func \(Iterator\[T\]\) SkipWhile

```go
func (it Iterator[T]) SkipWhile(f func(T) bool) Iterator[T]
```

SkipWhile discards all elements until the condition of the given function is met once.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
    SkipWhile(func(n int) bool { return n < 4 })
fmt.Println(it.Collect())
```

#### Output

```
[4 5 6]
```

</p>
</details>

### func \(Iterator\[T\]\) StepBy

```go
func (it Iterator[T]) StepBy(n uint) Iterator[T]
```

StepBy advances the Iterator by n elements every time something is taken.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
    StepBy(2)
fmt.Println(it.Collect())
```

#### Output

```
[1 3 5]
```

</p>
</details>

### func \(Iterator\[T\]\) Take

```go
func (it Iterator[T]) Take(n uint) Iterator[T]
```

Take takes the first n elements of the Iterator.

All elements after the first n elements will be discarded.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
    Take(3)
fmt.Println(it.Collect())
```

#### Output

```
[1 2 3]
```

</p>
</details>

### func \(Iterator\[T\]\) TakeWhile

```go
func (it Iterator[T]) TakeWhile(f func(T) bool) Iterator[T]
```

TakeWhile takes elements until the condition of the given function is false once.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
    TakeWhile(func(n int) bool { return n < 4 })
fmt.Println(it.Collect())
```

#### Output

```
[1 2 3]
```

</p>
</details>

### func \(Iterator\[T\]\) Windows

```go
func (it Iterator[T]) Windows(n uint) [][]T
```

Windows returns all overlapping subslices of length n of the original Iterator.

<details><summary>Example</summary>
<p>

```go
it := FromSlice([]int{1, 2, 3, 4, 5})
windows := it.Windows(2)
fmt.Println(windows)
```

#### Output

```
[[1 2] [2 3] [3 4] [4 5]]
```

</p>
</details>

## type Pair

Pair is used as a helper when an Iterator has to hold multiple values.

```go
type Pair[T, K any] struct {
    X   T
    Y   K
}
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
