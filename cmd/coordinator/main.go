package main

import "github.com/namelew/RPC/internal/coordinator"

func main() {
	cd := coordinator.Build()
	cd.Handler()
}
