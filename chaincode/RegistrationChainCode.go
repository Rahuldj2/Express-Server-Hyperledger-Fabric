package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"time"
	"strconv"

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

//private data for testing
type PrivateData struct {
    ID   string `json:"id"`
    Data string `json:"data"`
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

	// Fetch the health records only if consent is granted by the patient
	if consent {
		// Org2 queries the health records within a valid window
		privateData, err := s.QueryHealthRecords(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to query health records: %v", err)
		}

		// Validate the health records against the policy criteria
		if policy.Criteria.IsNonSmoker != isNonSmoker {
			return fmt.Errorf("user's smoking status does not match policy criteria")
		}
		if policy.Criteria.HasDisease != hasDisease {
			return fmt.Errorf("user's disease status does not match policy criteria")
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

		return ctx.GetStub().PutState(fmt.Sprintf("%s-%s", userID, policyID), registrationJSON)
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



// UploadHealthRecords - Function to upload health record to the collection
func (s *SmartContract) UploadHealthRecords(ctx contractapi.TransactionContextInterface, healthRecord string) (string, error) {
    // Retrieve the identity of the submitting organization
    id, err := cid.New(ctx.GetStub()) // Use GetStub() to pass ChaincodeStubInterface
    if err != nil {
        return "", fmt.Errorf("failed to get client identity: %v", err)
    }

    org, err := id.GetMSPID()
    if err != nil {
        return "", fmt.Errorf("failed to get MSP ID: %v", err)
    }

    // Generate a simple ID using the current timestamp (can be incremented if needed)
    timestamp := time.Now().Unix()
    recordID := strconv.FormatInt(timestamp, 10) // Convert Unix timestamp to string

    // Define the collection based on the organization
    var collectionName string
    if org == "Org1MSP" {
        collectionName = "Org1MSPPrivateCollection"
    } else if org == "Org2MSP" {
        collectionName = "Org2MSPPrivateCollection"
    } else {
        return "", fmt.Errorf("organization not recognized: %s", org)
    }

    // Prepare the health record data
    recordData := []byte(healthRecord)

    // Store the health record in the private collection of the organization
    err = ctx.GetStub().PutPrivateData(collectionName, recordID, recordData)
    if err != nil {
        return "", fmt.Errorf("failed to store record in private collection: %v", err)
    }

    // Also store the record in the shared collection to be accessible by both Org1 and Org2
    err = ctx.GetStub().PutPrivateData("InsuranceregistrationCollection", recordID, recordData)
    if err != nil {
        return "", fmt.Errorf("failed to store record in shared collection: %v", err)
    }

    // Return the record ID for further queries
    return recordID, nil
}

// QueryHealthRecords - Function to query health records
func (s *SmartContract) QueryHealthRecords(ctx contractapi.TransactionContextInterface, recordID string) (string, error) {
    // Retrieve the identity of the client
    id, err := cid.New(ctx.GetStub()) // Use GetStub() to pass ChaincodeStubInterface
    if err != nil {
        return "", fmt.Errorf("failed to get client identity: %v", err)
    }

    org, err := id.GetMSPID()
    if err != nil {
        return "", fmt.Errorf("failed to get MSP ID: %v", err)
    }

    // Define the collection based on the organization
    var collectionName string
    if org == "Org1MSP" {
        collectionName = "Org1MSPPrivateCollection"
    } else if org == "Org2MSP" {
        collectionName = "Org2MSPPrivateCollection"
    } else {
        return "", fmt.Errorf("organization not recognized: %s", org)
    }

    // Try to retrieve the health record from the private collection
    recordData, err := ctx.GetStub().GetPrivateData(collectionName, recordID)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve record from private collection: %v", err)
    }

    if recordData != nil {
        return string(recordData), nil
    }

    // If the record is not found in the private collection, query the shared collection
    recordData, err = ctx.GetStub().GetPrivateData("InsuranceregistrationCollection", recordID)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve record from shared collection: %v", err)
    }

    if recordData == nil {
        return "", fmt.Errorf("record not found")
    }

    // Return the record data
    return string(recordData), nil
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















// StorePrivateData: Allows Org1 to store private data
// func (s *SmartContract) StorePrivateData(ctx contractapi.TransactionContextInterface, id string, data string) error {
//     // Check if the caller is Org1
//     orgID, err := ctx.GetClientIdentity().GetMSPID()
//     if err != nil {
//         return fmt.Errorf("failed to get client identity: %v", err)
//     }
//     if orgID != "Org1MSP" {
//         return fmt.Errorf("Org1 is the only allowed organization to write this data")
//     }

//     privateData := PrivateData{
//         ID:   id,
//         Data: data,
//     }

//     // Store private data in Org1's private collection
//     privateDataJSON, err := json.Marshal(privateData)
//     if err != nil {
//         return fmt.Errorf("failed to marshal private data: %v", err)
//     }

//     return ctx.GetStub().PutPrivateData("Org1MSPPrivateCollection", id, privateDataJSON)
// }

// // GetPrivateData: Allows Org1 to read its own private data
// func (s *SmartContract) GetPrivateData(ctx contractapi.TransactionContextInterface, id string) (*PrivateData, error) {
//     // Check if the caller is Org1
//     orgID, err := ctx.GetClientIdentity().GetMSPID()
//     if err != nil {
//         return nil, fmt.Errorf("failed to get client identity: %v", err)
//     }
//     if orgID != "Org1MSP" {
//         return nil, fmt.Errorf("Org1 is the only allowed organization to read this data")
//     }

//     // Retrieve the private data from Org1's private collection
//     privateDataJSON, err := ctx.GetStub().GetPrivateData("Org1MSPPrivateCollection", id)
//     if err != nil {
//         return nil, fmt.Errorf("failed to get private data: %v", err)
//     }
//     if privateDataJSON == nil {
//         return nil, fmt.Errorf("no private data found with ID %s", id)
//     }

//     var privateData PrivateData
//     err = json.Unmarshal(privateDataJSON, &privateData)
//     if err != nil {
//         return nil, fmt.Errorf("failed to unmarshal private data: %v", err)
//     }

//     return &privateData, nil
// }