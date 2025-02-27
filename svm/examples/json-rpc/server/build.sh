#!/bin/bash
set -e

# 获取项目根目录
ROOT_DIR=$(pwd)
echo "项目根目录: $ROOT_DIR"

# 构建 Rust 库
echo "正在构建 Rust 库..."
cd "$ROOT_DIR"
cargo build --release -p json-rpc-server


# 创建 go 目录（如果不存在）
mkdir -p "$ROOT_DIR/go/binding"

# 编译 Go 程序
cd "$ROOT_DIR/go"
echo "正在编译 Go 程序..."
go build -o transaction_simulator main.go

echo "构建完成！运行 ./transaction_simulator 进行测试"