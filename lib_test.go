package iter

import (
	"fmt"
	"testing"
)

func TestFromChan(t *testing.T) {
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
	expected := []int{1, 2, 3, 4}
	if len(s) != len(expected) {
		t.Error("FromChan did not work for ints")
		return
	}
	for k := 0; k < len(s); k++ {
		if s[k] != expected[k] {
			t.Error("FromChan did not work for ints")
			return
		}
	}
}

func TestFromSlice(t *testing.T) {
	str := []string{"this", "is", "a", "test"}
	strIter := FromSlice(str)
	strSlice := strIter.Collect()
	if len(str) != len(strSlice) {
		t.Error("FromSlice did not work for strings")
		return
	}
	for i := 0; i < len(str); i++ {
		if strSlice[i] != str[i] {
			t.Error("FromSlice did not work for strings")
			return
		}
	}

	i := []int{1, 2, 3, 4, 5}
	iIter := FromSlice(i)
	iSlice := iIter.Collect()
	if len(i) != len(iSlice) {
		t.Error("FromSlice did not work for ints")
		return
	}
	for k := 0; k < len(i); k++ {
		if iSlice[k] != i[k] {
			t.Error("FromSlice did not work for ints")
			return
		}
	}
}

func TestFromMap(t *testing.T) {
	m := map[int]string{1: "1", 2: "2", 3: "3"}
	it := FromMap(m)
	if !it.All(func(pair Pair[int, string]) bool { return fmt.Sprintf("%d", pair.X) == pair.Y }) {
		t.Error("FromMap did not work")
	}
}

func TestFromMapKeys(t *testing.T) {
	m := map[int]string{1: "1", 2: "2", 3: "3"}
	it := FromMapKeys(m).Collect()
	expected := []int{1, 2, 3}
	if len(it) != len(expected) {
		t.Error("FromMapKeys did not work")
		return
	}
	for k := 0; k < len(it); k++ {
		found := false
		for _, v := range expected {
			if it[k] == v {
				found = true
			}
		}
		if !found {
			t.Error("FromMapKeys did not work")
			return
		}
	}
}

func TestFromMapValues(t *testing.T) {
	m := map[int]string{1: "1", 2: "2", 3: "3"}
	it := FromMapValues(m).Collect()
	expected := []string{"1", "2", "3"}
	if len(it) != len(expected) {
		t.Error("FromMapValues did not work")
		return
	}
	for k := 0; k < len(it); k++ {
		found := false
		for _, v := range expected {
			if it[k] == v {
				found = true
			}
		}
		if !found {
			t.Error("FromMapValues did not work")
			return
		}
	}
}

func TestIterator_Filter(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	expected := []int{2, 4, 6}
	filteredIter := it.Filter(func(i int) bool { return i%2 == 0 })
	iSlice := filteredIter.Collect()
	if len(expected) != len(iSlice) {
		t.Error("Filter did not work for ints")
	}
	for i := 0; i < len(expected); i++ {
		if iSlice[i] != expected[i] {
			t.Error("Filter did not work for ints")
		}
	}
}

func TestIterator_Map(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	expected := []int{1, 4, 9, 16, 25, 36}
	squaredIter := it.Map(func(i int) int { return i * i })
	iSlice := squaredIter.Collect()
	if len(expected) != len(iSlice) {
		t.Error("Map did not work for ints")
	}
	for i := 0; i < len(expected); i++ {
		if iSlice[i] != expected[i] {
			t.Error("Map did not work for ints")
		}
	}
}

func TestIterator_MapOtherType(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	expected := []string{"1", "2", "3", "4", "5", "6"}
	mappedIter := MapInto(it, func(i int) string { return fmt.Sprintf("%d", i) }).Collect()
	if len(expected) != len(mappedIter) {
		t.Error("Map did not work for ints")
	}
	for i := 0; i < len(expected); i++ {
		if mappedIter[i] != expected[i] {
			t.Error("Map did not work for ints")
		}
	}
}

func TestIterator_Chaining(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(i int) bool { return i%2 == 0 }).
		Map(func(i int) int { return i * i }).
		Collect()
	expected := []int{4, 16, 36}
	if len(expected) != len(it) {
		t.Error("Chaining did not work for ints")
		t.Errorf("iter: %v\nexpected: %v\n", it, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if it[i] != expected[i] {
			t.Error("Chaining did not work for ints")
			t.Errorf("iter: %v\nexpected: %v\n", it, expected)
			return
		}
	}
}

