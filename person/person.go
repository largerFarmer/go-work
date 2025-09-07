package person

import (
	"fmt"
)

/**
题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
*/

type Person struct {
	Name string
	Age  int
}
type Employee struct {
	EmployeeId int
	Person     Person
}

func (e Employee) PrintInfo() {

	fmt.Printf("员工姓名 d%, 年龄d% ,员工编号d%", e.Person.Name, e.Person.Age, e.EmployeeId)
}
