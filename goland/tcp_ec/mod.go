// Package tcp_ec contains a TCP server implementation that receives a
// hex-encoded ED25519 scalar, and compute s*G. It returns the point
// hex-encoded.

package tcp_ec

import (
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

const defaultTCPAddr = "127.0.0.1:1024"

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
			addr = defaultTCPAddr
		}
	}

	hexBuf := hex.EncodeToString(current)

	conn, err := net.DialTimeout("tcp", string(addr), dialTimeout)
	if err != nil {
		return res, xerrors.Errorf("failed to connect to tcp with %s: %v", addr, err)
	}

	dela.Logger.Info().Msgf("sending scalar value: %s", hexBuf)

	_, err = conn.Write([]byte(hexBuf))
	if err != nil {
		return res, xerrors.Errorf("failed to send to tcp server: %v", err)
	}

	readRes := make([]byte, 64)

	conn.SetReadDeadline(time.Now().Add(time.Second * 10))

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
