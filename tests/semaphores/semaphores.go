package main

import (
	"time"
	"github.com/matbabs/gofast"
	"fmt"
)

func worker(res gofast.Resolver){
	
	gofast.Acquire("mySem")

	fmt.Println("in semaphore")
	time.Sleep(2000 * time.Millisecond)

	gofast.Release("mySem")
	
	res.Done <- true
}

func main(){

	gofast.InitSemaphore("mySem",2)

	gofast.WorkerPool(6,worker)

	gofast.WaitAll()
}