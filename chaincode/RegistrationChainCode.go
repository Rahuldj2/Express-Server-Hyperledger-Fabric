package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing insurance policies
type SmartContract struct {
	contractapi.Contract
}

// InsurancePolicy represents an insurance policy with necessary details
type InsurancePolicy struct {
	PolicyID            string            `json:"policyId"`
	PolicyType          string            `json:"policyType"`
	CoverageAmount      float64           `json:"coverageAmount"`
	PremiumAmount       float64           `json:"premiumAmount"`
	PolicyStartDate     string            `json:"policyStartDate"`
	PolicyEndDate       string            `json:"policyEndDate"`
	TermsConditions     map[string]bool   `json:"termsConditions"` // Using a map for terms and conditions
}

// DefinePolicy allows insurance providers to define and register a new policy with a given PolicyID or a deterministic one.
func (s *SmartContract) DefinePolicy(ctx contractapi.TransactionContextInterface, policyId string, policyType string, coverageAmount float64, premiumAmount float64, startDate string, endDate string, termsConditions map[string]bool) (string, error) {
	// If the policyId is empty, create a deterministic PolicyID based on input fields
	if policyId == "" {
		deterministicID := fmt.Sprintf("%s-%s-%s-%f-%f", policyType, startDate, endDate, coverageAmount, premiumAmount)

		// Create a deterministic PolicyID based on the unique combination of the provided fields
		policyId = fmt.Sprintf("%x", sha256.Sum256([]byte(deterministicID)))
	}

	// Create an InsurancePolicy instance with TermsConditions as a map
	policy := InsurancePolicy{
		PolicyID:        policyId,
		PolicyType:      policyType,
		CoverageAmount:  coverageAmount,
		PremiumAmount:   premiumAmount,
		PolicyStartDate: startDate,
		PolicyEndDate:   endDate,
		TermsConditions: termsConditions,
	}

	// Serialize policy to JSON format
	policyBytes, err := json.Marshal(policy)
	if err != nil {
		return "", fmt.Errorf("failed to serialize policy: %v", err)
	}

	// Store policy on the ledger with policyId as the key
	err = ctx.GetStub().PutState(policyId, policyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to store policy on the ledger: %v", err)
	}

	return policyId, nil
}

// QueryPolicy retrieves an insurance policy by its policyId
func (s *SmartContract) QueryPolicy(ctx contractapi.TransactionContextInterface, policyId string) (*InsurancePolicy, error) {
	// Get the policy data from the ledger
	policyBytes, err := ctx.GetStub().GetState(policyId)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy from the ledger: %v", err)
	}
	if policyBytes == nil {
		return nil, fmt.Errorf("policy with ID %s does not exist", policyId)
	}

	// Deserialize policy data
	var policy InsurancePolicy
	err = json.Unmarshal(policyBytes, &policy)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize policy data: %v", err)
	}

	return &policy, nil
}

// main function to start the chaincode
func main() {
	// Create a new SmartContract object
	smartContract := new(SmartContract)

	// Create a new contract API that holds all the transactions
	chaincode, err := contractapi.NewChaincode(smartContract)
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err.Error())
		return
	}

	// Start the chaincode in the network
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
