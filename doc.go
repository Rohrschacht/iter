// Package iter implements a generic iterator type.
//
// # About
//
// This package implements Rust-inspired iterators using Go 1.18 generics.
// Internally, the iterators are implemented using Goroutines and channels. This,
// using the provided methods on the iterators, one can define a pipeline that
// automatically uses multiple threads.
//
// # Examples
//
//	it := iter.FromSlice([]int{1, 2, 3, 4, 5, 6}).
//		Filter(func(i int) bool { return i%2 == 0 }).
//		Map(func(i int) int { return i * i }).
//		Collect()
//	expected := []int{4, 16, 36}
package iter
