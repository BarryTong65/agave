package binding

/*
#cgo LDFLAGS: -L/Users/barry/binance/agave/svm/examples/target/release -ljson_rpc_server
#cgo darwin,arm64 LDFLAGS: -framework Security -framework CoreFoundation
#include <stdbool.h>

extern bool simulate_transaction_c();
*/
import "C"

// SimulateTransaction 调用 Rust 实现的交易模拟函数
func SimulateTransaction() bool {
	return bool(C.simulate_transaction_c())
}
