//go:build bitcoind && rpcpolling
// +build bitcoind,rpcpolling

package lntest

import (
	"github.com/brsuite/brond/chaincfg"
)

// NewBackend starts a bitcoind node without the txindex enabled and returns a
// BitoindBackendConfig for that node.
func NewBackend(miner string, netParams *chaincfg.Params) (
	*BitcoindBackendConfig, func() error, error) {

	extraArgs := []string{
		"-debug",
		"-regtest",
		"-txindex",
		"-disablewallet",
	}

	return newBackend(miner, netParams, extraArgs, true)
}
