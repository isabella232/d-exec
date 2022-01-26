package main

import (
	"encoding/binary"
	"testing"

	"github.com/dedis/d-exec/goland/evm"
	"github.com/dedis/d-exec/goland/tcp"
	"github.com/dedis/d-exec/goland/tcp_ec"
	"github.com/dedis/d-exec/goland/wasm"
	"github.com/dedis/d-exec/goland/unikernel_net_fs_ec"

	"github.com/stretchr/testify/require"
	"go.dedis.ch/dela/core/execution"
	"go.dedis.ch/dela/core/store"
	"go.dedis.ch/dela/core/txn"
	"go.dedis.ch/kyber/v3/suites"

	"encoding/base64"
	"encoding/json"
)

var storeKey = [32]byte{0, 0, 10}

const iterations = 1

var suite = suites.MustFind("Ed25519")

// Increment benchmark

func BenchmarkNative_Increment(b *testing.B) {
	for i := 0; i < iterations; i++ {
		k := 0
		k++
	}
}

func BenchmarkGraalvmTCP_Increment(b *testing.B) {
	testWithAddr(b, "127.0.0.1:12347")
}

func BenchmarkLocalTCP_Increment(b *testing.B) {
	testWithAddr(b, "127.0.0.1:12346")
}

func BenchmarkUnikernelTCP_Increment(b *testing.B) {
	testWithAddr(b, "172.44.0.2:12345")
}

func BenchmarkEVMLocal_Increment(b *testing.B) {
	n := iterations

	storage := newInmemory()
	step := execution.Step{Previous: []txn.Transaction{}, Current: tx{
		args: map[string][]byte{"contractName": []byte("increment")},
	}}

	exec, err := evm.NewExecution("increment")
	require.NoError(b, err)

	for i := 0; i < n; i++ {
		_, err := exec.Execute(storage, step)
		if err != nil {
			b.Logf("failed to execute: %+v", err)
			b.FailNow()
		}
	}
}

func BenchmarkEVMTCP_Increment(b *testing.B) {
	testWithAddr(b, "127.0.0.1:12347")
}

// Simple crypto (Elliptic curve - EC) benchmarks

func BenchmarkNative_EC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		scalar := suite.Scalar().Pick(suite.RandomStream())
		_, err := scalar.MarshalBinary()
		require.NoError(b, err)
		for k := 0; k < 1e6; k++ {

			point := suite.Point().Mul(scalar, nil)
			_, err = point.MarshalBinary()
			require.NoError(b, err)
		}
	}
}

func BenchmarkEVMLocal_EC(b *testing.B) {
	storage := newInmemory()
	step := execution.Step{Previous: []txn.Transaction{}, Current: tx{
		args: map[string][]byte{"contractName": []byte("Ed25519")},
	}}

	exec, err := evm.NewExecution("Ed25519")
	require.NoError(b, err)

	storeKey := [32]byte{0, 0, 10}
	gasUsageKey := [32]byte{0, 0, 20}
	runCountKey := [32]byte{0, 0, 30}
	//resultKey := [32]byte{0, 0, 40}

	storage.Set(gasUsageKey[:], make([]byte, 8))
	storage.Set(runCountKey[:], make([]byte, 8))

	for i := 0; i < b.N; i++ {
		scalar := suite.Scalar().Pick(suite.RandomStream())

		scalarBuf, err := scalar.MarshalBinary()
		require.NoError(b, err)

		storage.Set(storeKey[:], scalarBuf)
		_, err = exec.Execute(storage, step)
		if err != nil {
			b.Logf("failed to execute: %+v", err)
			b.FailNow()
		}
	}

	//	gasUsageBuf, err := storage.Get(gasUsageKey[:])
	//	require.NoError(b, err)

	//	gasUsage := float64(binary.LittleEndian.Uint64(gasUsageBuf))

	//	runCountBuf, err := storage.Get(runCountKey[:])
	//	require.NoError(b, err)

	//	runCount := float64(binary.LittleEndian.Uint64(runCountBuf))

	//	fmt.Printf("Did %f multiplications. Average Gas Usage=%.2f\n", runCount, gasUsage/runCount)

}

