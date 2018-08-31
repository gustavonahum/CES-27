package main

import (
	"fmt"
	"math"
)

func Sqrt1(x float64) float64 {
	z := 1.0
	for v:=0; v<10; v++ {
		z -= (z*z - x) / (2*z)
	}
	return z
}

func Sqrt2(x float64) (float64, int) {
	var z float64 = 1.0
	var z1, z2 float64
	var delta float64 = math.Abs(z - x)
	var cnt int = 0
	
	for delta > 0.00001 {
		z1 = z
		z -= (z*z - x) / (2*z)
		z2 = z
		delta = math.Abs(z1 - z2)
		cnt++
	}
	return z, cnt
}

func Sqrt3(x float64) (float64, int) {
	var z float64 = x
	var z1, z2 float64
	var delta float64 = 0.1
	var cnt int = 0
	
	for delta > 0.00001 {
		z1 = z
		z -= (z*z - x) / (2*z)
		z2 = z
		delta = math.Abs(z1 - z2)
		cnt++
	}
	return z, cnt
}

func Sqrt4(x float64) (float64, int) {
	var z float64 = x/2
	var z1, z2 float64
	var delta float64 = 0.1
	var cnt int = 0
	
	for delta > 0.00001 {
		z1 = z
		z -= (z*z - x) / (2*z)
		z2 = z
		delta = math.Abs(z1 - z2)
		cnt++
	}
	return z, cnt
}


func main() {
	z1 := Sqrt1(2)
	fmt.Printf("Método 1 (10 iterações): %v\n", z1)
	fmt.Printf("Módulo da diferença entre Sqrt1 e math.Sqrt: %v\n\n", math.Abs(z1 - math.Sqrt(2)))
	
	z2, cnt2 := Sqrt2(2)
	fmt.Printf("Método 2 (%v iterações): %v\n", cnt2, z2)
	fmt.Printf("Módulo da diferença entre Sqrt2 e math.Sqrt: %v\n\n", math.Abs(z2 - math.Sqrt(2)))
	
	z3, cnt3 := Sqrt3(2)
	fmt.Printf("Método 3 (%v iterações): %v\n", cnt3, z3)
	fmt.Printf("Módulo da diferença entre Sqrt3 e math.Sqrt: %v\n\n", math.Abs(z3 - math.Sqrt(2)))
	
	z4, cnt4 := Sqrt4(2)
	fmt.Printf("Método 4 (%v iterações): %v\n", cnt4, z3)
	fmt.Printf("Módulo da diferença entre Sqrt4 e math.Sqrt: %v\n\n", math.Abs(z4 - math.Sqrt(2)))
}
