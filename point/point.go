package point

import (
	"fmt"
)

const pkgName string = "pkg1"

var pp *int //声明一个指针类型变量需使用星号 * 标识：

func Process() {

	i := 19 //基础数据类型必须通过变量获取指针，因为字面量在编译期间会生成常量，不能获取到内存中的指针信息
	pp = &i //初始化指针必须通过另外一个变量，
	var s1 []int = []int{3, 5, 5, 6}

	zhiZhen(pp)
	sliceTest(&s1)
}

// 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
// 考察点 ：指针的使用、值传递与引用传递的区别。
// 在 Go 中，声明一个指针类型变量需使用星号 * 标识：
func zhiZhen(p1 *int) int {
	*p1 = *p1 + 10
	fmt.Println("指针加10", *p1)
	return *p1
}

// 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
// 考察点 ：指针运算、切片操作。
func sliceTest(slices *[]int) {

	for k, v := range *slices {
		(*slices)[k] *= 2
		fmt.Println(k, v, (*slices)[k])
	}
}
