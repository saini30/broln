//go:build !bitcoind && !neutrino
// +build !bitcoind,!neutrino

package lntest

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/brsuite/brond/btcjson"
	"github.com/brsuite/brond/chaincfg"
	"github.com/brsuite/brond/integration/rpctest"
	"github.com/brsuite/brond/rpcclient"
)

// logDirPattern is the pattern of the name of the temporary log directory.
const logDirPattern = "%s/.backendlogs"

// temp is used to signal we want to establish a temporary connection using the
// brond Node API.
//
// NOTE: Cannot be const, since the node API expects a reference.
var temp = "temp"

// brondBackendConfig is an implementation of the BackendConfig interface
// backed by a brond node.
type brondBackendConfig struct {
	// rpcConfig houses the connection config to the backing brond instance.
	rpcConfig rpcclient.ConnConfig

	// harness is the backing brond instance.
	harness *rpctest.Harness

	// minerAddr is the p2p address of the miner to connect to.
	minerAddr string
}

// A compile time assertion to ensure brondBackendConfig meets the BackendConfig
// interface.
var _ BackendConfig = (*brondBackendConfig)(nil)

// GenArgs returns the arguments needed to be passed to broln at startup for
// using this node as a chain backend.
func (b brondBackendConfig) GenArgs() []string {
	var args []string
	encodedCert := hex.EncodeToString(b.rpcConfig.Certificates)
	args = append(args, "--bitcoin.node=brond")
	args = append(args, fmt.Sprintf("--brond.rpchost=%v", b.rpcConfig.Host))
	args = append(args, fmt.Sprintf("--brond.rpcuser=%v", b.rpcConfig.User))
	args = append(args, fmt.Sprintf("--brond.rpcpass=%v", b.rpcConfig.Pass))
	args = append(args, fmt.Sprintf("--brond.rawrpccert=%v", encodedCert))

	return args
}

// ConnectMiner is called to establish a connection to the test miner.
func (b brondBackendConfig) ConnectMiner() error {
	return b.harness.Client.Node(btcjson.NConnect, b.minerAddr, &temp)
}

// DisconnectMiner is called to disconnect the miner.
func (b brondBackendConfig) DisconnectMiner() error {
	return b.harness.Client.Node(btcjson.NDisconnect, b.minerAddr, &temp)
}

// Credentials returns the rpc username, password and host for the backend.
func (b brondBackendConfig) Credentials() (string, string, string, error) {
	return b.rpcConfig.User, b.rpcConfig.Pass, b.rpcConfig.Host, nil
}

// Name returns the name of the backend type.
func (b brondBackendConfig) Name() string {
	return "brond"
}

// NewBackend starts a new rpctest.Harness and returns a brondBackendConfig for
// that node. miner should be set to the P2P address of the miner to connect
// to.
func NewBackend(miner string, netParams *chaincfg.Params) (
	*brondBackendConfig, func() error, error) {

	baseLogDir := fmt.Sprintf(logDirPattern, GetLogDir())
	args := []string{
		"--rejectnonstd",
		"--txindex",
		"--trickleinterval=100ms",
		"--debuglevel=debug",
		"--logdir=" + baseLogDir,
		"--nowinservice",
		// The miner will get banned and disconnected from the node if
		// its requested data are not found. We add a nobanning flag to
		// make sure they stay connected if it happens.
		"--nobanning",
		// Don't disconnect if a reply takes too long.
		"--nostalldetect",
	}
	chainBackend, err := rpctest.New(netParams, nil, args, GetbrondBinary())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create brond node: %v", err)
	}

	// We want to overwrite some of the connection settings to make the
	// tests more robust. We might need to restart the backend while there
	// are already blocks present, which will take a bit longer than the
	// 1 second the default settings amount to. Doubling both values will
	// give us retries up to 4 seconds.
	chainBackend.MaxConnRetries = rpctest.DefaultMaxConnectionRetries * 2
	chainBackend.ConnectionRetryTimeout = rpctest.DefaultConnectionRetryTimeout * 2

	if err := chainBackend.SetUp(false, 0); err != nil {
		return nil, nil, fmt.Errorf("unable to set up brond backend: %v", err)
	}

	bd := &brondBackendConfig{
		rpcConfig: chainBackend.RPCConfig(),
		harness:   chainBackend,
		minerAddr: miner,
	}

	cleanUp := func() error {
		var errStr string
		if err := chainBackend.TearDown(); err != nil {
			errStr += err.Error() + "\n"
		}

		// After shutting down the chain backend, we'll make a copy of
		// the log files, including any compressed log files from
		// logrorate, before deleting the temporary log dir.
		logDir := fmt.Sprintf("%s/%s", baseLogDir, netParams.Name)
		files, err := ioutil.ReadDir(logDir)
		if err != nil {
			errStr += fmt.Sprintf(
				"unable to read log directory: %v\n", err,
			)
		}

		for _, file := range files {
			logFile := fmt.Sprintf("%s/%s", logDir, file.Name())
			newFilename := strings.Replace(
				file.Name(), "brond.log", "output_brond_chainbackend.log", 1,
			)
			logDestination := fmt.Sprintf(
				"%s/%s", GetLogDir(), newFilename,
			)
			err := CopyFile(logDestination, logFile)
			if err != nil {
				errStr += fmt.Sprintf("unable to copy file: %v\n", err)
			}
		}

		if err = os.RemoveAll(baseLogDir); err != nil {
			errStr += fmt.Sprintf(
				"cannot remove dir %s: %v\n", baseLogDir, err,
			)
		}
		if errStr != "" {
			return errors.New(errStr)
		}
		return nil
	}

	return bd, cleanUp, nil
}
