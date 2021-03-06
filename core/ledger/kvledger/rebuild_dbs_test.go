/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"os"
	"path/filepath"
	"testing"

	configtxtest "github.com/hyperledger/fabric/common/configtx/test"
	"github.com/hyperledger/fabric/core/ledger/mock"
	"github.com/stretchr/testify/require"
)

func TestRebuildDBs(t *testing.T) {
	conf, cleanup := testConfig(t)
	defer cleanup()
	provider := testutilNewProvider(conf, t, &mock.DeployedChaincodeInfoProvider{})

	numLedgers := 3
	for i := 0; i < numLedgers; i++ {
		genesisBlock, _ := configtxtest.MakeGenesisBlock(constructTestLedgerID(i))
		_, err := provider.CreateFromGenesisBlock(genesisBlock)
		require.NoError(t, err)
	}

	// rebuild should fail when provider is still open
	err := RebuildDBs(conf)
	require.Error(t, err, "as another peer node command is executing, wait for that command to complete its execution or terminate it before retrying")
	provider.Close()

	err = RebuildDBs(conf)
	require.NoError(t, err)

	// verify blockstoreIndex, configHistory, history, state, bookkeeper dbs are deleted
	rootFSPath := conf.RootFSPath
	_, err = os.Stat(filepath.Join(BlockStorePath(rootFSPath), "index"))
	require.Equal(t, os.IsNotExist(err), true)
	_, err = os.Stat(ConfigHistoryDBPath(rootFSPath))
	require.Equal(t, os.IsNotExist(err), true)
	_, err = os.Stat(HistoryDBPath(rootFSPath))
	require.Equal(t, os.IsNotExist(err), true)
	_, err = os.Stat(StateDBPath(rootFSPath))
	require.Equal(t, os.IsNotExist(err), true)
	_, err = os.Stat(BookkeeperDBPath(rootFSPath))
	require.Equal(t, os.IsNotExist(err), true)

	// rebuild again should be successful
	err = RebuildDBs(conf)
	require.NoError(t, err)
}
