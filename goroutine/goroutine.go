package goroutine

import (
	"fmt"
	"sync"
)

//题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
//考察点 ： go 关键字的使用、协程的并发执行。

func ProTest() {
	channel := make(chan struct{}) //自定义结构体通道
	var wg sync.WaitGroup

	wg.Add(2) //两个需要等待的goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i += 2 {
			fmt.Printf("奇数 d%", i)
			//通知偶数协程
			channel <- struct{}{}
			if i < 9 {
				<-channel
			}

		}
	}()

	//打印偶数数据
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i += 2 {
			<-channel //接受奇数协程通知
			fmt.Printf("偶数", i)
			if i < 10 {
				channel <- struct{}{}

			}
		}
	}()
	wg.Wait()
	close(channel)

}
