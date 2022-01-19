// Package main implements a simple TCP server that computes s*G on ED25519.
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"

	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/xerrors"
)

var suite = suites.MustFind("Ed25519")

func main() {
	port := flag.String("port", "12346", "the port")

	flag.Parse()

	addr := "127.0.0.1:" + *port
	fmt.Println("listening on", addr)

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic("failed to create addr: " + err.Error())
	}

	conn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("failed to listen: " + err.Error())
	}

	for {
		tcpCon, err := conn.Accept()
		if err != nil {
			panic("failed to accept: " + err.Error())
		}

		fmt.Println("connected to", tcpCon.RemoteAddr().String())

		res := make([]byte, 64)
		_, err = tcpCon.Read(res)
		if err != nil {
			panic("failed to read: " + err.Error())
		}

		resultBuf, err := hex.DecodeString(string(res))
		if err != nil {
			panic(xerrors.Errorf("failed to read result: %v", err).Error())
		}

		s := suite.Scalar()

		err = s.UnmarshalBinary(resultBuf)
		if err != nil {
			panic(xerrors.Errorf("failed to unmarshal scalar: %v", err).Error())
		}

		fmt.Println("received", s.String())

		r := suite.Point().Mul(s, nil)

		output, err := r.MarshalBinary()
		if err != nil {
			panic(xerrors.Errorf("failed to marshal result: %v", err).Error())
		}

		outputHex := hex.EncodeToString(output)

		fmt.Println("send back:", outputHex)

		_, err = tcpCon.Write([]byte(outputHex))
		if err != nil {
			panic("failed to write back: " + err.Error())
		}

		err = tcpCon.Close()
		if err != nil {
			panic("failed to close: " + err.Error())
		}
	}
}
