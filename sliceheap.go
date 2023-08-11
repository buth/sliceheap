// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sliceheap provides heap operations for slices of any type. A heap is
// a tree with the property that each node is the minimum-valued node in its
// subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// A heap is a common way to implement a priority queue. To build a priority
// queue, implement the Heap interface with the (negative) priority as the
// ordering for the Less method, so Push adds items while Pop removes the
// highest-priority item from the queue. The Examples include such an
// implementation; the file example_pq_test.go has the complete source.
package sliceheap

import "cmp"

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = len(h).
func Init[T cmp.Ordered](h []T) {
	InitFunc(h, cmp.Less)
}

// InitFunc is like [Init] but uses a less function to compare elements.
func InitFunc[T any](h []T, less func(x, y T) bool) {
	// heapify
	n := len(h)
	for i := n/2 - 1; i >= 0; i-- {
		down(h, i, n, less)
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = len(h).
func Push[T cmp.Ordered](h *[]T, x T) {
	PushFunc(h, x, cmp.Less)
}

// PushFunc is like [Push] but uses a less function to compare elements.
func PushFunc[T any](h *[]T, x T, less func(x, y T) bool) {
	*h = append(*h, x)
	up(*h, len(*h)-1, less)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = len(h).
// Pop is equivalent to Remove(h, 0).
func Pop[T cmp.Ordered](h *[]T) T {
	return PopFunc(h, cmp.Less)
}

// PopFunc is like [Pop] but uses a less function to compare elements.
func PopFunc[T any](h *[]T, less func(x, y T) bool) T {
	n := len(*h) - 1
	x := (*h)[0]
	(*h)[0] = (*h)[n]
	*h = (*h)[:n]
	down(*h, 0, n, less)
	return x
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = len(h).
func Remove[T cmp.Ordered](h *[]T, i int) T {
	return RemoveFunc(h, i, cmp.Less)
}

// RemoveFunc is like [Remove] but uses a less function to compare elements.
func RemoveFunc[T any](h *[]T, i int, less func(x, y T) bool) T {
	n := len(*h) - 1
	x := (*h)[i]
	if n != i {
		(*h)[i] = (*h)[n]
		if !down(*h, i, n, less) {
			up(*h, i, less)
		}
	}
	*h = (*h)[:n]
	return x
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = len(h).
func Fix[T cmp.Ordered](h []T, i int) {
	FixFunc(h, i, cmp.Less)
}

// FixFunc is like [Fix] but uses a less function to compare elements.
func FixFunc[T any](h []T, i int, less func(x, y T) bool) {
	if !down(h, i, len(h), less) {
		up(h, i, less)
	}
}

func up[T any](h []T, j int, less func(x, y T) bool) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !less(h[j], h[i]) {
			break
		}
		h[i], h[j] = h[j], h[i]
		j = i
	}
}

func down[T any](h []T, i0, n int, less func(x, y T) bool) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && less(h[j2], h[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !less(h[j], h[i]) {
			break
		}
		h[i], h[j] = h[j], h[i]
		i = j
	}
	return i > i0
}
