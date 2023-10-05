package rocketpool

import (
	"fmt"
	"math/big"

	"github.com/goccy/go-json"

	"github.com/rocket-pool/smartnode/shared/types/api"
)

// Get protocol DAO proposals
func (c *Client) PDAOProposals() (api.PDAOProposalsResponse, error) {
	responseBytes, err := c.callAPI("pdao proposals")
	if err != nil {
		return api.PDAOProposalsResponse{}, fmt.Errorf("Could not get protocol DAO proposals: %w", err)
	}
	var response api.PDAOProposalsResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.PDAOProposalsResponse{}, fmt.Errorf("Could not decode protocol DAO proposals response: %w", err)
	}
	if response.Error != "" {
		return api.PDAOProposalsResponse{}, fmt.Errorf("Could not get protocol DAO proposals: %s", response.Error)
	}
	return response, nil
}

// Get protocol DAO proposal details
func (c *Client) PDAOProposalDetails(proposalID uint64) (api.PDAOProposalResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("pdao proposal-details %d", proposalID))
	if err != nil {
		return api.PDAOProposalResponse{}, fmt.Errorf("Could not get protocol DAO proposal: %w", err)
	}
	var response api.PDAOProposalResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.PDAOProposalResponse{}, fmt.Errorf("Could not decode protocol DAO proposal response: %w", err)
	}
	if response.Error != "" {
		return api.PDAOProposalResponse{}, fmt.Errorf("Could not get protocol DAO proposal: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can vote on a proposal
func (c *Client) PDAOCanVoteProposal(proposalID uint64, support bool) (api.CanVoteOnPDAOProposalResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("pdao can-vote-proposal %d %t", proposalID, support))
	if err != nil {
		return api.CanVoteOnPDAOProposalResponse{}, fmt.Errorf("Could not get protocol DAO can-vote-proposal: %w", err)
	}
	var response api.CanVoteOnPDAOProposalResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanVoteOnPDAOProposalResponse{}, fmt.Errorf("Could not decode protocol DAO can-vote-proposal response: %w", err)
	}
	if response.Error != "" {
		return api.CanVoteOnPDAOProposalResponse{}, fmt.Errorf("Could not get protocol DAO can-vote-proposal: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can vote on a proposal
func (c *Client) PDAOVoteProposal(proposalID uint64, support bool) (api.VoteOnPDAOProposalResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("pdao vote-proposal %d %t", proposalID, support))
	if err != nil {
		return api.VoteOnPDAOProposalResponse{}, fmt.Errorf("Could not get protocol DAO vote-proposal: %w", err)
	}
	var response api.VoteOnPDAOProposalResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.VoteOnPDAOProposalResponse{}, fmt.Errorf("Could not decode protocol DAO vote-proposal response: %w", err)
	}
	if response.Error != "" {
		return api.VoteOnPDAOProposalResponse{}, fmt.Errorf("Could not get protocol DAO vote-proposal: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can execute a proposal
func (c *Client) PDAOCanExecuteProposal(proposalID uint64) (api.CanExecutePDAOProposalResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("pdao can-execute-proposal %d", proposalID))
	if err != nil {
		return api.CanExecutePDAOProposalResponse{}, fmt.Errorf("Could not get protocol DAO can-execute-proposal: %w", err)
	}
	var response api.CanExecutePDAOProposalResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanExecutePDAOProposalResponse{}, fmt.Errorf("Could not decode protocol DAO can-execute-proposal response: %w", err)
	}
	if response.Error != "" {
		return api.CanExecutePDAOProposalResponse{}, fmt.Errorf("Could not get protocol DAO can-execute-proposal: %s", response.Error)
	}
	return response, nil
}

// Execute a proposal
func (c *Client) PDAOExecuteProposal(proposalID uint64) (api.ExecutePDAOProposalResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("pdao execute-proposal %d", proposalID))
	if err != nil {
		return api.ExecutePDAOProposalResponse{}, fmt.Errorf("Could not get protocol DAO execute-proposal: %w", err)
	}
	var response api.ExecutePDAOProposalResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.ExecutePDAOProposalResponse{}, fmt.Errorf("Could not decode protocol DAO execute-proposal response: %w", err)
	}
	if response.Error != "" {
		return api.ExecutePDAOProposalResponse{}, fmt.Errorf("Could not get protocol DAO execute-proposal: %s", response.Error)
	}
	return response, nil
}