func BenchmarkTCP_Simple_EC(b *testing.B) {
	addr := "127.0.0.1:12346"
	storeKey := [32]byte{0, 0, 10}

	storage := newInmemory()
	step := execution.Step{Previous: []txn.Transaction{}, Current: tx{
		args: map[string][]byte{"tcp:addr": []byte(addr)},
	}}
	exec := tcp_ec.NewExecution()

	for i := 0; i < iterations; i++ {
		scalar := suite.Scalar().Pick(suite.RandomStream())

		scalarBuf, err := scalar.MarshalBinary()
		require.NoError(b, err)

		storage.Set(storeKey[:], scalarBuf)

		_, err = exec.Execute(storage, step)
		if err != nil {
			b.Logf("failed to execute; %v", err)
			b.FailNow()
		}
	}
}

func BenchmarkUnikernel_Network_FS_Simple_EC(b *testing.B) {
	addr := "172.44.0.2:1024"
	storeKey := [32]byte{0, 0, 10}

	storage := newInmemory()
	step := execution.Step{Previous: []txn.Transaction{}, Current: tx{
		args: map[string][]byte{"tcp:addr": []byte(addr)},
	}}
	exec := unikernel_net_fs_ec.NewExecution()

	for i := 0; i < b.N; i++ {
		scalar := suite.Scalar().Pick(suite.RandomStream())

		scalarBuf, err := scalar.MarshalBinary()
		require.NoError(b, err)

		storage.Set(storeKey[:], scalarBuf)

		_, err = exec.Execute(storage, step)
		if err != nil {
			b.Logf("failed to execute; %v", err)
			b.FailNow()
		}
	}
}

func BenchmarkWASM_Go_EC(b *testing.B) {
	srvc := wasm.WASMService{}
	for i := 0; i < iterations; i++ {
		scalar := suite.Scalar().Pick(suite.RandomStream())
		scalarB, _ := scalar.MarshalBinary()
		args := map[string]interface{}{
			"scalar":           base64.StdEncoding.EncodeToString(scalarB),
			"contractName":     "simpleEC",
			"contractLanguage": "go",
		}
		marsh, err := json.Marshal(args)
		if err != nil {
			b.Error(err)
		}
		step := execution.Step{}
		step.Current = fakeTx{json: marsh}

		_, err = srvc.Execute(nil, step)
		if err != nil {
			b.Logf("failed to execute: %+v", err)
			b.FailNow()
		}
	}
}

func BenchmarkWASM_C_EC(b *testing.B) {
	var suite = suites.MustFind("Ed25519")
	step := execution.Step{}

	srvc := wasm.WASMService{}
	for i := 0; i < b.N; i++ {
		point1 := suite.Point().Pick(suite.RandomStream())
		scalar := suite.Scalar().Pick(suite.RandomStream())
		point1B, _ := point1.MarshalBinary()
		scalarB, _ := scalar.MarshalBinary()
		// encoding to base64 because JSON does not support raw bytes
		args := map[string]interface{}{
			"point1":           base64.StdEncoding.EncodeToString(point1B),
			"scalar":           base64.StdEncoding.EncodeToString(scalarB),
			"contractName":     "ed25519_gen_mul",
			"contractLanguage": "c",
		}
		marsh, err := json.Marshal(args)
		if err != nil {
			b.Error(err)
		}
		step.Current = fakeTx{json: marsh}
		_, err = srvc.Execute(nil, step)
	}
}

// -----------------------------------------------------------------------------
// Utility functions

type fakeExec struct {
	err error
}

type fakeTx struct {
	txn.Transaction
	json []byte
}

func (tx fakeTx) GetArg(key string) []byte {
	return []byte(tx.json)
}

func testWithAddr(b *testing.B, addr string) {
	n := iterations

	storage := newInmemory()
	step := execution.Step{Previous: []txn.Transaction{}, Current: tx{
		args: map[string][]byte{"tcp:addr": []byte(addr)},
	}}
	exec := tcp.NewExecution()

	initialCounter := uint64(1234)

	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, initialCounter)
	storage.Set(storeKey[:], buffer)

	for i := 0; i < n; i++ {
		_, err := exec.Execute(storage, step)
		if err != nil {
			b.Logf("failed to execute; %v", err)
			b.FailNow()
		}
	}
}

type inmemory struct {
	store.Readable
	store.Writable

	data map[string][]byte
}

func newInmemory() inmemory {
	return inmemory{
		data: make(map[string][]byte),
	}
}

func (i inmemory) Get(key []byte) ([]byte, error) {
	return i.data[string(key)], nil
}

func (i inmemory) Set(key []byte, value []byte) error {
	i.data[string(key)] = value
	return nil
}

func (i inmemory) Delete(key []byte) error {
	delete(i.data, string(key))
	return nil
}

type tx struct {
	txn.Transaction
	args map[string][]byte
}

func (t tx) GetArg(key string) []byte {
	return t.args[key]
}
