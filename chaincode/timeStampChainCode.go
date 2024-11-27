package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// Policy defines the structure for a policy
type Policy struct {
	PolicyID      string            `json:"policyId"`
	PolicyType    string            `json:"policyType"`
	CoverAmount   float64           `json:"coverAmount"`
	Premium       float64           `json:"premium"`
	StartDate     string            `json:"startDate"`
	EndDate       string            `json:"endDate"`
	Criteria      map[string]bool   `json:"criteria"`
}

// Registration defines the structure for a policy registration
type Registration struct {
	UserID       string  `json:"userId"`
	PolicyID     string  `json:"policyId"`
	PremiumPaid  float64 `json:"premiumPaid"`
	IsNonSmoker  bool    `json:"isNonSmoker"`
	HasDisease   bool    `json:"hasDisease"`
}

// PrivateData defines the structure for health records
type PrivateData struct {
	ID        string `json:"id"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

// DefinePolicy: Allows Org1 to define a policy
func (s *SmartContract) DefinePolicy(ctx contractapi.TransactionContextInterface, policyID, policyType string, coverAmount, premium float64, startDate, endDate string, criteriaJSON string) error {
	var criteria map[string]bool
	err := json.Unmarshal([]byte(criteriaJSON), &criteria)
	if err != nil {
		return fmt.Errorf("failed to parse criteria JSON: %v", err)
	}

	policy := Policy{
		PolicyID:    policyID,
		PolicyType:  policyType,
		CoverAmount: coverAmount,
		Premium:     premium,
		StartDate:   startDate,
		EndDate:     endDate,
		Criteria:    criteria,
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %v", err)
	}

	return ctx.GetStub().PutState(policyID, policyJSON)
}

// RegisterForPolicy: Allows users to register for a policy
func (s *SmartContract) RegisterForPolicy(ctx contractapi.TransactionContextInterface, userID, policyID string, premiumPaid float64, isNonSmoker, hasDisease bool) error {
	registration := Registration{
		UserID:      userID,
		PolicyID:    policyID,
		PremiumPaid: premiumPaid,
		IsNonSmoker: isNonSmoker,
		HasDisease:  hasDisease,
	}

	registrationJSON, err := json.Marshal(registration)
	if err != nil {
		return fmt.Errorf("failed to marshal registration: %v", err)
	}

	return ctx.GetStub().PutState(fmt.Sprintf("%s-%s", userID, policyID), registrationJSON)
}

// UploadHealthRecords: Allows Org1 to upload health records
func (s *SmartContract) UploadHealthRecords(ctx contractapi.TransactionContextInterface, id, data string) error {
	orgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client identity: %v", err)
	}
	if orgID != "Org1MSP" {
		return fmt.Errorf("only Org1 can upload health records")
	}

	privateData := PrivateData{
		ID:        id,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	privateDataJSON, err := json.Marshal(privateData)
	if err != nil {
		return fmt.Errorf("failed to marshal private data: %v", err)
	}

	return ctx.GetStub().PutPrivateData("Org1MSPPrivateCollection", id, privateDataJSON)
}

// QueryHealthRecords: Allows Org2 to query health records during a limited time window
func (s *SmartContract) QueryHealthRecords(ctx contractapi.TransactionContextInterface, id string) (*PrivateData, error) {
	orgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client identity: %v", err)
	}

	// Org2 can query within a time window, Org1 can always query
	if orgID != "Org1MSP" && orgID != "Org2MSP" {
		return nil, fmt.Errorf("only Org1 and Org2 can query health records")
	}

	// Retrieve the private data
	privateDataJSON, err := ctx.GetStub().GetPrivateData("Org1MSPPrivateCollection", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get private data: %v", err)
	}
	if privateDataJSON == nil {
		return nil, fmt.Errorf("no private data found with ID %s", id)
	}

	var privateData PrivateData
	err = json.Unmarshal(privateDataJSON, &privateData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private data: %v", err)
	}

	// Check if the record is within the time window (70 seconds)
	currentTime := time.Now().Unix()
	if orgID == "Org2MSP" && currentTime-privateData.Timestamp > 70 {
		return nil, fmt.Errorf("health records are no longer available for query by Org2")
	}

	return &privateData, nil
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
