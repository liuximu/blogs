//如果是main，那么会被编译成可执行命令
package main

//引入包既可以单条引入，也可以多条引入
import (
    "fmt"
    "math"
)

func Abs(v Vertex) float64 {
    return math.Sqrt(v.X * v.X + v.Y * v.Y)    
}


//func main() {
func main04() {
    v := Vertex{3, 4}
    //fmt.Println(Abs(v))

    valueReceiverTest(v, 1)
    fmt.Println(v)

    pointerReceiverTest(&v, 1)
    fmt.Println(v)

    (&v).pointerReceiverIndirection(1)
    fmt.Println(v)

    v.pointerReceiverIndirection(1)
    fmt.Println(v)

    result := v.valueReceiverIndirection(1)
    fmt.Println(result)

    result = (&v).valueReceiverIndirection(1)
    fmt.Println(result)
}

type Vertex struct {
    X, Y float64    
}

//当传值时，参数是个副本，其被修改但是不影响外界的变量
func valueReceiverTest(v Vertex, i float64) {
    v.X = v.X + i
}

//当使用指针传引用，值本身被改变
func pointerReceiverTest(v *Vertex, i float64) {
    //直接使用指针获取结构体成员是go语言支持的简便方式
    v.X = v.X + i
    //两者是等价的
    (*v).X = (*v).X + i
}

//当第一个参数是指针时，可以将其放到前面
func (v *Vertex) pointerReceiverIndirection(i float64) {
    v.X = v.X + i
}

//当第一个参数是指针时，可以将其放到前面
func (v Vertex) valueReceiverIndirection(i float64) float64{
   return v.X + i
}
