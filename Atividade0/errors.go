package main

import (
	"fmt"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprint("cannot Sqrt negative number: ", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x >= 0.0 {
		z := 1.0
		for v:=0; v<10; v++ {
			z -= (z*z - x) / (2*z)
		}
		return z, nil
	}
	return 0, ErrNegativeSqrt(x)
}

func main() {
	var x float64
	x = -2
	if v, err := Sqrt(x); err == nil {
		fmt.Println(v)
	} else {
		fmt.Println(err)
	}
}
