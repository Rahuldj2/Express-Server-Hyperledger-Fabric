package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Policy struct {
	PolicyID   string json:"policyID"
	HolderName string json:"holderName"
	Status     string json:"status"
}

type Claim struct {
	ClaimID   string  json:"claimID"
	PolicyID  string  json:"policyID"
	Amount    float64 json:"amount"
	Status    string  json:"status"
}

// Register a new policy
func (s *SmartContract) RegisterPolicy(ctx contractapi.TransactionContextInterface, policyID string, holderName string) error {
	policy := Policy{
		PolicyID:   policyID,
		HolderName: holderName,
		Status:     "Active",
	}

	policyAsBytes, _ := json.Marshal(policy)
	return ctx.GetStub().PutState(policyID, policyAsBytes)
}

// Process a claim
func (s *SmartContract) ProcessClaim(ctx contractapi.TransactionContextInterface, claimID string, policyID string, amount float64) error {
	policyAsBytes, err := ctx.GetStub().GetState(policyID)
	if err != nil {
		return fmt.Errorf("Failed to read policy: %s", err.Error())
	}
	if policyAsBytes == nil {
		return fmt.Errorf("Policy does not exist")
	}

	claim := Claim{
		ClaimID:  claimID,
		PolicyID: policyID,
		Amount:   amount,
		Status:   "Processed",
	}

	claimAsBytes, _ := json.Marshal(claim)
	return ctx.GetStub().PutState(claimID, claimAsBytes)
}

// Query a policy by ID
func (s *SmartContract) QueryPolicy(ctx contractapi.TransactionContextInterface, policyID string) (*Policy, error) {
	policyAsBytes, err := ctx.GetStub().GetState(policyID)
	if err != nil {
		return nil, fmt.Errorf("Failed to read policy: %s", err.Error())
	}
	if policyAsBytes == nil {
		return nil, fmt.Errorf("Policy does not exist")
	}

	var policy Policy
	err = json.Unmarshal(policyAsBytes, &policy)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal policy: %s", err.Error())
	}

	return &policy, nil
}

// Query all policies
func (s *SmartContract) QueryAllPolicies(ctx contractapi.TransactionContextInterface) ([]*Policy, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"policy\"}}")
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve policies: %s", err.Error())
	}
	defer resultsIterator.Close()

	var policies []*Policy
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var policy Policy
		err = json.Unmarshal(queryResponse.Value, &policy)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal policy: %s", err.Error())
		}
		policies = append(policies, &policy)
	}
	return policies, nil
}

// Query a claim by ID
func (s *SmartContract) QueryClaim(ctx contractapi.TransactionContextInterface, claimID string) (*Claim, error) {
	claimAsBytes, err := ctx.GetStub().GetState(claimID)
	if err != nil {
		return nil, fmt.Errorf("Failed to read claim: %s", err.Error())
	}
	if claimAsBytes == nil {
		return nil, fmt.Errorf("Claim does not exist")
	}

	var claim Claim
	err = json.Unmarshal(claimAsBytes, &claim)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal claim: %s", err.Error())
	}

	return &claim, nil
}

// Query claims by policy ID
func (s *SmartContract) QueryClaimsByPolicy(ctx contractapi.TransactionContextInterface, policyID string) ([]*Claim, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"policyID\":\"%s\"}}", policyID)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve claims: %s", err.Error())
	}
	defer resultsIterator.Close()

	var claims []*Claim
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var claim Claim
		err = json.Unmarshal(queryResponse.Value, &claim)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal claim: %s", err.Error())
		}
		claims = append(claims, &claim)
	}
	return claims, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating insurance chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting insurance chaincode: %s", err.Error())
	}
}

