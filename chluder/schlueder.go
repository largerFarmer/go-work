package goroutine

import (
	"fmt"
	"sync"
	"time"
)

/**
设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
*/

// Task 定义任务接口
// 任何想要被调度器执行的任务都需要实现这个接口
type Task interface {
	Execute() error // 执行任务的方法，返回可能的错误
	Name() string   // 获取任务名称的方法
}

// TaskResult 任务执行结果结构体
// 用于记录每个任务的执行情况和时间统计
type TaskResult struct {
	TaskName  string        // 任务名称
	StartTime time.Time     // 任务开始时间
	EndTime   time.Time     // 任务结束时间
	Duration  time.Duration // 任务执行时长
	Error     error         // 任务执行错误（如果有）
	Success   bool          // 任务是否执行成功
}

// Scheduler 任务调度器结构体
// 负责管理和执行所有任务
type Scheduler struct {
	tasks       []Task          // 存储所有待执行的任务
	results     []TaskResult    // 存储所有任务的执行结果
	concurrency int             // 并发执行的任务数量（工作协程数）
	wg          sync.WaitGroup  // 等待组，用于等待所有工作协程完成
	mu          sync.Mutex      // 互斥锁，用于保护共享数据的并发访问
	completed   int             // 已完成的任务数量
	totalTasks  int             // 总任务数量
	resultChan  chan TaskResult // 通道，用于工作协程向结果收集器发送结果
	done        chan bool       // 通道，用于通知结果收集已完成
}

// NewScheduler 创建新的调度器实例
// concurrency: 并发执行的任务数量
func NewScheduler(concurrency int) *Scheduler {
	// 返回一个新的调度器实例
	return &Scheduler{
		concurrency: concurrency,                // 设置并发度
		resultChan:  make(chan TaskResult, 100), // 创建带缓冲的结果通道
		done:        make(chan bool),            // 创建完成通知通道
	}
}

// AddTask 添加任务到调度器
// task: 要添加的任务实例
func (s *Scheduler) AddTask(task Task) {
	// 将任务追加到任务列表中
	s.tasks = append(s.tasks, task)
}

// worker 工作协程函数
// taskChan: 任务通道，从中获取要执行的任务
// workerID: 工作协程ID，用于标识不同的工作协程
func (s *Scheduler) worker(taskChan chan Task, workerID int) {
	// 在函数退出时通知等待组该协程已完成
	defer s.wg.Done()

	// 循环从任务通道中获取任务，直到通道关闭
	for task := range taskChan {
		start := time.Now() // 记录任务开始时间

		// 执行任务并获取可能的错误
		err := task.Execute()
		// 计算任务执行时长
		duration := time.Since(start)
		// 判断任务是否执行成功（没有错误即为成功）
		success := err == nil

		// 将任务执行结果发送到结果通道
		s.resultChan <- TaskResult{
			TaskName:  task.Name(), // 任务名称
			StartTime: start,       // 开始时间
			EndTime:   time.Now(),  // 结束时间
			Duration:  duration,    // 执行时长
			Error:     err,         // 错误信息
			Success:   success,     // 成功标志
		}
	}
}

// resultCollector 结果收集器函数
// 负责从结果通道接收结果并更新调度器状态
func (s *Scheduler) resultCollector() {
	// 循环从结果通道接收结果，直到通道关闭
	for result := range s.resultChan {
		// 加锁保护共享数据
		s.mu.Lock()
		// 将结果追加到结果列表中
		s.results = append(s.results, result)
		// 增加已完成任务计数
		s.completed++
		// 计算完成进度百分比
		progress := float64(s.completed) / float64(s.totalTasks) * 100
		// 打印实时进度信息（\r使光标回到行首实现覆盖效果）
		fmt.Printf("\r进度: %.2f%% (%d/%d)", progress, s.completed, s.totalTasks)
		// 解锁
		s.mu.Unlock()
	}
	// 所有结果收集完成后，发送完成信号
	s.done <- true
}

// Run 执行所有任务
// 启动工作协程和结果收集器，等待所有任务完成
func (s *Scheduler) Run() {
	// 设置总任务数
	s.totalTasks = len(s.tasks)
	// 初始化已完成任务数为0
	s.completed = 0
	// 初始化结果切片，容量为总任务数
	s.results = make([]TaskResult, 0, s.totalTasks)

	// 创建任务通道，容量为总任务数
	taskChan := make(chan Task, s.totalTasks)

	// 启动结果收集器协程
	go s.resultCollector()

	// 启动工作协程，数量等于并发度
	s.wg.Add(s.concurrency)
	for i := 0; i < s.concurrency; i++ {
		// 为每个工作协程分配一个唯一的ID（从1开始）
		go s.worker(taskChan, i+1)
	}

	// 将所有任务发送到任务通道
	for _, task := range s.tasks {
		taskChan <- task
	}
	// 关闭任务通道，表示所有任务已发送完毕
	close(taskChan)

	// 等待所有工作协程完成
	s.wg.Wait()
	// 关闭结果通道，表示所有结果已发送完毕
	close(s.resultChan)

	// 等待结果收集完成
	<-s.done
	// 打印完成信息
	fmt.Println("\n所有任务执行完成!")
}