// Get protocol DAO settings
func (c *Client) PDAOGetSettings() (api.GetPDAOSettingsResponse, error) {
	responseBytes, err := c.callAPI("pdao get-settings")
	if err != nil {
		return api.GetPDAOSettingsResponse{}, fmt.Errorf("Could not get protocol DAO get-settings: %w", err)
	}
	var response api.GetPDAOSettingsResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.GetPDAOSettingsResponse{}, fmt.Errorf("Could not decode protocol DAO get-settings response: %w", err)
	}
	if response.Error != "" {
		return api.GetPDAOSettingsResponse{}, fmt.Errorf("Could not get protocol DAO get-settings: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can vote on a proposal
func (c *Client) PDAOCanProposeSetting(setting string, value string) (api.CanProposePDAOSettingResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("pdao can-propose-setting %s %s", setting, value))
	if err != nil {
		return api.CanProposePDAOSettingResponse{}, fmt.Errorf("Could not get protocol DAO can-propose-setting: %w", err)
	}
	var response api.CanProposePDAOSettingResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanProposePDAOSettingResponse{}, fmt.Errorf("Could not decode protocol DAO can-propose-setting response: %w", err)
	}
	if response.Error != "" {
		return api.CanProposePDAOSettingResponse{}, fmt.Errorf("Could not get protocol DAO can-propose-setting: %s", response.Error)
	}
	return response, nil
}

// Propose updating a PDAO setting (use can-propose-setting to get the pollard)
func (c *Client) PDAOProposeSetting(setting string, value string, blockNumber uint32, pollard string) (api.ProposePDAOSettingResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("pdao propose-setting %s %s %d %s", setting, value, blockNumber, pollard))
	if err != nil {
		return api.ProposePDAOSettingResponse{}, fmt.Errorf("Could not get protocol DAO propose-setting: %w", err)
	}
	var response api.ProposePDAOSettingResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.ProposePDAOSettingResponse{}, fmt.Errorf("Could not decode protocol DAO propose-setting response: %w", err)
	}
	if response.Error != "" {
		return api.ProposePDAOSettingResponse{}, fmt.Errorf("Could not get protocol DAO propose-setting: %s", response.Error)
	}
	return response, nil
}

// Get the allocation percentages of RPL rewards for the Oracle DAO, the Protocol DAO, and the node operators
func (c *Client) PDAOGetRewardsPercentages() (api.PDAOGetRewardsPercentagesResponse, error) {
	responseBytes, err := c.callAPI("pdao get-rewards-percentages")
	if err != nil {
		return api.PDAOGetRewardsPercentagesResponse{}, fmt.Errorf("Could not get protocol DAO get-rewards-percentages: %w", err)
	}
	var response api.PDAOGetRewardsPercentagesResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.PDAOGetRewardsPercentagesResponse{}, fmt.Errorf("Could not decode protocol DAO get-rewards-percentages response: %w", err)
	}
	if response.Error != "" {
		return api.PDAOGetRewardsPercentagesResponse{}, fmt.Errorf("Could not get protocol DAO get-rewards-percentages: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can propose new RPL rewards allocation percentages for the Oracle DAO, the Protocol DAO, and the node operators
func (c *Client) PDAOCanProposeRewardsPercentages(node *big.Int, odao *big.Int, pdao *big.Int) (api.PDAOCanProposeRewardsPercentagesResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("pdao can-propose-rewards-percentages %s %s %s", node.String(), odao.String(), pdao.String()))
	if err != nil {
		return api.PDAOCanProposeRewardsPercentagesResponse{}, fmt.Errorf("Could not get protocol DAO can-propose-rewards-percentages: %w", err)
	}
	var response api.PDAOCanProposeRewardsPercentagesResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.PDAOCanProposeRewardsPercentagesResponse{}, fmt.Errorf("Could not decode protocol DAO can-propose-rewards-percentages response: %w", err)
	}
	if response.Error != "" {
		return api.PDAOCanProposeRewardsPercentagesResponse{}, fmt.Errorf("Could not get protocol DAO can-propose-rewards-percentages: %s", response.Error)
	}
	return response, nil
}

// Propose new RPL rewards allocation percentages for the Oracle DAO, the Protocol DAO, and the node operators
func (c *Client) PDAOProposeRewardsPercentages(node *big.Int, odao *big.Int, pdao *big.Int, blockNumber uint32, pollard string) (api.ProposePDAOSettingResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("pdao propose-rewards-percentages %s %s %s %d %s", node, odao, pdao, blockNumber, pollard))
	if err != nil {
		return api.ProposePDAOSettingResponse{}, fmt.Errorf("Could not get protocol DAO propose-rewards-percentages: %w", err)
	}
	var response api.ProposePDAOSettingResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.ProposePDAOSettingResponse{}, fmt.Errorf("Could not decode protocol DAO propose-rewards-percentages response: %w", err)
	}
	if response.Error != "" {
		return api.ProposePDAOSettingResponse{}, fmt.Errorf("Could not get protocol DAO propose-rewards-percentages: %s", response.Error)
	}
	return response, nil
}
