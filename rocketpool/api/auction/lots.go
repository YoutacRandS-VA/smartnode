package auction

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/auction"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/smartnode/rocketpool/common/server"
	"github.com/rocket-pool/smartnode/shared/types/api"
)

// ===============
// === Factory ===
// ===============

type auctionLotContextFactory struct {
	handler *AuctionHandler
}

func (f *auctionLotContextFactory) Create(vars map[string]string) (*auctionLotContext, error) {
	c := &auctionLotContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *auctionLotContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterSingleStageRoute[*auctionLotContext, api.AuctionLotsData](
		router, "lots", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type auctionLotContext struct {
	handler     *AuctionHandler
	rp          *rocketpool.RocketPool
	nodeAddress common.Address

	auctionMgr *auction.AuctionManager
}

func (c *auctionLotContext) Initialize() error {
	sp := c.handler.serviceProvider
	c.rp = sp.GetRocketPool()
	c.nodeAddress, _ = sp.GetWallet().GetAddress()

	// Requirements
	err := sp.RequireNodeRegistered()
	if err != nil {
		return err
	}

	// Bindings
	c.auctionMgr, err = auction.NewAuctionManager(c.rp)
	if err != nil {
		return fmt.Errorf("error creating auction manager binding: %w", err)
	}
	return nil
}

func (c *auctionLotContext) GetState(mc *batch.MultiCaller) {
	c.auctionMgr.LotCount.AddToQuery(mc)
}

func (c *auctionLotContext) PrepareData(data *api.AuctionLotsData, opts *bind.TransactOpts) error {
	// Get lot details
	lotCount := c.auctionMgr.LotCount.Formatted()
	lots := make([]*auction.AuctionLot, lotCount)
	details := make([]api.AuctionLotDetails, lotCount)

	// Load details
	err := c.rp.BatchQuery(int(lotCount), int(lotCountDetailsBatchSize), func(mc *batch.MultiCaller, i int) error {
		lot, err := auction.NewAuctionLot(c.rp, uint64(i))
		if err != nil {
			return fmt.Errorf("error creating lot %d binding: %w", i, err)
		}
		lots[i] = lot
		core.QueryAllFields(lot, mc)
		lot.GetLotAddressBidAmount(mc, &details[i].NodeBidAmount, c.nodeAddress)
		return nil
	}, nil)
	if err != nil {
		return fmt.Errorf("error getting lot details: %w", err)
	}

	// Process details
	for i := 0; i < int(lotCount); i++ {
		fullDetails := &details[i]
		lot := lots[i]
		fullDetails.Index = lot.Index
		fullDetails.Exists = lot.Exists.Get()
		fullDetails.StartBlock = lot.StartBlock.Formatted()
		fullDetails.EndBlock = lot.EndBlock.Formatted()
		fullDetails.StartPrice = lot.StartPrice.Formatted()
		fullDetails.ReservePrice = lot.ReservePrice.Formatted()
		fullDetails.PriceAtCurrentBlock = lot.PriceAtCurrentBlock.Formatted()
		fullDetails.PriceByTotalBids = lot.PriceByTotalBids.Formatted()
		fullDetails.CurrentPrice = lot.CurrentPrice.Formatted()
		fullDetails.TotalRplAmount = lot.TotalRplAmount.Get()
		fullDetails.ClaimedRplAmount = lot.ClaimedRplAmount.Get()
		fullDetails.RemainingRplAmount = lot.RemainingRplAmount.Get()
		fullDetails.TotalBidAmount = lot.TotalBidAmount.Get()
		fullDetails.IsCleared = lot.IsCleared.Get()
		fullDetails.RplRecovered = lot.RplRecovered.Get()

		// Check lot conditions
		addressHasBid := (fullDetails.NodeBidAmount.Cmp(big.NewInt(0)) > 0)
		hasRemainingRpl := (fullDetails.RemainingRplAmount.Cmp(big.NewInt(0)) > 0)

		fullDetails.ClaimAvailable = (addressHasBid && fullDetails.IsCleared)
		fullDetails.BiddingAvailable = (!fullDetails.IsCleared && hasRemainingRpl)
		fullDetails.RplRecoveryAvailable = (fullDetails.IsCleared && hasRemainingRpl && !fullDetails.RplRecovered)
	}
	data.Lots = details
	return nil
}
