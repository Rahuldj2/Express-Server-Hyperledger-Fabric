package main

// DEFINE POLICY
// REGISTER FOR POLICY
// UPLOAD HEALTH RECORDS
//QUERY HEALTH RECORDS


import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// Registration defines the structure for a policy registration
type Registration struct {
	UserID       string  `json:"userId"`
	PolicyID     string  `json:"policyId"`
	PremiumPaid  float64 `json:"premiumPaid"`
	IsNonSmoker  bool    `json:"isNonSmoker"`
	HasDisease   bool    `json:"hasDisease"`
}

// Modify the PrivateData struct to include the new boolean fields
type PrivateData struct {
	ID          string `json:"id"`
	IsNonSmoker bool   `json:"isNonSmoker"`
	HasDisease  bool   `json:"hasDisease"`
	Timestamp   int64  `json:"timestamp"`
}

// Policy defines the structure for a policy
type Policy struct {
	PolicyID      string            `json:"policyId"`
	PolicyType    string            `json:"policyType"`
	CoverAmount   float64           `json:"coverAmount"`
	Premium       float64           `json:"premium"`
	StartDate     string            `json:"startDate"`
	EndDate       string            `json:"endDate"`
	Criteria      Criteria          `json:"criteria"` // Changed to Criteria struct
	CoveredDiseases []string        `json:"coveredDiseases"`
}

// Criteria defines the structure for the criteria to be checked
type Criteria struct {
	IsNonSmoker  bool `json:"isNonSmoker"`
	HasDisease   bool `json:"hasDisease"`
}

// DefinePolicy: Allows Org2 to define a policy(insurance provider)
func (s *SmartContract) DefinePolicy(ctx contractapi.TransactionContextInterface, policyID, policyType string, coverAmount, premium float64, startDate, endDate, criteriaJSON string, diseasesJSON string) error {
	var criteria Criteria
	err := json.Unmarshal([]byte(criteriaJSON), &criteria)
	if err != nil {
		return fmt.Errorf("failed to parse criteria JSON: %v", err)
	}

	orgID, err1 := ctx.GetClientIdentity().GetMSPID()
	if err1 != nil {
		return fmt.Errorf("failed to get client identity: %v", err1)
	}
	if orgID != "Org2MSP" {
		return fmt.Errorf("only Org2 can define policies")
	}

	var coveredDiseases []string
	err = json.Unmarshal([]byte(diseasesJSON), &coveredDiseases)
	if err != nil {
		return fmt.Errorf("failed to parse diseases JSON: %v", err)
	}

	policy := Policy{
		PolicyID:      policyID,
		PolicyType:    policyType,
		CoverAmount:   coverAmount,
		Premium:       premium,
		StartDate:     startDate,
		EndDate:       endDate,
		Criteria:      criteria,
		CoveredDiseases: coveredDiseases,
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %v", err)
	}

	return ctx.GetStub().PutState(policyID, policyJSON)
}


// QueryPolicy: Retrieves the policy details by policyID
func (s *SmartContract) QueryPolicy(ctx contractapi.TransactionContextInterface, policyID string) (*Policy, error) {
	// Get policy JSON from the ledger
	policyJSON, err := ctx.GetStub().GetState(policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy from ledger: %v", err)
	}
	if policyJSON == nil {
		return nil, fmt.Errorf("policy with ID %s does not exist", policyID)
	}

	// Unmarshal the JSON into a Policy struct
	var policy Policy
	err = json.Unmarshal(policyJSON, &policy)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal policy JSON: %v", err)
	}

	return &policy, nil
}

