package channel

import (
	"fmt"
	"sync"
)

//一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。

func Channel() {
	ch := make(chan int)

	go func() {
		for i := 1; i < 10; i++ {
			ch <- i
		}
		defer close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
	}

}

// 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。

func ChangeWaitGroup() {

	ch := make(chan int, 10)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			ch <- i
		}
		close(ch)

	}()

	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Printf("消费者：%d\n", num)
		}
	}()
	wg.Wait()

}
