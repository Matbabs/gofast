package main

import (
	"fmt"
	"math"
	"github.com/matbabs/gofast"
)

type Step struct{
	start int
	inc int
}

func gofast_pi(res gofast.Resolver){
	step := <-scatter
	pi := 0.0
	for k := float64(step.start) ; k <= float64(step.inc); k++ {
		pi += 4 * math.Pow(-1, k) / (2*k + 1)
	}
	gather <- pi
	res.Done <- true
}

var NB_THREADS = 50
var scatter = make(chan Step,NB_THREADS)
var gather = make(chan float64,NB_THREADS)

func main(){
	var steps = 100000000
	pi := 0.0
	gofast.WorkerPool(NB_THREADS,gofast_pi)
	block := ((steps)/NB_THREADS)
	for i:= 0; i < NB_THREADS ; i++{ scatter<-Step{block*i,block} }
	for i:= 0; i < NB_THREADS ; i++{ pi += <-gather }
	fmt.Println(pi)
	gofast.WaitAll()
}