// RegisterForPolicy: Allows users to register for a policy, while Org2 queries health records for validation
func (s *SmartContract) RegisterForPolicy(ctx contractapi.TransactionContextInterface, userID, policyID string, premiumPaid float64, isNonSmoker, hasDisease bool, consent bool) error {
	// Fetch the policy to validate if criteria match
	policyJSON, err := ctx.GetStub().GetState(policyID)
	if err != nil {
		return fmt.Errorf("failed to fetch policy with ID %s: %v", policyID, err)
	}
	if policyJSON == nil {
		return fmt.Errorf("policy with ID %s not found", policyID)
	}

	var policy Policy
	err = json.Unmarshal(policyJSON, &policy)
	if err != nil {
		return fmt.Errorf("failed to unmarshal policy: %v", err)
	}

	// Validate if the premium paid matches the required premium
	if premiumPaid != policy.Premium {
		return fmt.Errorf("premium paid %.2f does not match the required premium %.2f", premiumPaid, policy.Premium)
	}

	// Fetch the health records only if consent is granted by the patient
	if consent {
		// Org2 queries the health records within a valid window
		healthRecord, err := s.QueryHealthRecords(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to query health records: %v", err)
		}

		// Perform validation with the actual health record data
		if healthRecord.IsNonSmoker != isNonSmoker {
			return fmt.Errorf("user's smoking status does not match actual health records")
		}
		if healthRecord.HasDisease != hasDisease {
			return fmt.Errorf("user's disease status does not match actual health records")
		}

		// If validation passes, register the user for the policy
		registration := Registration{
			UserID:      userID,
			PolicyID:    policyID,
			PremiumPaid: premiumPaid,
			IsNonSmoker: isNonSmoker,
			HasDisease:  hasDisease,
		}

		// Store the registration
		registrationJSON, err := json.Marshal(registration)
		if err != nil {
			return fmt.Errorf("failed to marshal registration: %v", err)
		}

		err = ctx.GetStub().PutState(fmt.Sprintf("%s-%s", userID, policyID), registrationJSON)
		if err != nil {
			return fmt.Errorf("failed to store registration: %v", err)
		}

		// Store the userID -> policyID mapping for cross-chaincode access
		err = s.UpdateUserPolicyMapping(ctx, userID, policyID)
		if err != nil {
			return fmt.Errorf("failed to update user-policy mapping: %v", err)
		}

		return nil
	} else {
		return fmt.Errorf("patient consent is required to query health records")
	}
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

// UpdateUserPolicyMapping: Updates the userID -> policyID mapping
func (s *SmartContract) UpdateUserPolicyMapping(ctx contractapi.TransactionContextInterface, userID, policyID string) error {
	// Create a composite key for user-policy mapping
	mappingKey, err := ctx.GetStub().CreateCompositeKey("UserPolicyMapping", []string{userID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	// Store the mapping as userID -> policyID
	err = ctx.GetStub().PutState(mappingKey, []byte(policyID))
	if err != nil {
		return fmt.Errorf("failed to store user-policy mapping: %v", err)
	}

	return nil
}

// QueryPolicyByUserID: Queries the policy ID linked to a specific user ID
func (s *SmartContract) QueryPolicyByUserID(ctx contractapi.TransactionContextInterface, userID string) (string, error) {
	// Create a composite key for user-policy mapping
	mappingKey, err := ctx.GetStub().CreateCompositeKey("UserPolicyMapping", []string{userID})
	if err != nil {
		return "", fmt.Errorf("failed to create composite key: %v", err)
	}

	// Retrieve the policy ID using the mapping key
	policyIDBytes, err := ctx.GetStub().GetState(mappingKey)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve policy ID: %v", err)
	}
	if policyIDBytes == nil {
		return "", fmt.Errorf("no policy found for user ID %s", userID)
	}

	return string(policyIDBytes), nil
}





// UploadHealthRecords: Allows Org1 to upload health records with boolean values like isNonSmoker and hasDisease
func (s *SmartContract) UploadHealthRecords(ctx contractapi.TransactionContextInterface, id string, isNonSmoker, hasDisease bool) error {
	orgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client identity: %v", err)
	}
	if orgID != "Org1MSP" {
		return fmt.Errorf("only Org1 can upload health records")
	}

	// Create a struct for health record with boolean fields
	privateData := PrivateData{
		ID:          id,
		IsNonSmoker: isNonSmoker,
		HasDisease:  hasDisease,
		Timestamp:   time.Now().Unix(),
	}

	// Marshal the private data to JSON format
	privateDataJSON, err := json.Marshal(privateData)
	if err != nil {
		return fmt.Errorf("failed to marshal private data: %v", err)
	}

	// Store the private data in Org1MSP's private collection
	return ctx.GetStub().PutPrivateData("Org1MSPPrivateCollection", id, privateDataJSON)
}


// QueryHealthRecords: Allows Org2 to query health records within a time window, or Org1 at any time
func (s *SmartContract) QueryHealthRecords(ctx contractapi.TransactionContextInterface, id string) (*PrivateData, error) {
	orgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client identity: %v", err)
	}

	// Org2 can query health records only within a time window
	if orgID != "Org1MSP" && orgID != "Org2MSP" {
		return nil, fmt.Errorf("only Org1 and Org2 can query health records")
	}

	// Fetch the private data (health record)
	privateDataJSON, err := ctx.GetStub().GetPrivateData("Org1MSPPrivateCollection", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get private data: %v", err)
	}
	if privateDataJSON == nil {
		return nil, fmt.Errorf("no private data found with ID %s", id)
	}

	// Unmarshal the private data
	var privateData PrivateData
	err = json.Unmarshal(privateDataJSON, &privateData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private data: %v", err)
	}

	// If Org2, check if it's within the time window (e.g., 70 seconds)
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
