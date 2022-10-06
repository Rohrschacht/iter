package iter

import "fmt"

func Example() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(i int) bool { return i%2 == 0 }).
		Map(func(i int) int { return i * i }).
		Collect()
	fmt.Println(it)
	// output:
	// [4 16 36]
}

func ExampleFromChan() {
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
	// output:
	// [1 2 3 4]
}

func ExampleFromSlice() {
	s := []int{1, 2, 3, 4, 5}
	it := FromSlice(s)
	fmt.Println(it.Collect())

	s2 := []string{"foo", "bar"}
	it2 := FromSlice(s2)
	fmt.Println(it2.Collect())
	// output:
	// [1 2 3 4 5]
	// [foo bar]
}

func ExampleFromMap() {
	m := map[int]string{1: "1", 2: "2", 3: "3"}
	it := FromMap(m)
	fmt.Println(it.Collect())
}

func ExampleFromMapKeys() {
	m := map[int]string{1: "1", 2: "2", 3: "3"}
	it := FromMapKeys(m)
	fmt.Println(it.Collect())
}

func ExampleFromMapValues() {
	m := map[int]string{1: "1", 2: "2", 3: "3"}
	it := FromMapValues(m)
	fmt.Println(it.Collect())
}

func ExampleIterator_Filter() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	filteredIter := it.Filter(func(i int) bool { return i%2 == 0 })
	fmt.Println(filteredIter.Collect())
	// output:
	// [2 4 6]
}

func ExampleIterator_Map() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	squaredIter := it.Map(func(i int) int { return i * i })
	fmt.Println(squaredIter.Collect())
	// output:
	// [1 4 9 16 25 36]
}

func ExampleMapInto() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	mappedIter := MapInto(it, func(i int) string { return fmt.Sprintf("%d", i) })
	fmt.Println(mappedIter.Collect())
	// output:
	// [1 2 3 4 5 6]
}

func ExampleIterator_Skip() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
		Skip(3)
	fmt.Println(it.Collect())
	// output:
	// [4 5 6]
}

func ExampleIterator_Nth() {
	n := FromSlice([]int{1, 2, 3, 4, 5, 6}).Nth(3)
	fmt.Println(*n)
	// output:
	// 3
}

func ExampleIterator_Count() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	fmt.Println(it.Count())
	// output:
	// 6
}

func ExampleIterator_Last() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	fmt.Println(it.Last())
	// output:
	// 6
}

func ExampleIterator_StepBy() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
		StepBy(2)
	fmt.Println(it.Collect())
	// output:
	// [1 3 5]
}

func ExampleIterator_Chain() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	other := FromSlice([]int{7, 8, 9})
	result := it.Chain(other)
	fmt.Println(result.Collect())
	// output:
	// [1 2 3 4 5 6 7 8 9]
}

func ExampleIterator_Intersperse() {
	it := FromSlice([]int{1, 2, 3}).
		Intersperse(5)
	fmt.Println(it.Collect())
	// output:
	// [1 5 2 5 3]
}

func ExampleIterator_ForEach() {
	it := FromSlice([]int{1, 2, 3})
	it.ForEach(func(x int) { fmt.Println(x) })
	// output:
	// 1
	// 2
	// 3
}

func ExampleIterator_SkipWhile() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
		SkipWhile(func(n int) bool { return n < 4 })
	fmt.Println(it.Collect())
	// output:
	// [4 5 6]
}

func ExampleIterator_TakeWhile() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
		TakeWhile(func(n int) bool { return n < 4 })
	fmt.Println(it.Collect())
	// output:
	// [1 2 3]
}

func ExampleIterator_Take() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
		Take(3)
	fmt.Println(it.Collect())
	// output:
	// [1 2 3]
}

func ExampleIterator_Inspect() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(x int) bool { return x%2 == 0 }).
		Inspect(func(x int) { fmt.Printf("got through filter: %d\n", x) }).
		Map(func(x int) int { return x * x })
	fmt.Println(it.Collect())
	// output:
	// got through filter: 2
	// got through filter: 4
	// got through filter: 6
	// [4 16 36]
}

func ExampleIterator_Partition() {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	even, odd := it.Partition(func(x int) bool { return x%2 == 0 })
	fmt.Println(even)
	fmt.Println(odd)
	// output:
	// [2 4 6]
	// [1 3 5]
}

func ExampleIterator_Fold() {
	it := FromSlice([]int{1, 2, 3})
	n := it.Fold(0, func(acc, x int) int { return acc + x })
	fmt.Println(n)
	// output:
	// 6
}

