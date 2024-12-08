package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Claim struct {
	UserID           string  `json:"userId"`
	PolicyID         string  `json:"policyId"`
	SettlementAmount float64 `json:"settlementAmount"`
	HospitalName     string  `json:"hospitalName"`
	Status           string  `json:"status"` // Example: "Processed", "Pending"
}

// PatientDetails defines the structure for storing patient details in the private data collection
type PatientDetails struct {
	UserID         string `json:"userId"`
	DiseaseDiagnosis string `json:"diseaseDiagnosis"`
	TreatmentPlan   string `json:"treatmentPlan"`
	HospitalName    string `json:"hospitalName"`
	AdmissionDate   string `json:"admissionDate"`
	DischargeDate   string `json:"dischargeDate"`
	ClaimStatus     string `json:"claimStatus"`
}
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

type Criteria struct {
	IsNonSmoker  bool `json:"isNonSmoker"`
	HasDisease   bool `json:"hasDisease"`
}

var patientDetailsList = make(map[string]PatientDetails) // Map to store policies by policyID

var claimDetailsList = make(map[string]Claim)

// UploadPatientDetails allows Org1 to upload patient details to the PDC
func (s *SmartContract) UploadPatientDetails(ctx contractapi.TransactionContextInterface, userID string, diseaseDiagnosis string, treatmentPlan string, hospitalName string, admissionDate string, dischargeDate string) error {
	patientDetails := PatientDetails{
		UserID:          userID,
		DiseaseDiagnosis: diseaseDiagnosis,
		TreatmentPlan:    treatmentPlan,
		HospitalName:     hospitalName,
		AdmissionDate:    admissionDate,
		DischargeDate:    dischargeDate,
		ClaimStatus:     "Pending",
	}

	patientDetailsList[userID]=patientDetails
	// Serialize the patient details to JSON
	patientDetailsJSON, err := json.Marshal(patientDetails)
	if err != nil {
		return fmt.Errorf("failed to serialize patient details: %v", err)
	}



	// Store the patient details in the private data collection
	return ctx.GetStub().PutPrivateData("Org1MSPPrivateCollection", userID, patientDetailsJSON)
}


// / QueryAllPolicies retrieves all policies from the in-memory list or ledger
func (s *SmartContract) QueryAllPatientData(ctx contractapi.TransactionContextInterface) ([]PatientDetails, error) {
	var patients []PatientDetails

	// First, check if patients are available in the in-memory store
	if len(patientDetailsList) > 0 {
		for _, patient := range patientDetailsList {
			patients = append(patients, patient)
		}
	} else {
		// If no patients in-memory, query the ledger
		startKey := ""
		endKey := ""

		iterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve patients from ledger: %v", err)
		}
		defer iterator.Close()

		for iterator.HasNext() {
			queryResponse, err := iterator.Next()
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve next patient entry during iteration: %v", err)
			}

			var patient PatientDetails
			err = json.Unmarshal(queryResponse.Value, &patient)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal patient JSON from value: %v", err)
			}

			patients= append(patients, patient)
		}

		if len(patients) == 0 {
			return nil, fmt.Errorf("no patients found in the ledger")
		}
	}

	return patients, nil
}


