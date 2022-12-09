package backfill_test

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/phayes/freeport"
	. "github.com/stretchr/testify/assert"
	"github.com/synapsecns/sanguine/core/ginhelper"
	"github.com/synapsecns/sanguine/ethergo/backends"
	"github.com/synapsecns/sanguine/ethergo/backends/geth"
	"github.com/synapsecns/sanguine/services/omnirpc/config"
	omniHTTP "github.com/synapsecns/sanguine/services/omnirpc/http"
	"github.com/synapsecns/sanguine/services/omnirpc/proxy"
	"github.com/synapsecns/sanguine/services/scribe/backfill"
	"k8s.io/apimachinery/pkg/util/sets"
	"math/big"
	"net/http"
	"sync"
	"testing"
)

func TestMakeRange(t *testing.T) {
	testIntRange := backfill.MakeRange(0, 4)
	Equal(t, []int{0, 1, 2, 3, 4}, testIntRange)

	testUintRange := backfill.MakeRange(uint16(10), uint16(12))
	Equal(t, testUintRange, []uint16{10, 11, 12})
}

// startOmnirpcServer boots an omnirpc server for an rpc address.
// the url for this rpc is returned.
func (b *BackfillSuite) startOmnirpcServer(ctx context.Context, backend backends.SimulatedTestBackend) string {
	// run an omnirpc proxy to our backend
	server := proxy.NewProxy(config.Config{
		Chains: map[uint32]config.ChainConfig{
			uint32(backend.GetChainID()): {
				RPCs:   []string{backend.RPCAddress()},
				Checks: 1,
			},
		},
		Port:            uint16(freeport.GetPort()),
		RefreshInterval: 0,
		ClientType:      omniHTTP.FastHTTP.String(),
	}, omniHTTP.FastHTTP)

	go func() {
		server.Run(ctx)
	}()

	baseHost := fmt.Sprintf("http://0.0.0.0:%d", server.Port())
	healthCheck := fmt.Sprintf("%s%s", baseHost, ginhelper.HealthCheck)

	// wait for server to start
	b.Eventually(func() bool {
		select {
		case <-ctx.Done():
			b.T().Error(b.GetTestContext().Err())
		default:
			// see below
		}

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, healthCheck, nil)
		Nil(b.T(), err)

		res, err := http.DefaultClient.Do(request)
		if err == nil {
			defer func() {
				_ = res.Body.Close()
			}()
			return true
		}

		return false
	})

	return fmt.Sprintf("%s/rpc/%d", baseHost, backend.GetChainID())
}

// ReachBlockHeight reaches a block height on a backend.
func (b *BackfillSuite) ReachBlockHeight(ctx context.Context, backend backends.SimulatedTestBackend, desiredBlockHeight uint64) {
	i := 0
	for {
		select {
		case <-ctx.Done():
			b.T().Log(ctx.Err())
			return
		default:
			// continue
		}
		i++
		backend.FundAccount(ctx, common.BigToAddress(big.NewInt(int64(i))), *big.NewInt(params.Wei))

		latestBlock, err := backend.BlockNumber(ctx)
		Nil(b.T(), err)

		if latestBlock >= desiredBlockHeight {
			return
		}
	}
}

// ReachBlockHeight reaches a block height on a backend.
func (b *BackfillSuite) PopuluateWithLogs(ctx context.Context, backend backends.SimulatedTestBackend, desiredBlockHeight uint64) common.Address {
	i := 0
	var address common.Address
	for {
		select {
		case <-ctx.Done():
			b.T().Log(ctx.Err())
			return address
		default:
			// continue
		}
		i++
		backend.FundAccount(ctx, common.BigToAddress(big.NewInt(int64(i))), *big.NewInt(params.Wei))
		testContract, testRef := b.manager.GetTestContract(b.GetTestContext(), backend)
		address = testContract.Address()
		transactOpts := backend.GetTxContext(b.GetTestContext(), nil)
		tx, err := testRef.EmitEventA(transactOpts.TransactOpts, big.NewInt(1), big.NewInt(2), big.NewInt(3))
		Nil(b.T(), err)
		backend.WaitForConfirmation(b.GetTestContext(), tx)

		latestBlock, err := backend.BlockNumber(ctx)
		Nil(b.T(), err)

		if latestBlock >= desiredBlockHeight {
			return address
		}
	}
}

func (b *BackfillSuite) TestBlockTimesInRange() {
	testBackend := geth.NewEmbeddedBackend(b.GetTestContext(), b.T())

	// start an omnirpc proxy and run 10 test tranactions so we can batch call blocks
	//  1-10
	var wg sync.WaitGroup
	wg.Add(2)

	const desiredBlockHeight = 10

	go func() {
		defer wg.Done()
		b.ReachBlockHeight(b.GetTestContext(), testBackend, desiredBlockHeight)
	}()

	var host string
	go func() {
		defer wg.Done()
		host = b.startOmnirpcServer(b.GetTestContext(), testBackend)
	}()

	wg.Wait()

	scribeBackend, err := backfill.DialBackend(b.GetTestContext(), host)
	Nil(b.T(), err)

	res, err := backfill.BlockTimesInRange(b.GetTestContext(), scribeBackend, 1, 10)
	Nil(b.T(), err)

	// use to make sure we don't double use values
	intSet := sets.NewInt64()

	itr := res.Iterator()

	var lastValue uint64

	for !itr.Done() {
		index, value, _ := itr.Next()

		Falsef(b.T(), intSet.Has(int64(index)), "%d appears at least twice", index)
		intSet.Insert(int64(index))
		GreaterOrEqual(b.T(), value, lastValue)
		lastValue = value
	}
}

func (b *BackfillSuite) TestLogsInRange() {
	testBackend := geth.NewEmbeddedBackend(b.GetTestContext(), b.T())

	// start an omnirpc proxy and run 10 test tranactions so we can batch call blocks
	//  1-10
	var wg sync.WaitGroup
	wg.Add(2)

	const desiredBlockHeight = 10
	var commonAddress common.Address
	go func() {
		defer wg.Done()
		commonAddress = b.PopuluateWithLogs(b.GetTestContext(), testBackend, desiredBlockHeight)
	}()

	var host string
	go func() {
		defer wg.Done()
		host = b.startOmnirpcServer(b.GetTestContext(), testBackend)
	}()

	wg.Wait()

	scribeBackend, err := backfill.DialBackend(b.GetTestContext(), host)
	Nil(b.T(), err)

	res, err := backfill.GetLogsInRange(b.GetTestContext(), scribeBackend, 1, 10, 1, commonAddress)
	Nil(b.T(), err)

	// use to make sure we don't double use values
	intSet := sets.NewInt64()

	itr := res.Iterator()

	for !itr.Done() {
		index, _ := itr.Next()

		Falsef(b.T(), intSet.Has(int64(index)), "%d appears at least twice", index)
		intSet.Insert(int64(index))
	}
}