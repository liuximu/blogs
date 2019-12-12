package main

import "fmt"
import "runtime"

/*func main() {
    forTest()
    switchTest()

    defersTest()
}*/

func forTest() {
    sum := 1
    for ; sum < 1000; {
        sum += sum    
    }    

    fmt.Println(sum)
}

func switchTest() {
    switch os :=runtime.GOOS; os {
        case "darwin":
            fmt.Println("OS X.")
        case "linux":
            fmt.Println("Linx")
        default:
            fmt.Printf("%s.", os)    
    }    
}

func defersTest() {
    fmt.Println("counting") 
    
    for i := 0; i < 10; i++ {
        defer fmt.Println(i)    
    }   

    fmt.Println("done")
}
