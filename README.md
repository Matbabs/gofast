![GoFast](assets/gofast-small.png)
# Multithread Programming Tool

## Contents

* [Install](#install)
* [Worker Pools](#worker-pools)
    * [Examples](#examples)
    * [Native Way](#native-way)
    * [GoFast Way](#gofast-way)
    * [Scatter & Gather](#scatter-&-gather)
    * [Sequential Aproximation of PI](#sequential-aproximation-of-pi)
    * [GoFast Aproximation of PI](#gofast-aproximation-of-pi)
    * [Benchmark Test](#benchmark-test)
* [Promises](#install)
* [Manage Errors](#manage-errors)
* [Logs Display](#logs-display)


# Install





First check that __Go__ is correctly installed. [Download and install Go](https://golang.org/doc/install#testing)

To use __GoFast__ in your project  use the commands below:

`go get github.com/fatih/color` (dependency)

`go get github.com/matbabs/gofast`


# Worker Pools

## Examples

### Native Way
In this example, we will look at how to implement a worker pool natively. We will use goroutines and channels.
The goal is to run a function with 10 threads and wait for the end of their work. 

```go
func worker(results chan<- bool) {
    /* ... */
    done <- 1
}

func main(){

    NB_THREADS := 10

    done := make(chan bool, NB_THREADS)

    for w := 0; w < NB_THREADS; w++ {
        go worker(done)
    }
    // make synchronization
    for w := 0; w < NB_THREADS; w++ {
        <-done
    }
    // ...
    // BUT
    // <program waits for the end of the synchronization>
    //...

}
```

### GoFast Way
Now we will show how GoFast can trigger 10 synchronized threads whose code is concurrent with the rest of the program.

```go
func worker(res gofast.Resolver) {
    /* ... */
    res.Done <- 1
}

func main(){

    NB_THREADS := 10

    gofast.WorkerPool(NB_THREADS,worker)
    // ...
    // <program is parallel and concurrent (no data to receive)>  
    // ...

    // secures the end of the threads at the end of the program
    gofast.WaitAll()
}
```

> __Here you can see that the rest of the code is parallel and concurrent until the synchronization is done.__


## Scatter & Gather

We are going to see how to parallelize the approximation of pi with 100,000,000 steps. The objective is to see the performance and how to send and retrieve data.

### Sequential Aproximation of PI

```go
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

```

### GoFast Aproximation of PI

```go
type Step struct{
	start int
	end int
}

func gofast_pi(res gofast.Resolver){
	n := <-scatter
	pi := 0.0
	for k := float64(n.start) ; k <= float64(n.end); k++ {
		pi += 4 * math.Pow(-1, k) / (2*k + 1)
	}
	gather <- pi
	res.Done <- true
}

func main(){

	var steps = 100000000
	var NB_THREADS = 10
	
	gofast.WorkerPool(NB_THREADS,gofast_pi)
	for i:= 0; i < NB_THREADS ; i++{ 
        scatter <- Step{((steps)/NB_THREADS)*i,((steps)/NB_THREADS)*(i+1)}}
	for i:= 0; i < NB_THREADS ; i++{ 
        pi += <-gather}
	fmt.Println(pi)
	
	gofast.WaitAll()
}
```


### Benchmark Test

|   |  Sequential PI | GoFast PI (10 threads)  | GoFast PI (20 threads) |
|:---:|:---:|:---:|:---:|
| time  |  9,051s | 2,008s | 1,942s |
|  res | 3.141592663589326  | 3.1415932293834268  | 3.1415940826855784 |

# Promises

Another use of multithreaded programming can be asynchronous requests. With __GoFast__ you can easily set up an asynchronous function. The main program will continue to run concurrently and in parallel.

```go
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
	gofast.WorkerPool(2,worker)
	gofast.Promise(asyncFunction,asyncFunctionThen,asyncFunctionCatch)
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("main program")
	gofast.WaitAll()
}
```

You can observe that this program is fully parralel.

```
new promise
worker
main program
worker
worker
worker
then
worker
```

# Manage Errors

__GoFast__ allows to the developper a manager error system.

You can indicate to the `gofast.Resolver` that no problem has been thrown in the program:

`res.Done <- true`

To indicate an error has been thrown you can use:

`res.Done <- false`

> NB: in case of a new `gofast.Promise()` the code will execute the catch fucntion. 

# Logs Display

`gofast.Resolver` can check if Worker Pool or Promise have succeed. 

To have synchronized finish job you can use: 

`gofast.ActivateLogs(true)`

```
[GOFAST]

new promise
main program
worker
worker
worker
worker
then
[GOFAST] DONE : Promise Then
worker
worker
[GOFAST] DONE : WorkerPool
```

