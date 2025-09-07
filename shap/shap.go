package shap

import (
	"math"
)

/**
定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
*/

type Shape interface {
	Area() float64      //面积
	Perimeter() float64 //周长
}

type Rectangle struct {
	chang float64
	kuan  float64
}

func (r Rectangle) Area() float64 {

	return float64(r.chang) * float64(r.kuan)
}
func (r Rectangle) Perimeter() float64 {
	return (r.chang + r.kuan) * 2
}

type Circle struct {
	Rdduis float64 //半径
}

func (r Circle) CircleArea() float64 {
	return math.Pi * r.Rdduis * r.Rdduis
}
func (r Circle) CirclePermter() float64 {
	return 2 * math.Pi * r.Rdduis
}