func TestIterator_Skip(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).Skip(3).Collect()
	expected := []int{4, 5, 6}
	if len(expected) != len(it) {
		t.Error("Skip did not work for ints")
		t.Errorf("iter: %v\nexpected: %v\n", it, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if it[i] != expected[i] {
			t.Error("Skip did not work for ints")
			t.Errorf("iter: %v\nexpected: %v\n", it, expected)
			return
		}
	}
}

func TestIterator_SkipStr(t *testing.T) {
	iter := FromSlice([]string{"ab", "cd", "ef", "gh"}).Skip(2).Collect()
	expected := []string{"ef", "gh"}
	if len(expected) != len(iter) {
		t.Error("Skip did not work for strings")
		t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if iter[i] != expected[i] {
			t.Error("Skip did not work for strings")
			t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
			return
		}
	}
}

func TestIterator_SkipMuch(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6}).Skip(10).Collect()
	expected := []int{}
	if len(expected) != len(iter) {
		t.Error("Skip did not work for ints")
		t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if iter[i] != expected[i] {
			t.Error("Skip did not work for ints")
			t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
			return
		}
	}
}

func TestIterator_SkipStrMuch(t *testing.T) {
	iter := FromSlice([]string{"ab", "cd", "ef", "gh"}).Skip(10).Collect()
	expected := []string{}
	if len(expected) != len(iter) {
		t.Error("Skip did not work for strings")
		t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if iter[i] != expected[i] {
			t.Error("Skip did not work for strings")
			t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
			return
		}
	}
}

func TestIterator_Nth(t *testing.T) {
	n3 := FromSlice([]int{1, 2, 3, 4, 5, 6}).Nth(3)
	if *n3 != 3 {
		t.Errorf("Nth did not work\niter: %v\nexpected: %v\n", *n3, 3)
	}

	n6 := FromSlice([]int{1, 2, 3, 4, 5, 6}).Nth(6)
	if *n6 != 6 {
		t.Errorf("Nth did not work\niter: %v\nexpected: %v\n", *n6, 6)
	}

	n7 := FromSlice([]int{1, 2, 3, 4, 5, 6}).Nth(7)
	if n7 != nil {
		t.Errorf("Nth did not work\niter: %v\nexpected: %v\n", n7, nil)
	}
}

func TestIterator_Count(t *testing.T) {
	c := FromSlice([]int{1, 2, 3, 4, 5, 6}).Count()
	if c != 6 {
		t.Errorf("Count did not work\ncounted: %d\nexpected: %d", c, 6)
	}
}

func TestIterator_Last(t *testing.T) {
	l := FromSlice([]int{1, 2, 3, 4, 5, 6}).Last()
	if l != 6 {
		t.Errorf("Last did not work\ncounted: %d\nexpected: %d", l, 6)
	}
}

func TestIterator_StepBy(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).StepBy(2).Collect()
	expected := []int{1, 3, 5}
	if len(expected) != len(it) {
		t.Error("StepBy did not work for ints")
		t.Errorf("it: %v\nexpected: %v\n", it, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if it[i] != expected[i] {
			t.Error("StepBy did not work for ints")
			t.Errorf("it: %v\nexpected: %v\n", it, expected)
			return
		}
	}
}

