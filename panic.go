//go:build !beta

package main

func TestingPanic() {
	panic("panic in non-beta build")
}
