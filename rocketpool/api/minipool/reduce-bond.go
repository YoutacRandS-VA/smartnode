package minipool

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/smartnode/rocketpool/common/server"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/input"
)

// ===============
// === Factory ===
// ===============

type minipoolReduceBondContextFactory struct {
	handler *MinipoolHandler
}

func (f *minipoolReduceBondContextFactory) Create(vars map[string]string) (*minipoolReduceBondContext, error) {
	c := &minipoolReduceBondContext{
		handler: f.handler,
	}
	inputErrs := []error{
		server.ValidateArg("addresses", vars, input.ValidateAddresses, &c.minipoolAddresses),
	}
	return c, errors.Join(inputErrs...)
}

func (f *minipoolReduceBondContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessRoute[*minipoolReduceBondContext, api.BatchTxInfoData](
		router, "reduce-bond", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type minipoolReduceBondContext struct {
	handler           *MinipoolHandler
	minipoolAddresses []common.Address
}

func (c *minipoolReduceBondContext) PrepareData(data *api.BatchTxInfoData, opts *bind.TransactOpts) error {
	return prepareMinipoolBatchTxData(c.handler.serviceProvider, c.minipoolAddresses, data, c.CreateTx, "reduce-bond")
}

func (c *minipoolReduceBondContext) CreateTx(mp minipool.IMinipool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	mpv3, success := minipool.GetMinipoolAsV3(mp)
	if !success {
		mpCommon := mp.GetCommonDetails()
		return nil, fmt.Errorf("cannot create v3 binding for minipool %s, version %d", mpCommon.Address.Hex(), mpCommon.Version)
	}
	return mpv3.ReduceBondAmount(opts)
}
