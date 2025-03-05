package main

import (
	"fmt"

	"github.com/barry/json-rpc-server/binding"
)

func main() {
	fmt.Println("Simulating transaction...")

	// Call the Rust implementation
	result := binding.SimulateTransaction()

	if result {
		fmt.Printf("%v", result)
	} else {
		fmt.Println("Transaction simulation failed!")
	}
}

// package main
//
// import (
// 	"fmt"
// 	"runtime"
// 	"sync"
// 	"sync/atomic"
// 	"time"
//
// 	"github.com/barry/json-rpc-server/binding"
// )
//
// // 定义结果结构
// type SimulationResult struct {
// 	ID        int
// 	Success   bool
// 	Timestamp time.Time
// }
//
// func main() {
// 	// 获取系统CPU核心数
// 	numCPU := runtime.NumCPU()
// 	fmt.Printf("系统CPU核心数: %d\n", numCPU)
//
// 	// 设置Go最大并发数
// 	runtime.GOMAXPROCS(numCPU)
//
// 	// 模拟请求总数
// 	totalRequests := 100000
// 	fmt.Printf("开始模拟 %d 个交易请求...\n", totalRequests)
//
// 	// 创建结果通道
// 	resultChan := make(chan SimulationResult, totalRequests)
//
// 	// 使用WaitGroup等待所有协程完成
// 	var wg sync.WaitGroup
//
// 	// 记录开始时间
// 	startTime := time.Now()
//
// 	// 创建工作池
// 	// 使用系统核心数作为工作池大小
// 	workerCount := numCPU
//
// 	// 任务通道
// 	taskChan := make(chan int, totalRequests)
//
// 	// 计数器用于记录完成的请求数
// 	var completedCount int32 = 0
//
// 	// 启动工作协程
// 	for i := 0; i < workerCount; i++ {
// 		wg.Add(1)
// 		go func(workerID int) {
// 			defer wg.Done()
//
// 			for taskID := range taskChan {
// 				// 调用Rust实现
// 				success := binding.SimulateTransaction()
//
// 				// 创建结果
// 				result := SimulationResult{
// 					ID:        taskID,
// 					Success:   success,
// 					Timestamp: time.Now(),
// 				}
//
// 				// 将结果发送到通道
// 				resultChan <- result
//
// 				// 更新计数器并打印进度
// 				completed := atomic.AddInt32(&completedCount, 1)
// 				if completed%1000 == 0 {
// 					fmt.Printf("已完成: %d / %d 请求 (%.2f%%)\n",
// 						completed, totalRequests, float64(completed)/float64(totalRequests)*100)
// 				}
// 			}
// 		}(i)
// 	}
//
// 	// 分发任务
// 	for i := 0; i < totalRequests; i++ {
// 		taskChan <- i
// 	}
// 	close(taskChan)
//
// 	// 等待所有工作协程完成
// 	go func() {
// 		wg.Wait()
// 		close(resultChan)
// 	}()
//
// 	// 统计结果
// 	var successCount int
// 	var failureCount int
//
// 	// 创建日志文件
// 	logFilePath := fmt.Sprintf("simulation_results_%s.log", time.Now().Format("20060102_150405"))
// 	// 这里可以添加文件写入逻辑，根据需要实现
//
// 	// 处理结果通道
// 	for result := range resultChan {
// 		if result.Success {
// 			successCount++
// 		} else {
// 			failureCount++
// 		}
//
// 		// 可以选择记录每个结果
// 		// 这里简化处理，不记录每个单独的结果，只在最后汇总
// 	}
//
// 	// 计算总耗时
// 	duration := time.Since(startTime)
//
// 	// 打印统计信息
// 	fmt.Println("\n========== 模拟结果 ==========")
// 	fmt.Printf("总请求数: %d\n", totalRequests)
// 	fmt.Printf("成功数: %d (%.2f%%)\n", successCount, float64(successCount)/float64(totalRequests)*100)
// 	fmt.Printf("失败数: %d (%.2f%%)\n", failureCount, float64(failureCount)/float64(totalRequests)*100)
// 	fmt.Printf("总耗时: %v\n", duration)
// 	fmt.Printf("每秒处理请求数 (QPS): %.2f\n", float64(totalRequests)/duration.Seconds())
// 	fmt.Printf("每个请求平均耗时: %.2f ms\n", duration.Seconds()*1000/float64(totalRequests))
// 	fmt.Println("================================")
// }