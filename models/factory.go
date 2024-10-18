package models

type Factory[T any] interface {
	New() T
}

type ArrayFactory[T any] interface {
	NewArray(size int) []T
	NewAt(index int) T
}
