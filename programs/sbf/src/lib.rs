// lib.rs
#![cfg(feature = "sbf_c")]
#![allow(clippy::uninlined_format_args)]
#![allow(clippy::arithmetic_side_effects)]
#![cfg_attr(
    any(target_os = "windows", not(target_arch = "x86_64")),
    allow(dead_code, unused_imports)
)]

use std::os::raw::c_int;
use std::fs::File;
use std::io::Read;
use std::path::PathBuf;
use solana_sbpf::program::{BuiltinProgram, SBPFVersion};
use {
    solana_feature_set::bpf_account_data_direct_mapping, solana_sbpf::memory_region::MemoryState,
    solana_sdk::signer::keypair::Keypair, std::slice,
};

use {
    byteorder::{ByteOrder, LittleEndian, WriteBytesExt},
    solana_bpf_loader_program::{
        create_vm, serialization::serialize_parameters,
        syscalls::create_program_runtime_environment_v1,
    },
    solana_compute_budget::compute_budget::ComputeBudget,
    solana_feature_set::FeatureSet,
    solana_measure::measure::Measure,
    solana_program_runtime::invoke_context::InvokeContext,
    solana_runtime::{
        bank::Bank,
        bank_client::BankClient,
        genesis_utils::{create_genesis_config, GenesisConfigInfo},
        loader_utils::{load_program_from_file, load_program_of_loader_v4},
    },
    solana_sbpf::{
        ebpf::MM_INPUT_START, elf::Executable, memory_region::MemoryRegion,
        verifier::RequisiteVerifier, vm::ContextObject,
    },
    solana_sdk::{
        account::AccountSharedData,
        bpf_loader,
        client::SyncClient,
        entrypoint::SUCCESS,
        instruction::{AccountMeta, Instruction},
        message::Message,
        native_loader,
        pubkey::Pubkey,
        signature::Signer,
        transaction_context::InstructionAccount,
    },
    std::{mem, sync::Arc},
};

const ARMSTRONG_LIMIT: u64 = 500;
const ARMSTRONG_EXPECTED: u64 = 5;

macro_rules! with_mock_invoke_context {
    ($invoke_context:ident, $loader_id:expr, $account_size:expr) => {
        let program_key = Pubkey::new_unique();
        let transaction_accounts = vec![
            (
                $loader_id,
                AccountSharedData::new(0, 0, &native_loader::id()),
            ),
            (program_key, AccountSharedData::new(1, 0, &$loader_id)),
            (
                Pubkey::new_unique(),
                AccountSharedData::new(2, $account_size, &program_key),
            ),
        ];
        let instruction_accounts = vec![InstructionAccount {
            index_in_transaction: 2,
            index_in_caller: 2,
            index_in_callee: 0,
            is_signer: false,
            is_writable: true,
        }];
        solana_program_runtime::with_mock_invoke_context!(
            $invoke_context,
            transaction_context,
            transaction_accounts
        );
        $invoke_context
            .transaction_context
            .get_next_instruction_context()
            .unwrap()
            .configure(&[0, 1], &instruction_accounts, &[]);
        $invoke_context.push().unwrap();
    };
}

fn clone_regions(regions: &[MemoryRegion]) -> Vec<MemoryRegion> {
    unsafe {
        regions
            .iter()
            .map(|region| match region.state.get() {
                MemoryState::Readable => MemoryRegion::new_readonly(
                    slice::from_raw_parts(region.host_addr.get() as *const _, region.len as usize),
                    region.vm_addr,
                ),
                MemoryState::Writable => MemoryRegion::new_writable(
                    slice::from_raw_parts_mut(
                        region.host_addr.get() as *mut _,
                        region.len as usize,
                    ),
                    region.vm_addr,
                ),
                MemoryState::Cow(id) => MemoryRegion::new_cow(
                    slice::from_raw_parts(region.host_addr.get() as *const _, region.len as usize),
                    region.vm_addr,
                    id,
                ),
            })
            .collect()
    }
}


