package main
import (
    "fmt"
)

func main1() {
//func main() {
    //pointerTest() 
    //structTest()
    //arrayTest()
    //sliceTest()
    //rangeTest()
    //mapTest()
    add := func(x, y int) int {
        return x + y    
    }
    fmt.Println(add( 1, 2))
    fmt.Println(compute(add, 1, 2))

    add = functionClosure();
    fmt.Println(add(1, 2))
}

func compute(fn func(int, int) int, x int, y int) int {
    return fn(x, y)    
}

func functionClosure() func(int, int) int {
    return func(x, y int) int {
        return x + y   
    }    
}

func mapTest(){
    //声明而未初始化的map为 nil，不能赋值
    var m map[string]OneStruct
    //m["first"] = OneStruct{1,2}   
    fmt.Println(m)

    //声明并初始化
    var m1 = map[string]OneStruct{
        "first": OneStruct{1,1},
        "two": OneStruct{2,2},
    }
    m1["third"] = OneStruct{3, 3}
    fmt.Println(m1)

    //使用make进行初始化
    var m2 map[string]OneStruct = make(map[string]OneStruct)
    //操作：插入元素
    m2["first"] = OneStruct{1,1}
    fmt.Println(m2)
    //操作：查找元素 m[key]
    elem := m2["one"]
    //操作：查找元素 若 key存在，ok = true
    elem, ok := m2["two"]
    fmt.Println(elem, ok)
    //操作：删除元素
    delete(m,"key")

}


func pointerTest() {
    i := 12
    var p *int = &i
    fmt.Println(p)
    fmt.Println(*p)
    *p = 21
    fmt.Println(p)
    fmt.Println(*p)
    fmt.Println(i)
}

type OneStruct struct {
    X int
    Y int    
}

func structTest() {
    //初始化一个变量，参数列表用 , 分隔
    v := OneStruct{1, 2}
    fmt.Println(v) 

    p := &v
    //使用指针变量对结构体成员进行修改
    p.X = 22

    fmt.Println(v) 

    //取值符然后调用成员
    fmt.Println((*p).X) 
    //两者是等价的
    fmt.Println(v.X)
    
    //用取值符太麻烦了，所以语言可以让我们简写 
    fmt.Println(p.X) 
}

func arrayTest(){
    //声明方式
    var a [2]string
    a[0] = "Hello" 
    a[1] = "workd" 
    fmt.Println(a); 
    //引用单个成员
    fmt.Println(a[0]);
    
    //声明并赋值
    var b = [6]int{1,2,3}
    fmt.Println(b);

    //可以省略数组长度
    q := []int{1,2,3}
    fmt.Println(q)

    //对于结构体可以酱紫：
    s := [3]struct{
        i int
        b bool    
    }{
        {2, true},
        {3, false},
    }

    fmt.Println(s) //>>> [{2 true} {3 false} {0 false}]
}

func sliceTest(){
    var b = [6]int{1,2,3}
    fmt.Println(b)

    //从下标1取到下标2
    var c = b[1:2]
    //b[begin:end] begin 不设置的话为0，end不设置的话为数组长度-1
    fmt.Println(c)

    //切片只是引用原数组
    c[0] = 10
    fmt.Println(b) //>>> [1 10 3 0 0 0]
 
    s := []int{1,2,3,4,5,6,7}
    printSlice(s)

    s= s[:0]
    printSlice(s) //>>> len=0 cap=7 []

    //扩容
    s= s[:4]
    printSlice(s) //>>> len=4 cap=7 [1 2 3 4]

    //截取
    s= s[2:3]
    printSlice(s) //>>> len=1 cap=5 [3]

    //扩容：截取是不可逆的，但是可以扩容
    s= s[0:5] //>>> len=0 cap=5 [3 4 5 6]
    printSlice(s)

   e := make([]int, 5)
   printSlice(e) //>>> len=5 cap=5 [0 0 0 0 0]
   
   //参数分别是 type length caption
   f := make([]int, 0, 5)
   printSlice(f) // >>> len=0 cap=5 [ ]
}

func printSlice(s []int) {
    fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)    
}

func rangeTest(){
    var pow = []int{1,2,4,8,16}

    for i, v := range pow{
        fmt.Printf("2 ^ %d = %d\n", i, v)
    }
    
}
