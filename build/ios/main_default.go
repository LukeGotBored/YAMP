//go:build !ios

package main

// main exists so `go build ./...` succeeds on non-iOS hosts.
func main() {}
