package main

import (
    "fmt"
    "math"
)

func sequential_pi(n float64) float64 {

    pi := 0.0
    
    for k := 0.0; k <= float64(n); k++ {
        pi += 4 * math.Pow(-1, k) / (2*k + 1)
    }
    
    return pi
}

func main(){

	var steps = 100000000	

	pi := sequential_pi(float64(steps))
	fmt.Println(pi)

}