#[cfg(feature = "sbf_c")]
pub fn test_create_vm_wrapper() -> Result<(), String> {
    use std::time::Instant;
    // We wrap the test in a function that returns a Result
    match std::panic::catch_unwind(|| {
        let total_start = Instant::now();

        let load_start = Instant::now();
        let elf = load_program_from_file("tuner");
        let load_duration = load_start.elapsed();
        println!("Step 1: Load program from file - {:?}", load_duration);

        let context_start = Instant::now();
        with_mock_invoke_context!(invoke_context, bpf_loader::id(), 10000001);
        const BUDGET: u64 = 200_000;
        invoke_context.mock_set_remaining(BUDGET);
        let context_duration = context_start.elapsed();
        println!("Step 2: Setup invoke context - {:?}", context_duration);

        let feature_start = Instant::now();
        let direct_mapping = invoke_context
            .get_feature_set()
            .is_active(&bpf_account_data_direct_mapping::id());
        let feature_duration = feature_start.elapsed();
        println!("Step 3: Get feature set - {:?}", feature_duration);

        let env_start = Instant::now();
        let program_runtime_environment = create_program_runtime_environment_v1(
            invoke_context.get_feature_set(),
            &ComputeBudget::default(),
            true,
            false,
        );
        let env_duration = env_start.elapsed();
        println!("Step 4: Create runtime environment - {:?}", env_duration);

        let executable_start = Instant::now();
        let executable =
            Executable::<InvokeContext>::from_elf(&elf, Arc::new(program_runtime_environment.unwrap()))
                .unwrap();
        let executable_duration = executable_start.elapsed();
        println!("Step 5: Create executable from ELF - {:?}", executable_duration);

        let verify_start = Instant::now();
        executable.verify::<RequisiteVerifier>().unwrap();
        let verify_duration = verify_start.elapsed();
        println!("Step 6: Verify executable - {:?}", verify_duration);

        // Serialize account data
        let serialize_start = Instant::now();
        let (_serialized, regions, account_lengths) = serialize_parameters(
            invoke_context.transaction_context,
            invoke_context
                .transaction_context
                .get_current_instruction_context()
                .unwrap(),
            !direct_mapping, // copy_account_data,
        )
            .unwrap();
        let serialize_duration = serialize_start.elapsed();
        println!("Step 7: Serialize parameters - {:?}", serialize_duration);

        let vm_start = Instant::now();
        create_vm!(
            vm,
            &executable,
            clone_regions(&regions),
            account_lengths.clone(),
            &mut invoke_context,
        );
        let vm_creation_duration = vm_start.elapsed();
        println!("Step 8: Create VM - {:?}", vm_creation_duration);

        let unwrap_start = Instant::now();
        vm.unwrap();
        let unwrap_duration = unwrap_start.elapsed();
        println!("Step 9: Unwrap VM - {:?}", unwrap_duration);

        let total_duration = total_start.elapsed();
        println!("Total execution time: {:?}", total_duration);
    }) {
        Ok(_) => {
            println!("Test completed successfully");
            Ok(())
        },
        Err(e) => {
            println!("Test failed with error: {:?}", e);
            Err("test_create_vm panicked".to_string())
        },
    }
}

#[no_mangle]
pub extern "C" fn call_test_create_vm() -> c_int {
    // Call the function and return a status code
    match test_create_vm_wrapper() {
        Ok(_) => 0,
        Err(_) => 1,
    }
}

#[cfg(feature = "sbf_c")]
pub fn test_instruction_count_tuner_wrapper() -> Result<(), String> {
    use std::time::Instant;
    // We wrap the test in a function that returns a Result
    match std::panic::catch_unwind(|| {
        println!("Program test_instruction_count_tuner");
        //let elf = load_program_from_file("tuner");
        println!("Program read_to_end");
        with_mock_invoke_context!(invoke_context, bpf_loader::id(), 10000001);
        const BUDGET: u64 = 900_000;
        invoke_context.mock_set_remaining(BUDGET);

        let direct_mapping = invoke_context
            .get_feature_set()
            .is_active(&bpf_account_data_direct_mapping::id());

        // Serialize account data
        let (_serialized, regions, account_lengths) = serialize_parameters(
            invoke_context.transaction_context,
            invoke_context
                .transaction_context
                .get_current_instruction_context()
                .unwrap(),
            !direct_mapping, // copy_account_data
        )
            .unwrap();

        let program_runtime_environment = create_program_runtime_environment_v1(
            invoke_context.get_feature_set(),
            &ComputeBudget::default(),
            true,
            false,
        );
        // let executable =
        //     Executable::<InvokeContext>::from_elf(&elf, Arc::new(program_runtime_environment.unwrap()))
        //         .unwrap();
        let loader = solana_type_overrides::sync::Arc::new(program_runtime_environment.unwrap());
        let function_registry = solana_sbpf::program::FunctionRegistry::default();
        let executable = solana_sbpf::elf::Executable::<InvokeContext>::from_text_bytes(
            &[0x9D, 0, 0, 0, 0, 0, 0, 0],
            loader,
            SBPFVersion::V3,
            function_registry,
        )
            .unwrap();

        executable.verify::<RequisiteVerifier>().unwrap();

        create_vm!(
        vm,
        &executable,
        regions,
        account_lengths,
        &mut invoke_context,
    );
        let (mut vm, _, _) = vm.unwrap();

        println!("create_vm");

        let mut measure = Measure::start("tune");
        let (instructions, _result) = vm.execute_program(&executable, true);
        println!("Program executed with result: {:?}", _result);
        measure.stop();

        // assert_eq!(
        //     0,
        //     vm.context_object_pointer.get_remaining(),
        //     "Tuner must consume the whole budget"
        // );
        println!(
            "{:?} compute units took {:?} us ({:?} instructions)",
            BUDGET - vm.context_object_pointer.get_remaining(),
            measure.as_us(),
            instructions,
        );
        println!("Finished bench_instruction_count_tuner test");
    }) {
        Ok(_) => {
            println!("Test completed successfully");
            Ok(())
        },
        Err(e) => {
            println!("Test failed with error: {:?}", e);
            Err("test_create_vm panicked".to_string())
        },
    }
}

#[no_mangle]
pub extern "C" fn call_test_instruction_count_tuner() -> c_int {
    // Call the function and return a status code
    match test_instruction_count_tuner_wrapper() {
        Ok(_) => 0,
        Err(_) => 1,
    }
}