package main

import (
	"fmt"

	channel "go-work/sync"
)

func main() {
	fmt.Print("ok")
	// goroutine.ProTest()
	// point.Process()
	// scheduler.TestTask()
	// emp := person.Employee{

	// 	EmployeeId: 1,
	// 	Person: person.Person{
	// 		Age:  18,
	// 		Name: "d大大",
	// 	},
	// }
	channel.SyncAuto()
}
