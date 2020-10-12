package gofast

import (
	"sync"
	"github.com/fatih/color"
)

var synchronizer sync.WaitGroup

var sems map[string]chan int

var mutexs map[string]chan int

var logger = false

type Resolver struct{
	Done chan bool
	component string
	capacity int
}

func WaitAll(){	
	defer synchronizer.Wait()
}

func WorkerPool(nbThreads int,funct func(res Resolver)){
	synchronizer.Add(1)
	go func(){
		res := Resolver{make(chan bool, nbThreads),"WorkerPool", nbThreads}
		for i:=0; i < nbThreads; i++ {
			go funct(res)
		}
		manageSynchro(res)
    }()
}

func Promise(funct func(res Resolver),then func(res Resolver),catch func(res Resolver)){	
	synchronizer.Add(1)
	go func(){
		res := Resolver{make(chan bool, 1),"Promise Init",1}
		res_then := Resolver{make(chan bool, 1),"Promise Then",1}
		res_catch := Resolver{make(chan bool, 1),"Promise Catch",1}
		go funct(res)
		if status := <-res.Done; status {
			go then(res_then)
			manageSynchro(res_then)
		} else {
			go catch(res_catch)
			manageSynchro(res_catch)
		}
	}()
}

func manageSynchro(res Resolver){
	for i:=0; i < res.capacity ; i++ {
		if status := <-res.Done; !status {
			errorLog(res.component)
		}
	}
	if logger {doneLog(res.component)}
	synchronizer.Done() 
}

func InitMutex(id string){
	if mutexs == nil {
		mutexs = make(map[string]chan int)
	}
	mutexs[id] = make(chan int, 1)
}

func DeleteMutex(id string){
	delete(mutexs,id)
}

func Lock(id string){
	mutexs[id] <- 1
	if logger {inCriticalLog()}
}

func Unlock(id string){
	if logger {outCriticalLog()}
	<-mutexs[id]
}

func InitSemaphore(id string, nbSemaphores int){
	if sems == nil {
		sems = make(map[string]chan int)
	}
	sems[id] = make(chan int, nbSemaphores)
}

func DeleteSemaphore(id string){
	delete(sems,id)
}

func Acquire(id string){
	sems[id] <- 1
	if logger {inSemLog()}
}

func Release(id string){
	if logger {outSemLog()}
	<-sems[id]
}

func ActivateLogs(act bool){
	logger = act
	if act {titleGoFast()}
}

func inCriticalLog(){
	color.Yellow("[GOFAST] IN CRITICAL SECTION")
}

func outCriticalLog(){
	color.Yellow("[GOFAST] OUT CRITICAL SECTION")
}

func inSemLog(){
	color.Yellow("[GOFAST] IN SEM SECTION")
}

func outSemLog(){
	color.Yellow("[GOFAST] OUT SEM SECTION")
}

func doneLog(component string){
	color.Green("[GOFAST] SYNCHRO DONE : "+component)
}

func errorLog(component string){
	color.Red("[GOFAST] ERROR : "+component)
}

func titleGoFast(){
	color.Cyan("[GOFAST]\n\n")
}