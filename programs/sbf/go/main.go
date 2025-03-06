package main

/*
#cgo LDFLAGS: -L/Users/barry/binance/agave/programs/sbf/target/release -lsbf_bridge
#include <stdlib.h>
extern int call_test_create_vm();
extern int call_test_instruction_count_tuner();
*/
import "C"
import (
    "fmt"
    "os"
    "time"
)

func main() {
    fmt.Println("开始调用Rust test_create_vm函数...")

    // 记录开始时间
    startTime := time.Now()

    // 调用Rust函数
    result := C.call_test_instruction_count_tuner()

    // 计算耗时
    duration := time.Since(startTime)

    if result == 0 {
       fmt.Printf("测试执行成功! 耗时: %v\n", duration)
    } else {
       fmt.Printf("测试失败，错误代码: %d, 耗时: %v\n", result, duration)
       os.Exit(1)
    }
}


// package main
//
// /*
// #cgo LDFLAGS: -L/Users/barry/binance/agave/programs/sbf/target/release -lsbf_bridge
// #include <stdlib.h>
// extern int call_test_create_vm();
// */
// import "C"
// import (
// 	"fmt"
// 	"runtime"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// )
//
// type Request struct {
// 	ID        int
// 	Timestamp time.Time
// }
//
// type Response struct {
// 	RequestID   int
// 	ProcessedBy int // 处理该请求的协程ID
// 	Success     bool
// 	RequestTime time.Time
// 	ResponseTime time.Time
// 	Duration    time.Duration
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
// 	totalRequests := 1000000
// 	fmt.Printf("准备处理 %d 个请求...\n", totalRequests)
//
// 	// 创建请求和响应通道
// 	requestChan := make(chan Request, totalRequests)
// 	responseChan := make(chan Response, totalRequests)
//
// 	// 使用WaitGroup等待所有协程完成
// 	var wg sync.WaitGroup
//
// 	// 记录开始时间
// 	startTime := time.Now()
//
// 	// 创建工作池 - 每个工作协程只启动一个VM
// 	vmCount := numCPU
// 	fmt.Printf("启动 %d 个VM协程...\n", vmCount)
//
// 	// 记录VM启动时间
// 	vmStartTime := time.Now()
//
// 	// 启动VM工作协程
// 	for i := 0; i < vmCount; i++ {
// 		wg.Add(1)
// 		go func(workerID int) {
// 			defer wg.Done()
//
// 			// 每个协程只调用一次call_test_create_vm
// 			result := C.call_test_create_vm()
//
// 			if result != 0 {
// 				return
// 			}
//
// 			// 处理请求
// 			for request := range requestChan {
// // 			    time.Sleep(time.Microsecond * 3)
// 				// 记录请求开始处理的时间
// 				// 模拟请求处理 - 实际情况中这里可能需要调用VM的特定方法
// 				// 这里不再调用call_test_create_vm()，而是使用已创建的VM实例
//
// 				// 创建响应
// 				response := Response{
// 					RequestID:    request.ID,
// 					ProcessedBy:  workerID,
// 					Success:      true,
// 				}
//
// 				// 将响应发送到通道
// 				responseChan <- response
// 			}
// 		}(i)
// 	}
//
// 	// 记录VM启动总耗时
// 	vmStartDuration := time.Since(vmStartTime)
// 	fmt.Printf("所有VM启动完成，耗时: %v\n", vmStartDuration)
//
// 	// 计数器用于记录已处理的请求数
// 	var processedCount int32 = 0
//
// 	// 启动日志记录协程
// 	go func() {
// 		for _ = range responseChan {
// 			processed := atomic.AddInt32(&processedCount, 1)
//
// 			// 每1000个请求打印一次进度，或者根据需要调整
// // 			if processed%1000 == 0 || processed == 1 {
// // 				fmt.Printf("已处理: %d/%d 请求 (%.2f%%), 最近请求由协程 #%d 处理，耗时: %v\n",
// // 					processed, totalRequests, float64(processed)/float64(totalRequests)*100,
// // 					response.ProcessedBy, response.Duration)
// // 			}
//
// 			// 所有请求处理完毕
// 			if processed >= int32(totalRequests) {
// 				break
// 			}
// 		}
//
// 		// 关闭请求通道，通知所有工作协程可以退出
// 		close(requestChan)
// 	}()
//
// 	// 发送请求
// 	fmt.Println("开始发送请求...")
// 	requestStartTime := time.Now()
//
// 	for i := 0; i < totalRequests; i++ {
// 		requestChan <- Request{
// 			ID:        i,
// // 			Timestamp: time.Now(),
// 		}
// 	}
//
// 	// 等待所有响应被处理
// // 	for atomic.LoadInt32(&processedCount) < int32(totalRequests) {
// // 		time.Sleep(100 * time.Millisecond)
// // 	}
//
// 	// 计算总耗时
// 	totalDuration := time.Since(startTime)
// 	requestDuration := time.Since(requestStartTime)
//
// 	// 打印统计信息
// 	fmt.Println("\n========== 性能测试结果 ==========")
// 	fmt.Printf("VM数量: %d (基于CPU核心数)\n", vmCount)
// 	fmt.Printf("VM启动总耗时: %v\n", vmStartDuration)
// 	fmt.Printf("总请求数: %d\n", totalRequests)
// 	fmt.Printf("请求发送耗时: %v\n", requestDuration)
// 	fmt.Printf("总耗时(包含VM启动): %v\n", totalDuration)
// 	fmt.Printf("每秒处理请求数 (QPS): %.2f\n", float64(totalRequests)/requestDuration.Seconds())
// 	fmt.Printf("每个请求平均耗时: %.2f μs\n", float64(requestDuration.Microseconds())/float64(totalRequests))
// 	fmt.Println("===================================")
//
// 	fmt.Println("测试完成!")
// }