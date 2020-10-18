// GoFast - Matbabs 2020

// It makes it possible to parallel pools of workers,
// to make promises, mutex, and semaphores.

// The goal is to deploy them in a simple and light way.

// The concept is to control all this via father threads
// synchronized with the "sync.WaitGroup" (limiting its
// use) on anonymous functions. This way the main code
// remains concurrent throughout the execution of the
// program.

package gofast

import (
	"github.com/fatih/color"
	"sync"
)

// I use the WaitGroup in a limited way because the goal
// is to have a lighter implementation. It is used for
// the synchronization of the father threads. The son
// threads being managed by simple chan ( easier).
var synchronizer sync.WaitGroup

// These maps serve to facilitate the naming of mutexes and
// semaphores. These are developed with simple chan. I use
// the syntactic capabilities of the language and make writing
// easier for the user.
var sems map[string]chan int
var mutexs map[string]chan int

// Mutexs to block data race on package global maps
var semsMutex = make(chan int, 1)
var mutexsMutex = make(chan int, 1)

// Enable / Disable debugger
var logger = false

// Lets you know the resolution status and the name of the
// component being parallelized.
type Resolver struct {
	Done      chan bool
	component string
	capacity  int
}

// Declares a thread pool concurrently (thanks to the anonymous
// function itself launched in a parent thread). Avoids the
// user to write the syntax of the "for" loop.
func WorkerPool(nbThreads int, funct func(res Resolver)) {
	synchronizer.Add(1)
	go func() {
		res := Resolver{make(chan bool, nbThreads), "WorkerPool", nbThreads}
		for i := 0; i < nbThreads; i++ {
			go funct(res)
		}
		manageSynchro(res)
	}()
}

// The channels are used to make asynchronous requests.
// As WorkerPool the complete block is parallelized.
func Promise(funct func(res Resolver), then func(res Resolver), catch func(res Resolver)) {
	synchronizer.Add(1)
	go func() {
		res := Resolver{make(chan bool, 1), "Promise Init", 1}
		res_then := Resolver{make(chan bool, 1), "Promise Then", 1}
		res_catch := Resolver{make(chan bool, 1), "Promise Catch", 1}
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

// Makes sure that the threads end correctly.
func manageSynchro(res Resolver) {
	for i := 0; i < res.capacity; i++ {
		if status := <-res.Done; !status {
			errorLog(res.component)
		}
	}
	if logger {
		doneLog(res.component)
	}
	synchronizer.Done()
}

// It is used for the ending synchronization of the father threads.
func WaitAll() {
	defer synchronizer.Wait()
}

// Init a Mutex chan in global map. Protect data race access
// with a private package Mutex.
func InitMutex(id string) {
	mutexsMutex <- 1
	if mutexs == nil {
		mutexs = make(map[string]chan int)
	}
	mutexs[id] = make(chan int, 1)
	<-mutexsMutex
}

func DeleteMutex(id string) {
	mutexsMutex <- 1
	delete(mutexs, id)
	<-mutexsMutex
}

func Lock(id string) {
	mutexs[id] <- 1
	if logger {
		inCriticalLog()
	}
}

func Unlock(id string) {
	if logger {
		outCriticalLog()
	}
	<-mutexs[id]
}

// Init a Semaphore chan in global map. Protect data race access
// with a private package Mutex.
func InitSemaphore(id string, nbSemaphores int) {
	semsMutex <- 1
	if sems == nil {
		sems = make(map[string]chan int)
	}
	sems[id] = make(chan int, nbSemaphores)
	<-semsMutex
}

func DeleteSemaphore(id string) {
	semsMutex <- 1
	delete(sems, id)
	<-semsMutex
}

func Acquire(id string) {
	sems[id] <- 1
	if logger {
		inSemLog()
	}
}

func Release(id string) {
	if logger {
		outSemLog()
	}
	<-sems[id]
}

// All the functions below are intended to handle logging and debugging.

func ActivateLogs(act bool) {
	logger = act
	if act {
		titleGoFast()
	}
}

func inCriticalLog() {
	color.Yellow("[GOFAST] IN CRITICAL SECTION")
}

func outCriticalLog() {
	color.Yellow("[GOFAST] OUT CRITICAL SECTION")
}

func inSemLog() {
	color.Yellow("[GOFAST] IN SEM SECTION")
}

func outSemLog() {
	color.Yellow("[GOFAST] OUT SEM SECTION")
}

func doneLog(component string) {
	color.Green("[GOFAST] SYNCHRO DONE : " + component)
}

func errorLog(component string) {
	color.Red("[GOFAST] ERROR : " + component)
}

func titleGoFast() {
	color.Cyan("[GOFAST] v1.0.0\n\n")
}
