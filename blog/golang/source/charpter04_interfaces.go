package main

import (
    "fmt"
    "math"
)

//为非本地类型起别名
type MyFloat float64

type Abser interface {
    Abs() float64
}

func (f MyFloat) Abs() float64 {
    if f < 0 {
        return float64(-f)    
    }
    return float64(f)
}


func(v *Vertex) Abs() float64 {
    return math.Sqrt(v.X * v.X + v.Y * v.Y)    
}

func main4_1() {
    f := MyFloat(-math.Sqrt2)
    fmt.Println(f);
    fmt.Println(f.Abs())   
    
    //声明一个接口类型的变量
    var a Abser; 
    //MyFloat 实现了 Abser
    //应该类型通过实现一个接口的函数来实现它，不需要使用关键字implement
    a = f
    fmt.Println(a.Abs())   

    //接口也是有值的，它是一个元组(value, type)
    fmt.Printf("%v, %T \n", a, a) //>> -1.4142135623730951, main.MyFloat5

    v := Vertex{3, 4}
    //不使用接口时也是可以用的
    fmt.Println(v.Abs())   
    // Vertex 实现接口 Abser
    a = &v
    fmt.Println(a.Abs())   

    fmt.Printf("%v, %T \n", a, a) //>> &{3 4}, *main.Vertex

    var empty Abser
    //空对象
    fmt.Printf("%v, %T \n", empty) //>> <nil>, %!T(MISSING)
    //会报运行时错误
    //empty.Abs()

    //值非空，类型空
    fmt.Printf("%v, %T \n", f) //>> -1.4142135623730951, %!T(MISSING) 

    //空接口，没有任何函数的接口，可以被任何类型实现，常用来处理未知类型的值
    var i interface { }
    i = 4
    fmt.Printf("%v, %T \n", i) //>> 4, %!T(MISSING)  

    //类型断言的语法和map的获取值的语法相似 value, ok := i.(T)
    //如果断言成功，value被正常赋值，ok为true，不然value为T的零值
    r, ok := i.(int)
    fmt.Println(r, ok) //>>> 4, true
    r1, ok := i.(float64)
    fmt.Println(r1, ok) //>>> 0, false

    //可以得到接口的类型，用在switch中
    switch v := i.(type){
        default:
            fmt.Printf("%T\n", v)
    }

    //Stringer 是最常用的接口，里面有一个函数是String
    /**
    type Stringer interface {
        String() string    
    }
    */
    p1 := Person{"aaa", 111}
    p2 := Person{"bbb", 11}
    fmt.Println(p1, p2)

    b1 := Person{"b1", 111}
    b2 := Person{"b2", 11}
    fmt.Println(b1, b2)

    //Errors接口是用来处理错误的
    /*
    type error interface {
        Error()    
    }
    */
    //更多常用包我们在3.\*中继续
}

type Boy struct {
    Name string
    Age int        
}

type Person struct {
    Name string
    Age int        
}
//对类型 Person 的String进行重载
func (p Person) String() string {
    return fmt.Sprintf("%v %v", p.Name, p.Age)
}

//对类型 Boy 类型 的String进行重载
func (p Boy) String() string {
    return fmt.Sprintf("%v %v", p.Name, p.Age)
}
