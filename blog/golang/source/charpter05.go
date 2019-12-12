package main

import (
    "sync"
    "fmt"
    "time"
)

func sum(s []int, c chan int) {
    sum := 0
    for _, v := range s {
        sum += v
        fmt.Println("sum", v, sum)
    }    

    c <- sum
}

func closeTest() {
    ch := make(chan int ,2)

    ch <- 1
    x, isOpen := <-ch
    fmt.Println(x, isOpen)

    ch <- 2
    close(ch)
    x, isOpen = <-ch
    fmt.Println(x, isOpen)

    x, isOpen = <-ch
    fmt.Println(x, isOpen)

    ch <- 3
    x, isOpen = <-ch
    fmt.Println(x, isOpen)
}

func main() {
    /*
    closeTest()

    bufferdChannel()

    c := make(chan int)
    s := []int{1,2,3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
    go sum(s[:len(s)/2], c)
    go sum(s[len(s)/2:], c)
    x, y := <-c, <-c
    fmt.Println(x, y)

    selectTest()
    */
    
    mutexTest()
}


func bufferdChannel() {
    ch := make(chan int, 2)
    ch <- 1    
    ch <- 2
    //再放进去就会造成死锁
    //ch <- 3
    fmt.Println(<-ch)
    //这时候可以放了
    ch <- 3
    fmt.Println(<-ch)
    fmt.Println(<-ch)
    //再取也会造成死锁
    //fmt.Println(<-ch)
}
    
func selectTest() {
    c := make(chan int)
    quit := make(chan int)
    go func() {
        for i := 0; i < 10; i++ {
            fmt.Println(<-c)   
        }
        quit <- 0
    }()
    fibonacci(c, quit)
}

func fibonacci(c, quit chan int) {
    x, y := 0, 1
    for{
        select {
        case c <-x:
            x, y = y, x+y
        case <-quit:
            fmt.Println("quit")
            return    
        }    
    }        
}

type SafeCounter struct {
    v map[string]int
    mux sync.Mutex    
}


func (c *SafeCounter) Value(key string) int {
    c.mux.Lock()
    defer c.mux.Unlock()    
    return c.v[key]
}
func (c *SafeCounter) Inc(key string) {
    c.mux.Lock()
    c.v[key]++
    c.mux.Unlock()    
}

func mutexTest() {
    c := SafeCounter{v: make(map[string]int)}
    for i :=0; i< 1000; i++ {
        go c.Inc("somekey")    
    }

    time.Sleep(time.Second)

    fmt.Println(c.Value("somekey"))
}
