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
	if !it.All(func(pair Pair[int, string]) bool { return fmt.Sprintf("%d", pair.x) == pair.y }) {
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
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6})
	expected := []int{2, 4, 6}
	filteredIter := iter.Filter(func(i int) bool { return i%2 == 0 })
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
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6})
	expected := []int{1, 4, 9, 16, 25, 36}
	squaredIter := iter.Map(func(i int) int { return i * i })
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
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(i int) bool { return i%2 == 0 }).
		Map(func(i int) int { return i * i }).
		Collect()
	expected := []int{4, 16, 36}
	if len(expected) != len(iter) {
		t.Error("Chaining did not work for ints")
		t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if iter[i] != expected[i] {
			t.Error("Chaining did not work for ints")
			t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
			return
		}
	}
}

func TestIterator_Skip(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6}).Skip(3).Collect()
	expected := []int{4, 5, 6}
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
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6}).StepBy(2).Collect()
	expected := []int{1, 3, 5}
	if len(expected) != len(iter) {
		t.Error("StepBy did not work for ints")
		t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if iter[i] != expected[i] {
			t.Error("StepBy did not work for ints")
			t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
			return
		}
	}
}

func TestIterator_Chain(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6})
	other := FromSlice([]int{7, 8, 9})
	result := iter.Chain(other).Collect()
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if len(expected) != len(result) {
		t.Error("Chain did not work for ints")
		t.Errorf("iter: %v\nexpected: %v\n", result, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {
			t.Error("Chain did not work for ints")
			t.Errorf("iter: %v\nexpected: %v\n", result, expected)
			return
		}
	}
}

func TestIterator_Intersperse(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3}).Intersperse(5).Collect()
	expected := []int{1, 5, 2, 5, 3}
	if len(expected) != len(iter) {
		t.Error("Intersperse did not work for ints")
		t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if iter[i] != expected[i] {
			t.Error("Intersperse did not work for ints")
			t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
			return
		}
	}
}

func TestIterator_ForEach(t *testing.T) {
	n := 0
	iter := FromSlice([]int{1, 2, 3})
	iter.ForEach(func(x int) { n += x })
	if n != 6 {
		t.Errorf("ForEach did not work for ints\niter: %d\nexpected: %d\n", n, 6)
	}
}

func TestIterator_SkipWhile(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6}).SkipWhile(func(n int) bool { return n < 4 }).Collect()
	expected := []int{4, 5, 6}
	if len(expected) != len(iter) {
		t.Error("SkipWhile did not work for ints")
		t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if iter[i] != expected[i] {
			t.Error("SkipWhile did not work for ints")
			t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
			return
		}
	}
}

func TestIterator_TakeWhile(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6}).TakeWhile(func(n int) bool { return n < 4 }).Collect()
	expected := []int{1, 2, 3}
	if len(expected) != len(iter) {
		t.Error("SkipWhile did not work for ints")
		t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if iter[i] != expected[i] {
			t.Error("SkipWhile did not work for ints")
			t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
			return
		}
	}
}

func TestIterator_Take(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6}).Take(3).Collect()
	expected := []int{1, 2, 3}
	if len(expected) != len(iter) {
		t.Error("Take did not work for ints")
		t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
		return
	}
	for i := 0; i < len(expected); i++ {
		if iter[i] != expected[i] {
			t.Error("Take did not work for ints")
			t.Errorf("iter: %v\nexpected: %v\n", iter, expected)
			return
		}
	}
}

func TestIterator_Inspect(t *testing.T) {
	n := 0
	iter := FromSlice([]int{1, 2, 3}).Inspect(func(x int) { n += x })
	iSlice := iter.Collect()
	expected := []int{1, 2, 3}
	if n != 6 {
		t.Errorf("Inspect did not work for ints\niter: %d\nexpected: %d\n", n, 6)
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
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6})
	even, odd := []int{2, 4, 6}, []int{1, 3, 5}
	left, right := iter.Partition(func(x int) bool { return x%2 == 0 })
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
	iter := FromSlice([]int{1, 2, 3})
	n := iter.Fold(0, func(acc, x int) int { return acc + x })
	if n != 6 {
		t.Errorf("Fold did not work for ints\niter: %d\nexpected: %d\n", n, 6)
	}
}

func TestIterator_Reduce(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3})
	n := iter.Reduce(func(acc, x int) int { return acc + x })
	if *n != 6 {
		t.Errorf("Reduce did not work for ints\niter: %d\nexpected: %d\n", *n, 6)
	}
}

func TestIterator_All(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3})
	b := iter.All(func(x int) bool { return x < 100 })
	if !b {
		t.Error("All did not work for ints")
		return
	}

	iter = FromSlice([]int{1, 2, 3})
	b = iter.All(func(x int) bool { return x < 3 })
	if b {
		t.Error("All did not work for ints")
		return
	}
}

func TestIterator_Any(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3})
	b := iter.Any(func(x int) bool { return x%2 == 0 })
	if !b {
		t.Error("Any did not work for ints")
		return
	}

	iter = FromSlice([]int{1, 2, 3})
	b = iter.Any(func(x int) bool { return x == 0 })
	if b {
		t.Error("Any did not work for ints")
		return
	}
}

func TestIterator_Find(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3})
	b := iter.Find(func(x int) bool { return x%2 == 0 })
	if *b != 2 {
		t.Error("Find did not work for ints")
		return
	}

	iter = FromSlice([]int{1, 2, 3})
	b = iter.Find(func(x int) bool { return x == 0 })
	if b != nil {
		t.Error("Find did not work for ints")
		return
	}
}

func TestIterator_Position(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3})
	p := iter.Position(func(x int) bool { return x%2 == 0 })
	if *p != 2 {
		t.Error("Position did not work for ints")
		return
	}

	iter = FromSlice([]int{1, 2, 3})
	p = iter.Position(func(x int) bool { return x == 0 })
	if p != nil {
		t.Error("Position did not work for ints")
		return
	}
}

func TestZip(t *testing.T) {
	it1 := FromSlice([]int{1, 2, 3})
	it2 := FromSlice([]int{4, 5, 6})
	it3 := Zip(it1, it2)
	if !it3.All(func(pair Pair[int, int]) bool { return pair.x+3 == pair.y }) {
		t.Error("Zip did not work for ints")
	}
}
