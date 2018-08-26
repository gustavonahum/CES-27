package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	allStrings := strings.Fields(s)
	m := make(map[string]int)
	for _, word := range allStrings {
		m[word]++
	}
	return m
}

func main() {
	wc.Test(WordCount)
}
