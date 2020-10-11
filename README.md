![GoFast](assets/gofast-small.png)
# Multithread Programming Tool

![](https://img.shields.io/static/v1.svg?label=&message=GoFast&color=2cb6aa)
![](https://img.shields.io/static/v1.svg?label=&message=Multithread&color=2cb6aa)
![](https://img.shields.io/static/v1.svg?label=Tool&message=v0.0.1&color=edca9c)
![](https://img.shields.io/static/v1.svg?label=Worker&message=Pools&color=e760a3)
![](https://img.shields.io/static/v1.svg?label=Worker&message=Promises&color=e760a3)
![](https://img.shields.io/static/v1.svg?label=Worker&message=Mutex&color=edca9c)
![](https://img.shields.io/static/v1.svg?label=Worker&message=Semaphores&color=edca9c)



## Contents

* [Install](#install)
* [Worker Pools](#worker-pools)
    * [Examples](#examples)
    * [Native Way](#native-way)
    * [GoFast Way](#gofast-way)
    * [Scatter & Gather](#scatter-&-gather)
    * [Sequential Aproximation of PI](#sequential-aproximation-of-pi)
	* [Pthread C Aproximation of PI](#pthread-c-aproximation-of-pi)
    * [GoFast Aproximation of PI](#gofast-aproximation-of-pi)
    * [Benchmark Test](#benchmark-test)
* [Promises](#promises)
* [Mutex](#mutex)
* [Semaphores](#semaphores)
* [Manage Errors](#manage-errors)
* [Logs Display](#logs-display)


# Install

First check that __Go__ is correctly installed. [Download and install Go](https://golang.org/doc/install#testing)

To use __GoFast__ in your project  use the commands below:

`go get github.com/fatih/color` (dependency)

`go get github.com/matbabs/gofast`

Import also:

`"github.com/matbabs/gofast"`


# Worker Pools

## Examples

### Native Way
In this example, we will look at how to implement a worker pool natively. We will use goroutines and channels.
The goal is to run a function with 10 threads and wait for the end of their work. 

```go
func worker(results chan<- bool) {
    /* ... */
    done <- true
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
    res.Done <- true
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

### Pthread C Aproximation of PI

```c
typedef struct Step{
    double start;
    long inc;
    double res;
}Step;

void* c_pi(void *st){
    Step *step = st;   
    for(long k=step->start;k<step->start+step->inc;k++){
        step->res += 4.0 * pow(-1, k) / (2*k + 1);
    }
}

#define NB_THREADS 50

int main () { 
    static long nb_pas = 100000000;
    pthread_t  p_thread[NB_THREADS];
    Step steps[NB_THREADS];
    double pi,bloc; 
    for(int i=0;i<NB_THREADS;i++){
        bloc = nb_pas/NB_THREADS;
        steps[i].start = bloc*i;
        steps[i].inc = bloc;
        steps[i].res = 0;
        pthread_create(&p_thread[i],NULL, c_pi, &steps[i]);
        pthread_join(p_thread[i],NULL);
    }
    for(int i=0;i<NB_THREADS;i++)
        pi += steps[i].res;
    printf("PI=%f\n",pi);
    return 0;
}
```

### GoFast Aproximation of PI

```go
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
var scatter = make(chan Step, NB_THREADS)
var gather = make(chan float64, NB_THREADS)

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
```


### Benchmark Test

Operations: 100,000,000

|   |  Sequential PI | Pthread PI (50 threads)  | GoFast PI (50 threads) |
|:---:|:---:|:---:|:---:|
| res  | 3.141592663589326 | 3.141593 | 3.1415941535892244 |
|  real | 9,114s  | 1,483s  | 0,196s |
|  user | 9,113s | 1,478s | 0,190s |
|  sys | 0,012s  | 0,004s | 0,008s |

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
	gofast.Promise(asyncFunction,asyncFunction_Then_,asyncFunction_Catch_)
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

# Mutex

__GoFast__ itself uses group synchronization. To use Mutex it is necessary to use the functions defined by __GoFast__.

Use `gofast.InitMutex("<mutex name>")` to initialize a mutex.

Lock and unlock a critical section with `gofast.Lock("<mutex name>")` & `gofast.Unlock("<mutex name>")`

```go
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
```

# Semaphores

__GoFast__ itself uses group synchronization. To use semaphores you can use the functions defined by __GoFast__. 

Use `gofast.InitSemaphore("<semaphore name>",<semaphore capacity>)` to initialize a semaphore.

Acquire and realease a section with `gofast.Acquire("<semaphore name>")` & `gofast.Release("<semaphore name>")`

```go
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
```


# Manage Errors

__GoFast__ allows to the developper a manager error system.

You can indicate to the `gofast.Resolver` that no problem has been thrown in the program:

`res.Done <- true`

To indicate an error has been thrown you can use:

`res.Done <- false`

> NB: in case of a new `gofast.Promise()` the code will execute the catch function. 

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

