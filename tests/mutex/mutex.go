package main

import (
	"time"
	"../../../gofast"
	"fmt"
)


func worker(res gofast.Resolver) {

	gofast.Lock()

	fmt.Println("critical section")
	time.Sleep(2000 * time.Millisecond)

	gofast.Unlock()


	res.Done <- true
}

func main(){
	gofast.ActivateLogs(true)
    gofast.WorkerPool(10,worker)
	fmt.Println("main")
    gofast.WaitAll()
}