func ExampleIterator_Reduce() {
	it := FromSlice([]int{1, 2, 3})
	n := it.Reduce(func(acc, x int) int { return acc + x })
	fmt.Println(*n)
	// output:
	// 6
}

func ExampleIterator_All() {
	it := FromSlice([]int{1, 2, 3})
	b := it.All(func(x int) bool { return x < 100 })
	fmt.Println(b)
	// output:
	// true
}

func ExampleIterator_Any() {
	it := FromSlice([]int{1, 2, 3})
	b := it.Any(func(x int) bool { return x%2 == 0 })
	fmt.Println(b)
	// output:
	// true
}

func ExampleIterator_Find() {
	it := FromSlice([]int{3, 4, 5})
	x := it.Find(func(x int) bool { return x%2 == 0 })
	fmt.Println(*x)
	// output:
	// 4
}

func ExampleIterator_Position() {
	it := FromSlice([]int{3, 4, 5})
	x := it.Position(func(x int) bool { return x%2 == 0 })
	fmt.Println(*x)
	// output:
	// 2
}

func ExampleIterator_Collect() {
	it := FromSlice([]int{3, 4, 5})
	fmt.Println(it.Collect())
	// output:
	// [3 4 5]
}

func ExampleZip() {
	it1 := FromSlice([]int{1, 2, 3})
	it2 := FromSlice([]int{4, 5, 6})
	it3 := Zip(it1, it2)
	fmt.Println(it3.Collect())
	// output:
	// [{1 4} {2 5} {3 6}]
}

func ExampleIterator_Interleave() {
	it1 := FromSlice([]int{1, 2, 3})
	it2 := FromSlice([]int{4, 5, 6, 7, 8})
	interleaved := it1.Interleave(it2)
	fmt.Println(interleaved.Collect())
	// output:
	// [1 4 2 5 3 6 7 8]
}

func ExampleIterator_InterleaveShortest() {
	it1 := FromSlice([]int{1, 2, 3})
	it2 := FromSlice([]int{4, 5, 6, 7, 8})
	interleaved := it1.InterleaveShortest(it2)
	fmt.Println(interleaved.Collect())
	// output:
	// [1 4 2 5 3 6]
}

func ExampleIterator_GroupBy() {
	it := FromSlice([]int{-2, -1, 1, 2, 3, -4, -5, 7, 8})
	grouped := it.GroupBy(func(x int) bool { return x > 0 })
	fmt.Println(grouped)
	// output:
	// [[-2 -1] [1 2 3] [-4 -5] [7 8]]
}

func ExampleIterator_Chunks() {
	it := FromSlice([]int{-2, -1, 1, 2, 3, -4, -5, 7, 8})
	chunks := it.Chunks(3)
	fmt.Println(chunks)
	// output:
	// [[-2 -1 1] [2 3 -4] [-5 7 8]]
}

func ExampleIterator_Windows() {
	it := FromSlice([]int{1, 2, 3, 4, 5})
	windows := it.Windows(2)
	fmt.Println(windows)
	// output:
	// [[1 2] [2 3] [3 4] [4 5]]
}

func ExampleCartesianProduct() {
	it1 := FromSlice([]int{1, 2, 3})
	it2 := FromSlice([]int{4, 5, 6, 7, 8})
	cp := CartesianProduct(it1, it2)
	fmt.Println(cp.Collect())
	// output:
	// [{1 4} {1 5} {1 6} {1 7} {1 8} {2 4} {2 5} {2 6} {2 7} {2 8} {3 4} {3 5} {3 6} {3 7} {3 8}]
}

func ExampleIterator_Dedup() {
	it := FromSlice([]int{1, 1, 1, 2, 2, 3, 4, 5, 6, 6}).
		Dedup(func(x, y int) bool { return x == y })
	fmt.Println(it.Collect())
	// output:
	// [1 2 3 4 5 6]
}

func ExampleUnique() {
	it := FromSlice([]int{1, 1, 2, 2, 1, 3, 4, 5, 6, 1, 6})
	uniq := Unique(it, func(x int) int { return x }) // int is already comparable
	fmt.Println(uniq.Collect())
	// output:
	// [1 2 3 4 5 6]
}

func ExampleIterator_Join() {
	it := FromSlice([]int{1, 2, 3, 4})
	fmt.Println(it.Join(","))
	// output:
	// 1,2,3,4
}
