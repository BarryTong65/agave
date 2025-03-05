package main

/*
#cgo LDFLAGS: -L/Users/barry/binance/agave/programs/sbf/target/release -lsbf_bridge
#include <stdlib.h>
extern int call_test_create_vm();
*/
import "C"
import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("调用Rust test_create_vm函数...")

	result := C.call_test_create_vm()

	if result == 0 {
		fmt.Println("测试执行成功!")
	} else {
		fmt.Println("测试失败，错误代码:", result)
		os.Exit(1)
	}
}