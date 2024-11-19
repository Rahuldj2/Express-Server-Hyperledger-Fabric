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

// Registration represents a user's registration for a policy
type Registration struct {
	UserID           string  `json:"userId"`
	PolicyID         string  `json:"policyId"`
	PremiumPaid      float64 `json:"premiumPaid"`
	IsNonSmoker      bool    `json:"isNonSmoker"`
	HasDisease       bool    `json:"hasDisease"`
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


// RegisterForPolicy allows a user to register for a policy
func (s *SmartContract) RegisterForPolicy(ctx contractapi.TransactionContextInterface, userId string, policyId string, premiumPaid float64, isNonSmoker bool, hasDisease bool) (string, error) {
	// Query the policy details to check the premium amount requirement
	policy, err := s.QueryPolicy(ctx, policyId)
	if err != nil {
		return "", err
	}

	// Check if the premium paid matches the required premium amount
	if premiumPaid != policy.PremiumAmount {
		return "", fmt.Errorf("the premium paid does not match the policy's required premium amount")
	}

	// Create the registration entry
	registration := Registration{
		UserID:      userId,
		PolicyID:    policyId,
		PremiumPaid: premiumPaid,
		IsNonSmoker: isNonSmoker,
		HasDisease:  hasDisease,
	}

	// Serialize the registration data
	registrationBytes, err := json.Marshal(registration)
	if err != nil {
		return "", fmt.Errorf("failed to serialize registration: %v", err)
	}

	// Store the registration data in the ledger
	registrationKey := fmt.Sprintf("%s-%s", userId, policyId)
	err = ctx.GetStub().PutState(registrationKey, registrationBytes)
	if err != nil {
		return "", fmt.Errorf("failed to store registration on the ledger: %v", err)
	}

	return "Registration successful", nil
}

// QueryRegistration retrieves the registration details for a user and a policy
func (s *SmartContract) QueryRegistration(ctx contractapi.TransactionContextInterface, userId string, policyId string) (*Registration, error) {
	// Query the registration using the user ID and policy ID
	registrationKey := fmt.Sprintf("%s-%s", userId, policyId)
	registrationBytes, err := ctx.GetStub().GetState(registrationKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read registration from the ledger: %v", err)
	}
	if registrationBytes == nil {
		return nil, fmt.Errorf("registration for user %s and policy %s does not exist", userId, policyId)
	}

	var registration Registration
	err = json.Unmarshal(registrationBytes, &registration)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize registration data: %v", err)
	}

	return &registration, nil
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
