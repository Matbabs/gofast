package main

import (
	"time"
	"github.com/matbabs/gofast"
	"fmt"
)

func worker(res gofast.Resolver) {

	gofast.Lock("myMutex")

	fmt.Println("critical section")
	time.Sleep(1000 * time.Millisecond)

	gofast.Unlock("myMutex")

	res.Done <- true
}

func main(){

	gofast.InitMutex("myMutex")

	gofast.WorkerPool(10,worker)

    gofast.WaitAll()
}