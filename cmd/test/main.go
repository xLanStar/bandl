package main

import (
	"fmt"
)

type A struct {
	data map[string]bool
}

func f(data *A) {
	data.data["a"] = true
}

func g(data *A) {
	data.data["b"] = true
}

func main() {
	v := &A{
		data: make(map[string]bool),
	}
	fmt.Println(v)
	f(v)
	fmt.Println(v)
	g(v)
	fmt.Println(v)
}
