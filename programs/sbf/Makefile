SBF_SDK_PATH := ./solana-release/bin/sdk/sbf
SRC_DIR := c/src
OUT_DIR := target/sbpf-solana-solana/release

test: rust all
	SBF_OUT_DIR=$(OUT_DIR) cargo test --features="sbf_rust,sbf_c" $(TEST_ARGS)

rust:
	./solana-release/bin/cargo-build-sbf --release --workspace

.PHONY: rust

include $(SBF_SDK_PATH)/c/sbf.mk
