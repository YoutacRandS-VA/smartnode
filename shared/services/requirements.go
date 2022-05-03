package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/urfave/cli"
)

// Settings
const EthClientSyncTimeout = 16    // 16 seconds
const BeaconClientSyncTimeout = 16 // 16 seconds
var checkNodePasswordInterval, _ = time.ParseDuration("15s")
var checkNodeWalletInterval, _ = time.ParseDuration("15s")
var checkRocketStorageInterval, _ = time.ParseDuration("15s")
var checkNodeRegisteredInterval, _ = time.ParseDuration("15s")
var ethClientSyncPollInterval, _ = time.ParseDuration("5s")
var beaconClientSyncPollInterval, _ = time.ParseDuration("5s")
var ethClientRecentBlockThreshold, _ = time.ParseDuration("5m")
var ethClientStatusRefreshInterval, _ = time.ParseDuration("60s")

//
// Service requirements
//

func RequireNodePassword(c *cli.Context) error {
	nodePasswordSet, err := getNodePasswordSet(c)
	if err != nil {
		return err
	}
	if !nodePasswordSet {
		return errors.New("The node password has not been set. Please run 'rocketpool wallet init' and try again.")
	}
	return nil
}

func RequireNodeWallet(c *cli.Context) error {
	if err := RequireNodePassword(c); err != nil {
		return err
	}
	nodeWalletInitialized, err := getNodeWalletInitialized(c)
	if err != nil {
		return err
	}
	if !nodeWalletInitialized {
		return errors.New("The node wallet has not been initialized. Please run 'rocketpool wallet init' and try again.")
	}
	return nil
}

func RequireEthClientSynced(c *cli.Context) error {
	ethClientSynced, err := waitEthClientSynced(c, false, EthClientSyncTimeout)
	if err != nil {
		return err
	}
	if !ethClientSynced {
		return errors.New("The Eth 1.0 node is currently syncing. Please try again later.")
	}
	return nil
}

func RequireBeaconClientSynced(c *cli.Context) error {
	beaconClientSynced, err := waitBeaconClientSynced(c, false, BeaconClientSyncTimeout)
	if err != nil {
		return err
	}
	if !beaconClientSynced {
		return errors.New("The Eth 2.0 node is currently syncing. Please try again later.")
	}
	return nil
}

func RequireRocketStorage(c *cli.Context) error {
	if err := RequireEthClientSynced(c); err != nil {
		return err
	}
	rocketStorageLoaded, err := getRocketStorageLoaded(c)
	if err != nil {
		return err
	}
	if !rocketStorageLoaded {
		return errors.New("The Rocket Pool storage contract was not found; the configured address may be incorrect, or the Eth 1.0 node may not be synced. Please try again later.")
	}
	return nil
}

func RequireOneInchOracle(c *cli.Context) error {
	if err := RequireEthClientSynced(c); err != nil {
		return err
	}
	oneInchOracleLoaded, err := getOneInchOracleLoaded(c)
	if err != nil {
		return err
	}
	if !oneInchOracleLoaded {
		return errors.New("The 1inch oracle contract was not found; the configured address may be incorrect, or the mainnet Eth 1.0 node may not be synced. Please try again later.")
	}
	return nil
}

func RequireRplFaucet(c *cli.Context) error {
	if err := RequireEthClientSynced(c); err != nil {
		return err
	}
	rplFaucetLoaded, err := getRplFaucetLoaded(c)
	if err != nil {
		return err
	}
	if !rplFaucetLoaded {
		return errors.New("The RPL faucet contract was not found; the configured address may be incorrect, or the Eth 1.0 node may not be synced. Please try again later.")
	}
	return nil
}

func RequireNodeRegistered(c *cli.Context) error {
	if err := RequireNodeWallet(c); err != nil {
		return err
	}
	if err := RequireRocketStorage(c); err != nil {
		return err
	}
	nodeRegistered, err := getNodeRegistered(c)
	if err != nil {
		return err
	}
	if !nodeRegistered {
		return errors.New("The node is not registered with Rocket Pool. Please run 'rocketpool node register' and try again.")
	}
	return nil
}

