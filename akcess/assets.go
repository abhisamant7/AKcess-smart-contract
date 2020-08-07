package main

import "time"

// User describes basic details of user
type User struct {
	ObjectType    string                    `json:"docType"`
	AkcessID      string                    `json:"akcessid"`
	Verifications map[string][]Verification `json:"verifications"`
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
	VerifierAKcessID  string    `json:"verifierAKcessId"`
	VerifierName      string    `json:"verifierName"`
	ExpirtyDate       time.Time `json:"expiryDate"`
	VerificationGrade string    `json:"verificationGrade"`
}

// Document structure
type Document struct {
	ObjectType        string               `json:"docType"`
	DocumentID        string               `json:"documentid"`
	DocumentHash      []string             `json:"documenthash"`
	Signature         []Signature          `json:"signature"`
	AkcessID          string               `json:"akcessid"`
	VerifiedBy        map[string]time.Time `json:"verifiedby"`
	VerificationGrade []string             `json:"verificationGrade"`
}

// Signature structure
type Signature struct {
	SignatureHash string    `json:"signatureHash"`
	OTP           string    `json:"otp"`
	AkcessID      string    `json:"akcessId"`
	TimeStamp     time.Time `json:"timeStamp"`
}

// DocumentShare document object for share doc
type DocumentShare struct {
	ObjectType string `json:"docType"`
	SharingID  string `json:"sharingid"`
	Sender     string `json:"sender"`
	Verifier   string `json:"verifier"`
	DocumentID string `json:"documentid"`
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
		list = append(list, verification.VerifierAKcessID)
	}
	return list
}