func TestIterator_Chain(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	other := FromSlice([]int{7, 8, 9})
	result := it.Chain(other).Collect()
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if len(expected) != len(result) {
		t.Error("Chain did not work for ints")
		t.Errorf("it: %v\nexpected: %v\n", result, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {
			t.Error("Chain did not work for ints")
			t.Errorf("it: %v\nexpected: %v\n", result, expected)
			return
		}
	}
}

func TestIterator_Intersperse(t *testing.T) {
	it := FromSlice([]int{1, 2, 3}).Intersperse(5).Collect()
	expected := []int{1, 5, 2, 5, 3}
	if len(expected) != len(it) {
		t.Error("Intersperse did not work for ints")
		t.Errorf("it: %v\nexpected: %v\n", it, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if it[i] != expected[i] {
			t.Error("Intersperse did not work for ints")
			t.Errorf("it: %v\nexpected: %v\n", it, expected)
			return
		}
	}
}

func TestIterator_ForEach(t *testing.T) {
	n := 0
	it := FromSlice([]int{1, 2, 3})
	it.ForEach(func(x int) { n += x })
	if n != 6 {
		t.Errorf("ForEach did not work for ints\nit: %d\nexpected: %d\n", n, 6)
	}
}

func TestIterator_SkipWhile(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).SkipWhile(func(n int) bool { return n < 4 }).Collect()
	expected := []int{4, 5, 6}
	if len(expected) != len(it) {
		t.Error("SkipWhile did not work for ints")
		t.Errorf("it: %v\nexpected: %v\n", it, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if it[i] != expected[i] {
			t.Error("SkipWhile did not work for ints")
			t.Errorf("it: %v\nexpected: %v\n", it, expected)
			return
		}
	}
}

func TestIterator_TakeWhile(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).TakeWhile(func(n int) bool { return n < 4 }).Collect()
	expected := []int{1, 2, 3}
	if len(expected) != len(it) {
		t.Error("SkipWhile did not work for ints")
		t.Errorf("it: %v\nexpected: %v\n", it, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if it[i] != expected[i] {
			t.Error("SkipWhile did not work for ints")
			t.Errorf("it: %v\nexpected: %v\n", it, expected)
			return
		}
	}
}

func TestIterator_Take(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6}).Take(3).Collect()
	expected := []int{1, 2, 3}
	if len(expected) != len(it) {
		t.Error("Take did not work for ints")
		t.Errorf("it: %v\nexpected: %v\n", it, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if it[i] != expected[i] {
			t.Error("Take did not work for ints")
			t.Errorf("it: %v\nexpected: %v\n", it, expected)
			return
		}
	}
}

func TestIterator_Inspect(t *testing.T) {
	n := 0
	it := FromSlice([]int{1, 2, 3}).Inspect(func(x int) { n += x })
	iSlice := it.Collect()
	expected := []int{1, 2, 3}
	if n != 6 {
		t.Errorf("Inspect did not work for ints\nit: %d\nexpected: %d\n", n, 6)
	}
	if len(expected) != len(iSlice) {
		t.Error("Inspect did not work for ints")
	}
	for i := 0; i < len(expected); i++ {
		if iSlice[i] != expected[i] {
			t.Error("Inspect did not work for ints")
		}
	}
}

func TestIterator_Partition(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5, 6})
	even, odd := []int{2, 4, 6}, []int{1, 3, 5}
	left, right := it.Partition(func(x int) bool { return x%2 == 0 })
	if len(left) != len(even) {
		t.Error("Partition did not work for ints")
	}
	for i := 0; i < len(left); i++ {
		if left[i] != even[i] {
			t.Error("Partition did not work for ints")
		}
	}
	if len(right) != len(odd) {
		t.Error("Partition did not work for ints")
	}
	for i := 0; i < len(right); i++ {
		if right[i] != odd[i] {
			t.Error("Partition did not work for ints")
		}
	}
}

func TestIterator_Fold(t *testing.T) {
	it := FromSlice([]int{1, 2, 3})
	n := it.Fold(0, func(acc, x int) int { return acc + x })
	if n != 6 {
		t.Errorf("Fold did not work for ints\nit: %d\nexpected: %d\n", n, 6)
	}
}

func TestIterator_Reduce(t *testing.T) {
	it := FromSlice([]int{1, 2, 3})
	n := it.Reduce(func(acc, x int) int { return acc + x })
	if *n != 6 {
		t.Errorf("Reduce did not work for ints\nit: %d\nexpected: %d\n", *n, 6)
	}
}

func TestIterator_All(t *testing.T) {
	it := FromSlice([]int{1, 2, 3})
	b := it.All(func(x int) bool { return x < 100 })
	if !b {
		t.Error("All did not work for ints")
		return
	}

	it = FromSlice([]int{1, 2, 3})
	b = it.All(func(x int) bool { return x < 3 })
	if b {
		t.Error("All did not work for ints")
		return
	}
}

