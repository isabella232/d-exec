// Package unikernel_net_fs_ec implements an execution environment that uses the
// network and filesystem to perform s*G on ED25519.
package unikernel_net_fs_ec

import (
	"fmt"
	"net"
	"os"
	"time"

	"encoding/hex"

	"go.dedis.ch/dela"
	"go.dedis.ch/dela/core/execution"
	"go.dedis.ch/dela/core/store"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/xerrors"
)

const defaultUnikernelAddr = "172.44.0.2:1024"
const defaultFilePath = "/home/nkcr/d-exec/unikernel/apps/simple_crypto_network_fs/mnt/ec_multiply"

var storeKey = [32]byte{0, 0, 10}
var resultKey = [32]byte{0, 0, 40}

var suite = suites.MustFind("Ed25519")

const addrArg = "tcp:addr"
const dialTimeout = time.Second * 1

// Service ...
type Service struct {
}

// NewExecution ...
func NewExecution() *Service {
	return &Service{}
}

// Execute uses a unikernel to perform s*G. A random scalar is hex encoded and
// save in the relevant file that the unikernel use. A TCP command is then sent
// to call the unikernel. The result is a point returned by the unikernel.
func (hs *Service) Execute(snap store.Snapshot, step execution.Step) (execution.Result, error) {
	res := execution.Result{}

	current, err := snap.Get(storeKey[:])
	if err != nil {
		return res, xerrors.Errorf("failed to get store value: %v", err)
	}

	if len(current) == 0 {
		current = make([]byte, 8)
	}

	addr := string(step.Current.GetArg(addrArg))
	if addr == "" {
		addr = os.Getenv("UNIKERNEL_TCP")
		if addr == "" {
			addr = defaultUnikernelAddr
		}
	}

	hexBuf := hex.EncodeToString(current)
	err = os.WriteFile(defaultFilePath, []byte(fmt.Sprintf("ec_multiply\n%s", hexBuf)), 0644)
	if err != nil {
		return res, xerrors.Errorf("failed to write scalar: %v", err)
	}

	conn, err := net.DialTimeout("tcp", string(addr), dialTimeout)
	if err != nil {
		return res, xerrors.Errorf("failed to connect to tcp with %s: %v", addr, err)
	}

	dela.Logger.Info().Msgf("sending value command: ec_multiply")

	_, err = conn.Write([]byte("ec_multiply"))
	if err != nil {
		return res, xerrors.Errorf("failed to send to unikernel: %v", err)
	}

	readRes := make([]byte, 64)

	conn.SetReadDeadline(time.Now().Add(time.Second * 100))

	_, err = conn.Read(readRes)
	if err != nil {
		return res, xerrors.Errorf("failed to read result: %v", err)
	}

	snap.Set(resultKey[:], readRes)

	// resultBuf, err := hex.DecodeString(string(readRes))
	// if err != nil {
	// 	return res, xerrors.Errorf("failed to read result: %v", err)
	// }

	// err = suite.Point().UnmarshalBinary(resultBuf)
	// if err != nil {
	// 	return res, xerrors.Errorf("failed to unmarshal point: %v", err)
	// }

	dela.Logger.Info().Msgf("got result: %s", readRes)

	return execution.Result{
		Accepted: true,
	}, nil
}
