package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// User describes basic details of user
type User struct {
	ObjectType    string                    `json:"docType"`
	AkcessID      string                    `json:"akcessid"`
	Verifications map[string][]Verification `json:"verifications"`
}

// Verifier schema
type Verifier struct {
	ObjectType    string `json:"docType"`
	AkcessID      string `json:"akcessid"` // AKcessID of a verifier
	VerifierName  string `json:"verifierName"`
	VerifierGrade string `json:"verifierGrade"`
}

// Verification schema
type Verification struct {
	VerifierObj Verifier  `json:"veriier"`
	ExpirtyDate time.Time `json:"expiryDate"` // when verification will expire
}

// Document structure
type Document struct {
	ObjectType    string         `json:"docType"`
	DocumentID    string         `json:"documentid"`
	DocumentHash  []string       `json:"documenthash"`
	Signature     []Signature    `json:"signature"`
	AkcessID      string         `json:"akcessid"` // AKcessID of user who owns the document
	Verifications []Verification `json:"verifications"`
}

// Signature structure
type Signature struct {
	SignatureHash string    `json:"signatureHash"`
	OTP           string    `json:"otp"`
	AkcessID      string    `json:"akcessId"`  // AKcessID of user who signs
	TimeStamp     time.Time `json:"timeStamp"` // timestamp when signature is performed
}

// DocumentShare document object for share doc
type DocumentShare struct {
	ObjectType string   `json:"docType"`
	SharingID  string   `json:"sharingid"`
	Sender     string   `json:"sender"`
	Receivers  []string `json:"receivers"`
	DocumentID string   `json:"documentid"`
}

// DigitalAsset AKcess digital asset
type DigitalAsset struct {
	UniqueAssetID string            `json:"uniqueAssetID"`
	AssetType     string            `json:"assetType"`
	Owner         string            `json:"owner"`
	Metadata      map[string]string `json:"metadata"`
	LinkedDocs    []string          `json:"linkedDocs"`
	Verifications []Verification    `json:"verifications"`
	Description   string            `json:"description"`
	AssetDocHash  string            `json:"assetDocHash"`
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

// Remove deletes and element at peticular index from slice
func Remove(s []Verification, i int) []Verification {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