func TestIterator_Any(t *testing.T) {
	it := FromSlice([]int{1, 2, 3})
	b := it.Any(func(x int) bool { return x%2 == 0 })
	if !b {
		t.Error("Any did not work for ints")
		return
	}

	it = FromSlice([]int{1, 2, 3})
	b = it.Any(func(x int) bool { return x == 0 })
	if b {
		t.Error("Any did not work for ints")
		return
	}
}

func TestIterator_Find(t *testing.T) {
	it := FromSlice([]int{1, 2, 3})
	b := it.Find(func(x int) bool { return x%2 == 0 })
	if *b != 2 {
		t.Error("Find did not work for ints")
		return
	}

	it = FromSlice([]int{1, 2, 3})
	b = it.Find(func(x int) bool { return x == 0 })
	if b != nil {
		t.Error("Find did not work for ints")
		return
	}
}

func TestIterator_Position(t *testing.T) {
	it := FromSlice([]int{1, 2, 3})
	p := it.Position(func(x int) bool { return x%2 == 0 })
	if *p != 2 {
		t.Error("Position did not work for ints")
		return
	}

	it = FromSlice([]int{1, 2, 3})
	p = it.Position(func(x int) bool { return x == 0 })
	if p != nil {
		t.Error("Position did not work for ints")
		return
	}
}

func TestZip(t *testing.T) {
	it1 := FromSlice([]int{1, 2, 3})
	it2 := FromSlice([]int{4, 5, 6})
	it3 := Zip(it1, it2)
	if !it3.All(func(pair Pair[int, int]) bool { return pair.X+3 == pair.Y }) {
		t.Error("Zip did not work for ints")
	}
}

func TestIterator_Interleave(t *testing.T) {
	it1 := FromSlice([]int{1, 2, 3})
	it2 := FromSlice([]int{4, 5, 6, 7, 8})
	expected := []int{1, 4, 2, 5, 3, 6, 7, 8}
	interleaved := it1.Interleave(it2).Collect()
	if len(interleaved) != len(expected) {
		t.Error("Interleave did not work for ints")
		return
	}
	for i := 0; i < len(interleaved); i++ {
		if interleaved[i] != expected[i] {
			t.Error("Interleave did not work for ints")
			return
		}
	}
}

func TestIterator_InterleaveShortest(t *testing.T) {
	it1 := FromSlice([]int{1, 2, 3})
	it2 := FromSlice([]int{4, 5, 6, 7, 8})
	expected := []int{1, 4, 2, 5, 3, 6}
	interleaved := it1.InterleaveShortest(it2).Collect()
	if len(interleaved) != len(expected) {
		t.Error("InterleaveShortest did not work for ints")
		return
	}
	for i := 0; i < len(interleaved); i++ {
		if interleaved[i] != expected[i] {
			t.Error("InterleaveShortest did not work for ints")
			return
		}
	}
}

func TestIterator_GroupBy(t *testing.T) {
	it := FromSlice([]int{-2, -1, 1, 2, 3, -4, -5, 7, 8})
	expected := [][]int{{-2, -1}, {1, 2, 3}, {-4, -5}, {7, 8}}
	grouped := it.GroupBy(func(x int) bool { return x > 0 })
	if len(grouped) != len(expected) {
		t.Error("GroupBy did not work for ints")
		return
	}
	for i := 0; i < len(grouped); i++ {
		if len(grouped[i]) != len(expected[i]) {
			t.Error("GroupBy did not work for ints")
			return
		}
		for k := 0; k < len(grouped[i]); k++ {
			if grouped[i][k] != grouped[i][k] {
				t.Error("GroupBy did not work for ints")
				return
			}
		}
	}
}

