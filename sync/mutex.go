package sync

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

//编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

// 计数器
func getId() int {

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.Trim(string(buf[:n]), "goroutine"))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		fmt.Print(id)
	}
	return getId()
}

type Counter struct {
	mu    sync.Mutex
	count int
}

func (c *Counter) Incream() (int, int) {

	c.mu.Lock()
	defer c.mu.Unlock()
	id := getId()
	old := c.count

	c.count++
	return id, old
}
func SyncMutex() {

	var wg sync.WaitGroup
	var mx sync.Mutex //互斥锁
	var counter Counter
	goAddtion := make(map[int]int)

	for i := 1; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 1; i < 10; i++ {
				id, old := counter.Incream()
				mx.Lock()
				//开始计数新增

				goAddtion[id] += 1

				mx.Unlock()
				print(old)

			}

		}()
	}
	wg.Wait()

	for id, count := range goAddtion {
		fmt.Printf("新增 %d ： %d ", id, count)
	}
}

// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值
func SyncAuto() {
	var wg sync.WaitGroup
	var count int64

	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			for i := 0; i < 1000; i++ {
				atomic.AddInt64(&count, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println(count)

}