// ProcessClaim processes a claim for a user and stores the claim details
func (s *SmartContract) ProcessClaim(ctx contractapi.TransactionContextInterface, userID string) error {
	// Step 1: Query the policy ID associated with the user from the RegistrationContract
	args := [][]byte{[]byte("QueryPolicyByUserID"), []byte(userID)}
	response := ctx.GetStub().InvokeChaincode("registration", args, "mychannel") // Use the channel name where RegistrationContract is deployed
	
	if response.Status != 200 {
		return fmt.Errorf("failed to query policy for user %s from RegistrationContract: %v", userID, response.Message)
	}
	
	policyID := string(response.Payload)
	if policyID == "" {
		return fmt.Errorf("no policy found for user %s", userID)
	}

	// Step 2: Retrieve the policy details using the policyID from RegistrationContract
	args = [][]byte{[]byte("QueryPolicy"), []byte(policyID)}
	response = ctx.GetStub().InvokeChaincode("registration", args, "mychannel") // channel name
	
	if response.Status != 200 {
		return fmt.Errorf("failed to query policy details for policyID %s from RegistrationContract: %v", policyID, response.Message)
	}

	var policy Policy
	err := json.Unmarshal(response.Payload, &policy)
	if err != nil {
		return fmt.Errorf("failed to unmarshal policy details: %v", err)
	}

	// Step 3: Fetch patient details from Org1's PDC
	patientDetailsJSON, err := ctx.GetStub().GetPrivateData("Org1MSPPrivateCollection", userID)
	if err != nil {
		return fmt.Errorf("failed to fetch patient details: %v", err)
	}
	if patientDetailsJSON == nil {
		return fmt.Errorf("patient details not found for user %s", userID)
	}

	var patientDetails PatientDetails
	err = json.Unmarshal(patientDetailsJSON, &patientDetails)
	if err != nil {
		return fmt.Errorf("failed to unmarshal patient details: %v", err)
	}

	// Step 4: Check if the disease diagnosed is covered by the policy
	diseaseCovered := false
	for _, disease := range policy.CoveredDiseases {
		if disease == patientDetails.DiseaseDiagnosis {
			diseaseCovered = true
			break
		}
	}
	if !diseaseCovered {
		return fmt.Errorf("disease %s is not covered by policy %s", patientDetails.DiseaseDiagnosis, policy.PolicyID)
	}

	// Step 5: Calculate the settlement amount (e.g., 50% of the cover amount for simplicity)
	settlementAmount := policy.CoverAmount * 0.5

	// Step 6: Store the claim details
	claim := Claim{
		UserID:           userID,
		PolicyID:         policyID,
		SettlementAmount: settlementAmount,
		HospitalName:     patientDetails.HospitalName,
		Status:           "Processed",
	}

	
	claimDetailsList[userID] = claim
	// Serialize the claim to JSON
	claimJSON, err := json.Marshal(claim)
	if err != nil {
		return fmt.Errorf("failed to serialize claim: %v", err)
	}

	patientDetails.ClaimStatus = "Processed"

	// Update the in-memory map (if the patient exists in the map)
	if _, exists := patientDetailsList[userID]; exists {
		patientDetailsList[userID] = patientDetails
	}

	// Serialize the updated patient details to JSON
	updatedPatientDetailsJSON, err := json.Marshal(patientDetails)
	if err != nil {
		return fmt.Errorf("failed to serialize updated patient details: %v", err)
	}

	// Update the private data collection with the new patient details
	err = ctx.GetStub().PutPrivateData("Org1MSPPrivateCollection", userID, updatedPatientDetailsJSON)
	if err != nil {
		return fmt.Errorf("failed to update patient details in PDC: %v", err)
	}



	// Store the claim details in a collection (e.g., "claims")
	err = ctx.GetStub().PutPrivateData("ClaimsPrivateCollection", userID, claimJSON)
	if err != nil {
		return fmt.Errorf("failed to store claim details: %v", err)
	}

	return nil
}



// QueryAllClaims retrieves all claims from the in-memory list or ledger
func (s *SmartContract) QueryAllClaims(ctx contractapi.TransactionContextInterface) ([]Claim, error) {
	var claims []Claim

	// First, check if claims are available in the in-memory store
	if len(claimDetailsList) > 0 {
		for _, claim := range claimDetailsList {
			claims = append(claims, claim)
		}
	} else {
		// If no claims in-memory, query the ledger
		startKey := ""
		endKey := ""

		iterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve claims from ledger: %v", err)
		}
		defer iterator.Close()

		for iterator.HasNext() {
			queryResponse, err := iterator.Next()
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve next claim entry during iteration: %v", err)
			}

			var claim Claim
			err = json.Unmarshal(queryResponse.Value, &claim)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal claim JSON from value: %v", err)
			}

			claims = append(claims, claim)
		}

		if len(claims) == 0 {
			return nil, fmt.Errorf("no claims found in the ledger")
		}
	}

	return claims, nil
}

// QueryClaim retrieves claim details by userID
func (s *SmartContract) QueryClaim(ctx contractapi.TransactionContextInterface, userID string) (*Claim, error) {
	claimJSON, err := ctx.GetStub().GetPrivateData("ClaimsPrivateCollection", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch claim: %v", err)
	}
	if claimJSON == nil {
		return nil, fmt.Errorf("no claim found for user %s", userID)
	}

	var claim Claim
	err = json.Unmarshal(claimJSON, &claim)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal claim: %v", err)
	}

	return &claim, nil
}



func main() {
	claimsContract := new(SmartContract)

	chaincode, err := contractapi.NewChaincode(claimsContract)
	if err != nil {
		fmt.Printf("Error creating ClaimsContract: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting ClaimsContract: %s", err)
	}
}