func RequireNodeTrusted(c *cli.Context) error {
	if err := RequireNodeWallet(c); err != nil {
		return err
	}
	if err := RequireRocketStorage(c); err != nil {
		return err
	}
	nodeTrusted, err := getNodeTrusted(c)
	if err != nil {
		return err
	}
	if !nodeTrusted {
		return errors.New("The node is not a member of the oracle DAO. Nodes can only join the oracle DAO by invite.")
	}
	return nil
}

//
// Service synchronization
//

func WaitNodePassword(c *cli.Context, verbose bool) error {
	for {
		nodePasswordSet, err := getNodePasswordSet(c)
		if err != nil {
			return err
		}
		if nodePasswordSet {
			return nil
		}
		if verbose {
			log.Printf("The node password has not been set, retrying in %s...\n", checkNodePasswordInterval.String())
		}
		time.Sleep(checkNodePasswordInterval)
	}
}

func WaitNodeWallet(c *cli.Context, verbose bool) error {
	if err := WaitNodePassword(c, verbose); err != nil {
		return err
	}
	for {
		nodeWalletInitialized, err := getNodeWalletInitialized(c)
		if err != nil {
			return err
		}
		if nodeWalletInitialized {
			return nil
		}
		if verbose {
			log.Printf("The node wallet has not been initialized, retrying in %s...\n", checkNodeWalletInterval.String())
		}
		time.Sleep(checkNodeWalletInterval)
	}
}

func WaitEthClientSynced(c *cli.Context, verbose bool) error {
	_, err := waitEthClientSynced(c, verbose, 0)
	return err
}

func WaitBeaconClientSynced(c *cli.Context, verbose bool) error {
	_, err := waitBeaconClientSynced(c, verbose, 0)
	return err
}

func WaitRocketStorage(c *cli.Context, verbose bool) error {
	if err := WaitEthClientSynced(c, verbose); err != nil {
		return err
	}
	for {
		rocketStorageLoaded, err := getRocketStorageLoaded(c)
		if err != nil {
			return err
		}
		if rocketStorageLoaded {
			return nil
		}
		if verbose {
			log.Printf("The Rocket Pool storage contract was not found, retrying in %s...\n", checkRocketStorageInterval.String())
		}
		time.Sleep(checkRocketStorageInterval)
	}
}

func WaitNodeRegistered(c *cli.Context, verbose bool) error {
	if err := WaitNodeWallet(c, verbose); err != nil {
		return err
	}
	if err := WaitRocketStorage(c, verbose); err != nil {
		return err
	}
	for {
		nodeRegistered, err := getNodeRegistered(c)
		if err != nil {
			return err
		}
		if nodeRegistered {
			return nil
		}
		if verbose {
			log.Printf("The node is not registered with Rocket Pool, retrying in %s...\n", checkNodeRegisteredInterval.String())
		}
		time.Sleep(checkNodeRegisteredInterval)
	}
}

//
// Helpers
//

// Check if the node password is set
func getNodePasswordSet(c *cli.Context) (bool, error) {
	pm, err := GetPasswordManager(c)
	if err != nil {
		return false, err
	}
	return pm.IsPasswordSet(), nil
}

// Check if the node wallet is initialized
func getNodeWalletInitialized(c *cli.Context) (bool, error) {
	w, err := GetWallet(c)
	if err != nil {
		return false, err
	}
	return w.GetInitialized()
}

// Check if the RocketStorage contract is loaded
func getRocketStorageLoaded(c *cli.Context) (bool, error) {
	cfg, err := GetConfig(c)
	if err != nil {
		return false, err
	}
	ec, err := GetEthClient(c)
	if err != nil {
		return false, err
	}
	code, err := ec.CodeAt(context.Background(), common.HexToAddress(cfg.Smartnode.GetStorageAddress()), nil)
	if err != nil {
		return false, err
	}
	return (len(code) > 0), nil
}

