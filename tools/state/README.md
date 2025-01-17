# state tool

This tool allows to rerun a set of batches, you could set a flag to persist changes in hashDB

# Usage



## Network configuration
If you want to avoid passing network configuration (`--network` and `--custom-network-file`) you need to provide the L2ChainID (`--l2_chain_id`)

## Reprocess a set of batches and compare with state database
This re-execute a batch/batches and check if match the data on DB.
It have some flags to allow:
- `--write_on_hash_db`: this for each execution create the corresponding MT if possible
it override state_db
- `--first_batch`: first batch to process (default: 1)
- `--last_batch`: last batch to process (default: the highest batch on batch table)
- `--l2_chain_id`:  Instead of asking to SMC you can set it 
- `--dont_stop_on_error`: If a batch have an error the process doesn't stop
- `--prefer_execution_state_root`: The oldStateRoot used to process a batch is usually is the stateRoot of the previous batch on database but, with this flag, you could use the calculated stateRoot from the execution result from previous batch instead

To see the full flags execute:
```
go run ./tools/state/. reprocess 
```

# Examples:

- You need to set the right `State` config, `Executor` config and `MTClient`. You can override the parameters with environment variables: 
```
KEVM_NODE_MTCLIENT_URI="127.0.0.1:50061" ZKEVM_NODE_STATE_DB_HOST="127.0.0.1" ZKEVM_NODE_EXECUTOR_URI="127.0.0.1:50071" go run ./tools/state/. reprocess -cfg test/config/test.node.config.toml   -l2_chain_id 1440 --last_batch_number 5000
```
- We are setting the `chain_id` directly so we don't need the genesis data.
- All this examples redirect the log info to `/dev/null` for that reason if the command returns an error (`$? -ne 1`) relaunch without the redirection part (`2> /dev/null`) to see the full output

### Rebuild hashdb entries for first 5000 batches

```
go run ./tools/state/. reprocess -cfg test/config/test.node.config.toml   -l2_chain_id 1440 --last_batch_number 5000 --write_on_hash_db 2> /dev/null
```
expected output: 
```
         batch     91 1.80%: ... ntx:   1 WRITE (flush: 5955)  ETA:       38s speed:127.8 batch/s  StateRoot:0x9f2db3f7775f30f1e79b4c0d876b8094a839cdba2cc51a48359b817a1c07e09f [OK]
         batch     92 1.82%: ... ntx:   0 WRITE (flush: 5956)  ETA:       44s speed:112.7 batch/s  StateRoot:0x9f2db3f7775f30f1e79b4c0d876b8094a839cdba2cc51a48359b817a1c07e09f [OK]
         batch     93 1.84%: ... ntx:  11 WRITE (flush: 5957)  ETA:       49s speed:99.3 batch/s  StateRoot:0xf77f6df21cbb5455ebae4dd9275bf5753f6e7e94250afe537192e624b7291854 [OK]
         batch     94 1.86%: ... ntx:   0 WRITE (flush: 5958)  ETA:       54s speed:90.4 batch/s  StateRoot:0xf77f6df21cbb5455ebae4dd9275bf5753f6e7e94250afe537192e624b7291854 [OK]
         batch     95 1.88%: ... ntx:   0 WRITE (flush: 5959)  ETA:       54s speed:91.2 batch/s  StateRoot:0xf77f6df21cbb5455ebae4dd9275bf5753f6e7e94250afe537192e624b7291854 [OK]
         batch     96 1.90%: ... ntx:   0 WRITE (flush: 5960)  ETA:       53s speed:92.0 batch/s  StateRoot:0xf77f6df21cbb5455ebae4dd9275bf5753f6e7e94250afe537192e624b7291854 [OK]
         batch     97 1.92%: ... ntx:   0 WRITE (flush: 5961)  ETA:       53s speed:92.7 batch/s  StateRoot:0xf77f6df21cbb5455ebae4dd9275bf5753f6e7e94250afe537192e624b7291854 [OK]
```

### Check that the batches from 1000 to  5000 match stateRoot
```
 go run ./tools/state/. reprocess -cfg test/config/test.node.config.toml -l2_chain_id 1440 --first_batch_number 1000 --last_batch_number 5000  2> /dev/null
```
