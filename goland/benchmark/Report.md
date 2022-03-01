# Benchmark

This document is a technical report for internal purposes to keep track of our
progress. We describe two series of benchmarks that measure the efficiency of
different smart contract execution environment. The purpose of those benchmarks
is to get a comparison between a Unikernel-based solution and other common ones.

## Setups

The benchmark are run in go, using the native benchmark solution. The benchmark
reports the time by operation: ns/op.

Inside each benchmark we iterate the operations 50 times.

The code used for those results uses version
`d936178eccb2264109ecd8fc59fc9cfd8dfb79e0` of
https://github.com/dedis/dela/tree/unikernel-initial.

The experiments are run in the same virtual machine running `Ubuntu 20.04.1
LTS` with 4 GB of RAM:

```bash
$ lscpu
Architecture:                    x86_64
CPU op-mode(s):                  32-bit, 64-bit
Byte Order:                      Little Endian
Address sizes:                   45 bits physical, 48 bits virtual
CPU(s):                          2
On-line CPU(s) list:             0,1
Thread(s) per core:              1
Core(s) per socket:              1
Socket(s):                       2
NUMA node(s):                    1
Vendor ID:                       GenuineIntel
CPU family:                      6
Model:                           142
Model name:                      Intel(R) Core(TM) i7-7660U CPU @ 2.50GHz
Stepping:                        9
CPU MHz:                         2496.000
```

Results show the average among 5 runs.

Bellow we describe the 5 different setups.

### Native

This setup runs the smart contract code directly from the benchmark, in go
code. This simulates a "native" smart contract.

```
Bench
┌───────────────┐
│       Smart C.│
└───────────────┘
```

### Unikernel with TCP

This setups uses a Unikernel smart contract that is accessed via a TCP
connection.

```
Bench                Unikernel
┌─────────┐       ┌────────────┐
│         ├──────►│   Smart C. │
└─────────┘       └────────────┘
```

### Local TCP server

In this setup, the smart contract is run in go via a TCP connection.

```
Bench               go TCP server
┌──────────┐       ┌──────────┐
│          ├──────►│  Smart C.│
└──────────┘       └──────────┘
```

### Solidity in local

This setup uses the go-ethereum library to execute a solidity smart contract in
its VM from go.

```
Bench
┌───────────────┐
│               │
│         sol VM│
│     ┌─────────┤
│     │ Smart C.│
└─────┴─────────┘
 ```

### Solidity via TCP

This setup uses the go-ethereum library to execute a solidity smart contract via
a TCP server.

```
                    go TCP server
                   ┌────────────┐
Bench              │      Sol VM│
┌──────────┐       │  ┌─────────┤
│          ├──────►│  │ Smart C.│
└──────────┘       └──┴─────────┘
```

### WASM

This setup uses the Web-Assembly virtual machine.

Based on the work in https://github.com/dedis/student_21_dela-wasm version
`fc381375cb68a26d46cdbd070243e7244c2e54a1`, with the update that the point
selection is made inside the iteration loop.

```
Bench                WASM
┌─────────┐       ┌────────────┐
│         ├──────►│   Smart C. │
└─────────┘       └────────────┘
```

## Experiments

We perform two series of experiment:

1. **Increment**: in this series of experiment the smart contract increments
   the input it receives by 1.
2. **Simple crypto**: in this series of experiment the smart contract performs a
   simple operation on ed25519 elliptic curves. The smart contract takes a
   scalar (s) in argument. It then performs `s*G`, where `G` is the ed25519 base
   point.

## Results

|   [ns/op]    |Increment  |Simple crypto|
|--------------|-----------|-------------|
| Native       |0.00       |0.005        |
| Unikernel    |2342209082 |6072048716   |
| TCP          |0.014      |0.03646      |
| Solidity     |0.001      |0.237        |
| Solidity TCP |0.014      |             |
| WASM Go      |0.0238     |0.058        |
| WASM C       |0.0162     |0.052        |

# Repetitive simple crypto operation

In this series of experiments with execute the `s*G` operation a number of times
inside the execution environment: the unikernel, tcp server, EVM. The benchmark
runs the execution mutiple times, until having a stable result. The result is
based on an average over 6 runs.

## Result

| [ns/op]            | 1e0       | 1e1       | 1e2       | 1e3         | 1e4         | 1e5          | 1e6           |
|--------------------|-----------|-----------|-----------|-------------|-------------|--------------|---------------|
| Native             |     85'237|    796'121|  7'524'008|   84'807'264|  754'702'346| 7'688'306'495| 84'108'037'580|
| Unikernel (tcp+fs) |101'907'643|103'753'734|112'467'733|  146'313'097|  340'759'106| 2'314'039'681| 26'744'189'506|
| EVM (local)        |  3'917'656| 39'468'078|367'289'249|4'143'285'328| OutOfMemory |              |               |
| WASM C             |  1'617'968|  3'240'597| 16'243'341|  147'322'609|1'404'344'166|14'450'936'473|177'133'558'246|
## Procedures to setup and run the tests

### WASM

- edit `wasm_env/c/ed25519_gen_mul.c` and the loop condition at line 95:

```c
        for (int i = 0; i < 1e0; ++i)
        {
            crypto_scalarmult_ed25519_base(point, scalar);
        }
```

- compile again with the instructions from `wasm_env/c/README.md`.

- run the app, from `wasm_env`: `node app.js`

- run the benchmark, from this folder: `go test --bench BenchmarkWASM_C_EC`

### EVM

- edit `goland/evm/contracts/Ed25519.sol` and the loop condition at line 115:

```sol
for (uint i=0; i < 1e4; i++){
```

- compile the contract, from `goland/evm/contracts`: `./make_ed.sh`

- run the benchmark, from this folder: `go test --bench BenchmarkEVMLocal_EC`

### Unikernel

- edit `unikernel/apps/simple_crypto_network_js/main.c` and the `ITERATION`
  constant at line 42:

```c
#define ITERATIONS 1e7
```

- compile the unikernel, from `unikernel/apps/simple_crypto_network_fs`: `make`

- run the unikernel, from `unikernel/apps/simple_crypto_network_fs`: `./run`

- run the benchmark, from this folder: `go test --bench BenchmarkUnikernel_Network_FS_Simple_EC`

### Native

- run the benchmark, from this folder: `go test --bench BenchmarkNative_EC`