/**
%d：格式化整数（用于i+1）

%s：格式化字符串（用于result.TaskName和result.Duration）

%v：以默认格式格式化result.Error（可能是error类型或任何其他类型）

%v的特点
通用性：可以格式化任何类型的值

智能输出：

对于基本类型（int, string, bool等），输出其自然形式

对于结构体，输出字段名和值

对于指针，输出指针指向的值（而不是地址）

对于error类型，输出错误的描述信息
*/

// PrintReport 打印执行报告
// 显示每个任务的执行情况和总体统计信息
func (s *Scheduler) PrintReport() {
	// 打印报告标题
	fmt.Println("\n===== 任务执行报告 =====")

	// 初始化统计变量
	var totalDuration time.Duration // 总执行时长
	successCount := 0               // 成功任务计数
	failCount := 0                  // 失败任务计数

	// 遍历所有结果
	for i, result := range s.results {
		// 根据成功标志设置状态文本
		if !result.Success {
			// 增加失败计数
			failCount++
			// 打印失败任务的详细信息（包括错误信息）
			fmt.Printf("%d. %s: %s (错误: %v)\n", i+1, result.TaskName, result.Duration, result.Error)
		} else {
			// 增加成功计数
			successCount++
			// 打印成功任务的执行时长
			fmt.Printf("%d. %s: %s\n", i+1, result.TaskName, result.Duration)
		}
		// 累加总执行时长
		totalDuration += result.Duration
	}

	// 计算平均执行时长
	avgDuration := totalDuration / time.Duration(len(s.results))
	// 打印统计摘要
	fmt.Printf("\n统计: 成功 %d, 失败 %d, 总计 %d\n", successCount, failCount, len(s.results))
	// 打印时间统计
	fmt.Printf("总耗时: %v, 平均耗时: %v\n", totalDuration, avgDuration)
}

// GetResults 获取所有任务结果
// 返回所有任务的执行结果切片
func (s *Scheduler) GetResults() []TaskResult {
	return s.results
}

// 示例任务实现

// SimpleTask 简单任务结构体
// 实现Task接口，模拟一个简单的延时任务
type SimpleTask struct {
	name string        // 任务名称
	dur  time.Duration // 任务模拟执行时长
}

// Execute 执行简单任务
// 模拟任务执行，只是简单地休眠指定时长
func (t *SimpleTask) Execute() error {
	time.Sleep(t.dur) // 模拟任务执行
	return nil        // 总是返回nil，表示执行成功
}

// Name 获取任务名称
func (t *SimpleTask) Name() string {
	return t.name
}

// ErrorTask 可能失败的任务结构体
// 实现Task接口，模拟一个可能随机失败的任务
type ErrorTask struct {
	name string        // 任务名称
	dur  time.Duration // 任务模拟执行时长
}

// Execute 执行可能失败的任务
// 模拟任务执行，有一定概率随机失败
func (t *ErrorTask) Execute() error {
	time.Sleep(t.dur) // 模拟任务执行
	// 模拟随机失败（约20%的概率）
	if time.Now().UnixNano()%5 == 0 {
		return fmt.Errorf("随机失败") // 返回错误表示执行失败
	}
	return nil // 返回nil表示执行成功
}

// Name 获取任务名称
func (t *ErrorTask) Name() string {
	return t.name
}

func TestTask() {
	// 创建调度器，设置并发度为3
	scheduler := NewScheduler(3)

	// 添加10个简单任务到调度器
	for i := 1; i <= 10; i++ {
		// 创建简单任务实例
		task := &SimpleTask{
			name: fmt.Sprintf("简单任务-%d", i),               // 设置任务名称
			dur:  time.Duration(100*i) * time.Millisecond, // 设置执行时长（递增）
		}
		// 将任务添加到调度器
		scheduler.AddTask(task)
	}

	// 添加5个可能失败的任务到调度器
	for i := 1; i <= 5; i++ {
		// 创建可能失败的任务实例
		task := &ErrorTask{
			name: fmt.Sprintf("可能失败的任务-%d", i),            // 设置任务名称
			dur:  time.Duration(150*i) * time.Millisecond, // 设置执行时长（递增）
		}
		// 将任务添加到调度器
		scheduler.AddTask(task)
	}

	// 记录开始时间
	startTime := time.Now()
	// 执行所有任务（阻塞直到所有任务完成）
	scheduler.Run()
	// 计算总运行时间
	totalTime := time.Since(startTime)

	// 打印详细的执行报告
	scheduler.PrintReport()
	// 打印调度器总运行时间（包括调度开销）
	fmt.Printf("调度器总运行时间: %v\n", totalTime)

	// 获取所有任务执行结果（可用于进一步处理或存储）
	results := scheduler.GetResults()
	_ = results // 使用空白标识符忽略未使用的变量（实际应用中可进一步处理）

}
