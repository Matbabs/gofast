package main

import (
	"time"
	"github.com/matbabs/gofast"
	"fmt"
)

func worker(res gofast.Resolver){
	for i := 0; i < 3; i++ {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("worker")
	}
	res.Done <- true
}

func asyncFunction(res gofast.Resolver){
	fmt.Println("new promise")
	time.Sleep(3000 * time.Millisecond)
	res.Done <- true
}

func asyncFunction_Then_(res gofast.Resolver){
	fmt.Println("then")
	res.Done <- true
}

func asyncFunction_Catch_(res gofast.Resolver){
	fmt.Println("catch")
	res.Done <- true
}


func main(){
	gofast.ActivateLogs(true)
	gofast.WorkerPool(2,worker)
	gofast.Promise(asyncFunction,asyncFunction_Then_,asyncFunction_Catch_)
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("main program")
	gofast.WaitAll()
}