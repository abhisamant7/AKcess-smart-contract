package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Eform structure
type Eform struct {
	ObjectType    string         `json:"docType"`
	EformID       string         `json:"eformId"`
	EformHash     []string       `json:"eformHash"`
	Signature     []Signature    `json:"signature"`
	AkcessID      string         `json:"akcessid"`
	Verifications []Verification `json:"verifications"`
}

// Signature structure
type Signature struct {
	SignatureHash string    `json:"signatureHash"`
	OTP           string    `json:"otp"`
	AkcessID      string    `json:"akcessId"`
	TimeStamp     time.Time `json:"timeStamp"`
}

// EformShare eform object for share eform
type EformShare struct {
	ObjectType string `json:"docType"`
	SharingID  string `jaon:"sharingid"`
	Sender     string `json:"sender"`
	Verifier   string `json:"verifier"`
	EformID    string `json:"eformId"`
}

// Verifier schema
type Verifier struct {
	ObjectType    string `json:"docType"`
	AkcessID      string `json:"akcessid"`
	VerifierName  string `json:"verifierName"`
	VerifierGrade string `json:"verifierGrade"`
}

// Verification schema
type Verification struct {
	VerifierObj Verifier  `json:"veriier"`
	ExpirtyDate time.Time `json:"expiryDate"`
}

// Find check if item already exists in slice
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// VerifiersList get list of verifiers
func VerifiersList(v []Verification) []string {
	var list []string
	for _, verification := range v {
		list = append(list, verification.VerifierObj.AkcessID)
	}
	return list
}

// IsVerifier checks if user who is invoking transaction is verifier or not
func IsVerifier(ctx contractapi.TransactionContextInterface) bool {
	isVerifier, attr, err := ctx.GetClientIdentity().GetAttributeValue("isVerifier")
	if err != nil {
		fmt.Println("Error while getting attribute from verifier identity")
	}
	if attr == false {
		fmt.Println("isVerifier attribute for this identity is not set")
	}
	isverifier, err := strconv.ParseBool(isVerifier)
	return isverifier
}
