//go:build beta

package main

func TestingPanic() {
	panic("panic in beta build")
}
