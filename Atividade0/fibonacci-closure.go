package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func(int) int {
	var atual, ant, antant int

	return func(n int) int {
		if n == 0 {
			return 0
		}
		if n == 1 {
			return 1
		}
		antant = 0
		ant = 1
		atual = 1
		for i := 2; i < n; i++ {
			antant = ant
			ant = atual
			atual = ant + antant
		}
		return atual
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f(i))
	}
}