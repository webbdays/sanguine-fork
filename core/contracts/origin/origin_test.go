package origin_test

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	. "github.com/stretchr/testify/assert"
	"github.com/synapsecns/sanguine/core/contracts/origin"
	"github.com/synapsecns/sanguine/core/types"
	"math/big"
	"time"
)

func (h OriginSuite) TestDispatchTopic() {
	// init the dispatch event
	txContext := h.testBackend.GetTxContext(h.GetTestContext(), nil)

	dispatchSink := make(chan *origin.OriginDispatch)
	sub, err := h.originContract.WatchDispatch(&bind.WatchOpts{Context: h.GetTestContext()}, dispatchSink, [][32]byte{}, []*big.Int{}, []uint64{})
	Nil(h.T(), err)

	enodedTips, err := types.EncodeTips(types.NewTips(big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)))
	Nil(h.T(), err)

	tx, err := h.originContract.Dispatch(txContext.TransactOpts, 1, [32]byte{}, 1, enodedTips, nil)
	Nil(h.T(), err)

	h.testBackend.WaitForConfirmation(h.GetTestContext(), tx)

	watchCtx, cancel := context.WithTimeout(h.GetTestContext(), time.Second*10)
	defer cancel()

	select {
	// check for errors and fail
	case <-watchCtx.Done():
		h.T().Error(h.T(), fmt.Errorf("test context completed %w", h.GetTestContext().Err()))
	case <-sub.Err():
		h.T().Error(h.T(), sub.Err())
	// get dispatch event
	case item := <-dispatchSink:
		parser, err := origin.NewParser(h.originContract.Address())
		Nil(h.T(), err)

		eventType, ok := parser.EventType(item.Raw)
		True(h.T(), ok)
		Equal(h.T(), eventType, origin.DispatchEvent)

		break
	}
}

func (h OriginSuite) TestUpdateTopic() {
	h.T().Skip("TODO: test this. Mocker should be able to mock this out")
}