// Check if the 1inch exchange oracle contract is loaded
func getOneInchOracleLoaded(c *cli.Context) (bool, error) {
	cfg, err := GetConfig(c)
	if err != nil {
		return false, err
	}
	ec, err := GetEthClient(c)
	if err != nil {
		return false, err
	}
	code, err := ec.CodeAt(context.Background(), common.HexToAddress(cfg.Smartnode.GetOneInchOracleAddress()), nil)
	if err != nil {
		return false, err
	}
	return (len(code) > 0), nil
}

// Check if the RPL faucet contract is loaded
func getRplFaucetLoaded(c *cli.Context) (bool, error) {
	cfg, err := GetConfig(c)
	if err != nil {
		return false, err
	}
	ec, err := GetEthClient(c)
	if err != nil {
		return false, err
	}
	code, err := ec.CodeAt(context.Background(), common.HexToAddress(cfg.Smartnode.GetRplFaucetAddress()), nil)
	if err != nil {
		return false, err
	}
	return (len(code) > 0), nil
}

// Check if the node is registered
func getNodeRegistered(c *cli.Context) (bool, error) {
	w, err := GetWallet(c)
	if err != nil {
		return false, err
	}
	rp, err := GetRocketPool(c)
	if err != nil {
		return false, err
	}
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return false, err
	}
	return node.GetNodeExists(rp, nodeAccount.Address, nil)
}

// Check if the node is a member of the oracle DAO
func getNodeTrusted(c *cli.Context) (bool, error) {
	w, err := GetWallet(c)
	if err != nil {
		return false, err
	}
	rp, err := GetRocketPool(c)
	if err != nil {
		return false, err
	}
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return false, err
	}
	return trustednode.GetMemberExists(rp, nodeAccount.Address, nil)
}

// Wait for the eth client to sync
// timeout of 0 indicates no timeout
var ethClientSyncLock sync.Mutex

func checkExecutionClientStatus(ecMgr *ExecutionClientManager) (bool, rocketpool.ExecutionClient, error) {

	// Check the EC status
	mgrStatus := ecMgr.CheckStatus()
	if ecMgr.primaryReady {
		return true, nil, nil
	}

	// If the primary isn't synced but there's a fallback and it is, return true
	if ecMgr.fallbackReady {
		if mgrStatus.PrimaryEcStatus.Error != "" {
			log.Printf("Primary execution client is unavailable (%s), using fallback execution client...\n", mgrStatus.PrimaryEcStatus.Error)
		} else {
			log.Printf("Primary execution client is still syncing (%.2f%%), using fallback execution client...\n", mgrStatus.PrimaryEcStatus.SyncProgress*100)
		}
		return true, nil, nil
	}

	// If neither is synced, go through the status to figure out what to do

	// Is the primary working and syncing? If so, wait for it
	if mgrStatus.PrimaryEcStatus.IsWorking && mgrStatus.PrimaryEcStatus.Error == "" {
		log.Printf("Fallback execution client is not configured or unavailable, waiting for primary execution client to finish syncing (%.2f%%)\n", mgrStatus.PrimaryEcStatus.SyncProgress)
		return false, ecMgr.primaryEc, nil
	}

	// Is the fallback working and syncing? If so, wait for it
	if mgrStatus.FallbackEnabled && mgrStatus.FallbackEcStatus.IsWorking && mgrStatus.FallbackEcStatus.Error == "" {
		log.Printf("Primary execution client is unavailable (%s), waiting for the fallback execution client to finish syncing (%.2f%%)\n", mgrStatus.PrimaryEcStatus.Error, mgrStatus.FallbackEcStatus.SyncProgress)
		return false, ecMgr.fallbackEc, nil
	}

	// If neither client is working, report the errors
	if mgrStatus.FallbackEnabled {
		return false, nil, fmt.Errorf("Primary execution client is unavailable (%s) and fallback execution client is unavailable (%s), no execution clients are ready.", mgrStatus.PrimaryEcStatus.Error, mgrStatus.FallbackEcStatus.Error)
	} else {
		return false, nil, fmt.Errorf("Primary execution client is unavailable (%s) and no fallback execution client is configured.", mgrStatus.PrimaryEcStatus.Error)
	}
}

