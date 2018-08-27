package main

import "golang.org/x/tour/pic"

func Pic1(dx, dy int) [][]uint8 {
	pic := make([][]uint8, dy)
	for y := range pic {
		pic[y] = make([]uint8, dx)
		for x := range pic[y] {
			pic[y][x] = uint8((x+y)/2)
		}
	}
	return pic
}

func Pic2(dx, dy int) [][]uint8 {
	pic := make([][]uint8, dy)
	for y := range pic {
		pic[y] = make([]uint8, dx)
		for x := range pic[y] {
			pic[y][x] = uint8(x*y)
		}
	}
	return pic
}

func Pic3(dx, dy int) [][]uint8 {
	pic := make([][]uint8, dy)
	for y := range pic {
		pic[y] = make([]uint8, dx)
		for x := range pic[y] {
			pic[y][x] = uint8(x^y)
		}
	}
	return pic
}

func main() {
	pic.Show(Pic1)
	// pic.Show(Pic2)
	// pic.Show(Pic3)
}