func TestIterator_Chunks(t *testing.T) {
	it := FromSlice([]int{-2, -1, 1, 2, 3, -4, -5, 7, 8})
	expected := [][]int{{-2, -1, 1}, {2, 3, -4}, {-5, 7, 8}}
	chunked := it.Chunks(3)
	if len(chunked) != len(expected) {
		t.Error("Chunks did not work for ints")
		return
	}
	for i := 0; i < len(chunked); i++ {
		if len(chunked[i]) != len(expected[i]) {
			t.Error("Chunks did not work for ints")
			return
		}
		for k := 0; k < len(chunked[i]); k++ {
			if chunked[i][k] != chunked[i][k] {
				t.Error("Chunks did not work for ints")
				return
			}
		}
	}

	it = FromSlice([]int{-2, -1, 1, 2, 3, -4, -5, 7, 8})
	expected = [][]int{{-2, -1}, {1, 2}, {3, -4}, {-5, 7}, {8}}
	chunked = it.Chunks(2)
	if len(chunked) != len(expected) {
		t.Error("Chunks did not work for ints")
		return
	}
	for i := 0; i < len(chunked); i++ {
		if len(chunked[i]) != len(expected[i]) {
			t.Error("Chunks did not work for ints")
			return
		}
		for k := 0; k < len(chunked[i]); k++ {
			if chunked[i][k] != chunked[i][k] {
				t.Error("Chunks did not work for ints")
				return
			}
		}
	}
}

func TestIterator_Windows(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4, 5})
	expected := [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}
	windows := it.Windows(2)
	if len(windows) != len(expected) {
		t.Error("Windows did not work for ints")
		return
	}
	for i := 0; i < len(windows); i++ {
		if len(windows[i]) != len(expected[i]) {
			t.Error("Windows did not work for ints")
			return
		}
		for k := 0; k < len(windows[i]); k++ {
			if windows[i][k] != windows[i][k] {
				t.Error("Windows did not work for ints")
				return
			}
		}
	}

	it = FromSlice([]int{1, 2, 3, 4, 5})
	expected = [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}
	windows = it.Windows(3)
	if len(windows) != len(expected) {
		t.Error("Windows did not work for ints")
		return
	}
	for i := 0; i < len(windows); i++ {
		if len(windows[i]) != len(expected[i]) {
			t.Error("Windows did not work for ints")
			return
		}
		for k := 0; k < len(windows[i]); k++ {
			if windows[i][k] != windows[i][k] {
				t.Error("Windows did not work for ints")
				return
			}
		}
	}
}

func TestCartesianProduct(t *testing.T) {
	it1 := FromSlice([]int{1, 2, 3})
	it2 := FromSlice([]int{4, 5, 6, 7, 8})
	expected := []Pair[int, int]{
		{X: 1, Y: 4},
		{X: 1, Y: 5},
		{X: 1, Y: 6},
		{X: 1, Y: 7},
		{X: 1, Y: 8},
		{X: 2, Y: 4},
		{X: 2, Y: 5},
		{X: 2, Y: 6},
		{X: 2, Y: 7},
		{X: 2, Y: 8},
		{X: 3, Y: 4},
		{X: 3, Y: 5},
		{X: 3, Y: 6},
		{X: 3, Y: 7},
		{X: 3, Y: 8},
	}
	cp := CartesianProduct(it1, it2).Collect()
	if len(cp) != len(expected) {
		t.Error("CartesianProduct did not work for ints")
		return
	}
	for i := 0; i < len(cp); i++ {
		if cp[i] != expected[i] {
			t.Error("CartesianProduct did not work for ints")
			return
		}
	}
}

func TestIterator_Dedup(t *testing.T) {
	it := FromSlice([]int{1, 1, 1, 2, 2, 3, 4, 5, 6, 6})
	expected := []int{1, 2, 3, 4, 5, 6}
	deduped := it.Dedup(func(x, y int) bool { return x == y }).Collect()
	if len(deduped) != len(expected) {
		t.Error("Dedup did not work for ints")
		return
	}
	for i := 0; i < len(deduped); i++ {
		if deduped[i] != expected[i] {
			t.Error("Dedup did not work for ints")
			return
		}
	}
}

func TestUnique(t *testing.T) {
	it := FromSlice([]int{1, 1, 2, 2, 1, 3, 4, 5, 6, 1, 6})
	expected := []int{1, 2, 3, 4, 5, 6}
	uniq := Unique(it, func(x int) int { return x }).Collect()
	if len(uniq) != len(expected) {
		t.Error("Unique did not work for ints")
		return
	}
	for i := 0; i < len(uniq); i++ {
		if uniq[i] != expected[i] {
			t.Error("Unique did not work for ints")
			return
		}
	}
}

func TestIterator_Join(t *testing.T) {
	it := FromSlice([]int{1, 2, 3, 4})
	s := it.Join(",")
	if s != "1,2,3,4" {
		t.Error("Join did not work for ints")
		return
	}
}
