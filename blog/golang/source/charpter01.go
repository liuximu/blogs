//如果是main，那么会被编译成可执行命令
package main

//引入包既可以单条引入，也可以多条引入
import (
    "fmt"
    "math"
)

func main01() {
//func main() {
    fmt.Println(math.Pi)
    fmt.Println(add(1, 3))

    // := 简洁语法在函数内部可用an
    x, y := swap(1, 3)
    fmt.Println(x, y)

    typeConvert();

    constTest() // >>> 3.14
    const PI = 4
    //常量只在相应范围有效
    fmt.Println(PI) //>>> 4
    constTest() // >>> 3.14
}


//函数的格式就是这样 注意形参列表和返回列表
func add(x, y int)(int) {
    return x + y    
}

func swap(x int, y int)(a int, b int) {
    a = y
    b = x
    return   
}


//类型转换需要显式转换
func typeConvert() {
    var i int = 42
    fmt.Println(float64(i))   
}

func constTest(){
    const PI = 3.14
    fmt.Println(PI)  
}
