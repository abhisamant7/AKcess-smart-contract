package main

import "time"

// Eform structure
type Eform struct {
	ObjectType        string               `json:"docType"`
	EformID           string               `json:"eformId"`
	EformHash         []string             `json:"eformHash"`
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

// EformShare eform object for share eform
type EformShare struct {
	ObjectType string `json:"docType"`
	SharingID  string `jaon:"sharingid"`
	Sender     string `json:"sender"`
	Verifier   string `json:"verifier"`
	EformID    string `json:"eformId"`
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