func waitEthClientSynced(c *cli.Context, verbose bool, timeout int64) (bool, error) {

	// Prevent multiple waiting goroutines from requesting sync progress
	ethClientSyncLock.Lock()
	defer ethClientSyncLock.Unlock()

	// Get eth client
	var err error
	ecMgr, err := GetEthClient(c)
	if err != nil {
		return false, err
	}

	synced, clientToCheck, err := checkExecutionClientStatus(ecMgr)
	if err != nil {
		return false, err
	}
	if synced {
		return true, nil
	}

	// Get wait start time
	startTime := time.Now()

	// Get EC status refresh time
	ecRefreshTime := startTime

	// Wait for sync
	for {

		// Check timeout
		if (timeout > 0) && (time.Since(startTime).Seconds() > float64(timeout)) {
			return false, nil
		}

		// Check if the EC status needs to be refreshed
		if time.Since(ecRefreshTime) > ethClientStatusRefreshInterval {
			log.Println("Refreshing primary / fallback execution client status...")
			ecRefreshTime = time.Now()
			synced, clientToCheck, err = checkExecutionClientStatus(ecMgr)
			if err != nil {
				return false, err
			}
			if synced {
				return true, nil
			}
		}

		// Get sync progress
		progress, err := clientToCheck.SyncProgress(context.Background())
		if err != nil {
			return false, err
		}

		// Check sync progress
		if progress != nil {
			if verbose {
				p := float64(progress.CurrentBlock-progress.StartingBlock) / float64(progress.HighestBlock-progress.StartingBlock)
				if p > 1 {
					log.Println("Eth 1.0 node syncing...")
				} else {
					log.Printf("Eth 1.0 node syncing: %.2f%%\n", p*100)
				}
			}
		} else {
			// Eth 1 client is not in "syncing" state but may be behind head
			// Get the latest block it knows about and make sure it's recent compared to system clock time
			isUpToDate, _, err := IsSyncWithinThreshold(clientToCheck)
			if err != nil {
				return false, err
			}
			// Only return true if the last reportedly known block is within our defined threshold
			if isUpToDate {
				return true, nil
			}
		}

		// Pause before next poll
		time.Sleep(ethClientSyncPollInterval)

	}

}

// Wait for the beacon client to sync
// timeout of 0 indicates no timeout
var beaconClientSyncLock sync.Mutex

func waitBeaconClientSynced(c *cli.Context, verbose bool, timeout int64) (bool, error) {

	// Prevent multiple waiting goroutines from requesting sync progress
	beaconClientSyncLock.Lock()
	defer beaconClientSyncLock.Unlock()

	// Get beacon client
	bc, err := GetBeaconClient(c)
	if err != nil {
		return false, err
	}

	// Get wait start time
	startTime := time.Now().Unix()

	// Wait for sync
	for {

		// Check timeout
		if (timeout > 0) && (time.Now().Unix()-startTime > timeout) {
			return false, nil
		}

		// Get sync status
		syncStatus, err := bc.GetSyncStatus()
		if err != nil {
			return false, err
		}

		// Check sync status
		if syncStatus.Syncing {
			if verbose {
				log.Println("Eth 2.0 node syncing...")
			}
		} else {
			return true, nil
		}

		// Pause before next poll
		time.Sleep(beaconClientSyncPollInterval)

	}

}

// Confirm the EC's latest block is within the threshold of the current system clock
func IsSyncWithinThreshold(ec rocketpool.ExecutionClient) (bool, time.Time, error) {
	timestamp, err := GetEthClientLatestBlockTimestamp(ec)
	if err != nil {
		return false, time.Time{}, err
	}

	// Return true if the latest block is under the threshold
	blockTime := time.Unix(int64(timestamp), 0)
	if time.Since(blockTime) < ethClientRecentBlockThreshold {
		return true, blockTime, nil
	} else {
		return false, blockTime, nil
	}
}
