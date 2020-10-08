package gofast

import (
	"sync"
	"github.com/fatih/color"
)

var synchronizer sync.WaitGroup

var logger = false

type Resolver struct{
	Done chan bool
	component string
}

func WaitAll(){	
	defer synchronizer.Wait()
}

func WorkerPool(nbThreads int,funct func(res Resolver)){
	synchronizer.Add(1)
	go func(){
		res := Resolver{make(chan bool, nbThreads),"WorkerPool"}
		for i:=0; i < nbThreads; i++ {
			go funct(res)
		}
		if status := <-res.Done; status == false{
			errorLog(res.component)
		}
		if logger {doneLog(res.component)}
		synchronizer.Done()
    }()
}

func Promise(funct func(res Resolver),then func(res Resolver),catch func(res Resolver)){	
	synchronizer.Add(1)
	go func(){
		res := Resolver{make(chan bool, 1),"Promise Init"}
		res_then := Resolver{make(chan bool, 1),"Promise Then"}
		res_catch := Resolver{make(chan bool, 1),"Promise Catch"}
		go funct(res)
		if status := <-res.Done; status != false{
			go then(res_then)
			promiseError(&res_then.Done,res_then.component)
		} else {
			go catch(res_catch)
			promiseError(&res_catch.Done,res_catch.component)
		}
	}()
}

func promiseError(res *chan bool, component string){
	if status := <-*res; status == false{
		errorLog(component)
	}
	if logger {doneLog(component)}
	synchronizer.Done() 
}

func ActivateLogs(act bool){
	logger = act
	if act {titleGoFast()}
}

func doneLog(component string){
	color.Green("[GOFAST] DONE : "+component)
}

func errorLog(component string){
	color.Red("[GOFAST] ERROR : "+component)
}

func titleGoFast(){
	color.Cyan("[GOFAST]\n\n")
}