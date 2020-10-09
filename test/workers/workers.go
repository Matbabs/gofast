package main

import (
	"time"
	"github.com/matbabs/gofast"
	"fmt"
)


func worker(res gofast.Resolver) {
	fmt.Println("worker")
    time.Sleep(2000 * time.Millisecond)
    res.Done <- true
}

func main(){
	gofast.ActivateLogs(true)
    gofast.WorkerPool(10,worker)
	fmt.Println("main")
    gofast.WaitAll()
}