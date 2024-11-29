package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ClaimsContract struct {
	contractapi.Contract
}

// Policy represents the policy details
type Policy struct {
    PolicyID       string   `json:"policyId"`
    PolicyType     string   `json:"policyType"`
    CoverageAmount float64  `json:"coverageAmount"`
    Premium        float64  `json:"premium"`
    StartDate      string   `json:"startDate"`
    EndDate        string   `json:"endDate"`
    Conditions     string   `json:"conditions"`
    CoveredDiseases []string `json:"coveredDiseases"`
}

// Claim represents the claim request structure
type Claim struct {
	PatientID      string  `json:"patientId"`
	Disease        string  `json:"disease"`
	TreatmentCost  float64 `json:"treatmentCost"`
}

// ProcessClaim: Verify claim details with the policy
// ProcessClaim: Verify claim details with the policy
func (c *ClaimsContract) ProcessClaim(ctx contractapi.TransactionContextInterface, patientID, disease string, treatmentCost float64) (string, error) {
	// Interact with the Registration chaincode to get policy details
	chaincodeName := "registration" // Name of the registration chaincode

	response := ctx.GetStub().InvokeChaincode(chaincodeName, [][]byte{[]byte("QueryPolicy"), []byte("Policy123")}, "")
	if response.Status != 200 {
		return "", fmt.Errorf("failed to query policy details: %s", response.Message)
	}

	// Unmarshal the policy details
	var policy Policy
	err := json.Unmarshal(response.Payload, &policy)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal policy details: %v", err)
	}

	// Verify if the disease is covered in the policy
	for _, coveredDisease := range policy.CoveredDiseases {
		if coveredDisease == disease {
			// Return "Approve" if the disease is covered
			return "Approve", nil
		}
	}

	// Return "Reject" if the disease is not covered
	return "Reject", nil
}


// main function to start the chaincode
func main() {
	claimsContract := new(ClaimsContract)

	cc, err := contractapi.NewChaincode(claimsContract)
	if err != nil {
		fmt.Printf("Error creating ClaimsProcessing chaincode: %s", err.Error())
		return
	}

	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting ClaimsProcessing chaincode: %s", err.Error())
	}
}
