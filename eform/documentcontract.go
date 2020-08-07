package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// EformContract contract for storing user in blockchain
type EformContract struct {
	contractapi.Contract
}

// CreateEform creates eform
func (d *EformContract) CreateEform(ctx contractapi.TransactionContextInterface, eformid string, eformHash []string, akcessid string) (string, error) {
	eformAsBytes, err := ctx.GetStub().GetState(eformid)
	txid := ctx.GetStub().GetTxID()

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if eformAsBytes != nil {
		return txid, fmt.Errorf("EformID with %s already exist", eformid)
	}

	eform := Eform{
		ObjectType:        "eform",
		EformID:           eformid,
		EformHash:         eformHash,
		Signature:         []Signature{},
		AkcessID:          akcessid,
		VerifiedBy:        map[string]time.Time{},
		VerificationGrade: []string{},
	}

	newEformAsBytes, _ := json.Marshal(eform)
	fmt.Printf("%s: Eform with %s id created\n", txid, eformid)
	return txid, ctx.GetStub().PutState(eformid, newEformAsBytes)
}

// SignEform signs the eform
func (d *EformContract) SignEform(ctx contractapi.TransactionContextInterface, eformid string, signhash string, signDate string, otpCode string, akcessid string) (string, error) {
	eformAsBytes, err := ctx.GetStub().GetState(eformid)
	txid := ctx.GetStub().GetTxID()

	signdate, err := time.Parse(time.RFC3339, signDate)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if eformAsBytes == nil {
		return txid, fmt.Errorf("eform with eformid %s doesn't exist", eformid)
	}

	var eform Eform
	json.Unmarshal(eformAsBytes, &eform)

	signature := Signature{
		SignatureHash: signhash,
		OTP:           otpCode,
		AkcessID:      akcessid,
		TimeStamp:     signdate,
	}

	eform.Signature = append(eform.Signature, signature)

	eformAsBytes, _ = json.Marshal(eform)
	fmt.Printf("%s: Eform %s signed by %s\n", txid, eformid, akcessid)
	return txid, ctx.GetStub().PutState(eformid, eformAsBytes)
}

// SendEform shares eform from sender to verifier
func (d *EformContract) SendEform(ctx contractapi.TransactionContextInterface, sharingid string, sender string, verifier string, eformid string) (string, error) {
	eformAsBytes, err := ctx.GetStub().GetState(eformid)
	txid := ctx.GetStub().GetTxID()

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if eformAsBytes == nil {
		return txid, fmt.Errorf("Eform with eformid %s doesn't exist", eformid)
	}

	shareeform := EformShare{
		ObjectType: "eformshare",
		SharingID:  sharingid,
		Sender:     sender,
		Verifier:   verifier,
		EformID:    eformid,
	}

	shareEformAsBytes, _ := json.Marshal(shareeform)

	fmt.Printf("%s: Eform %s shared from %s to %s\n", txid, eformid, sender, verifier)
	return txid, ctx.GetStub().PutState(sharingid, shareEformAsBytes)
}

// VerifyEform verify the eform
func (d *EformContract) VerifyEform(ctx contractapi.TransactionContextInterface, akcessid string, eformid string, expiryDate string, verificationGrade string) (string, error) {
	eformAsBytes, err := ctx.GetStub().GetState(eformid)
	txid := ctx.GetStub().GetTxID()

	expirydate, err := time.Parse(time.RFC3339, expiryDate)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if eformAsBytes == nil {
		return txid, fmt.Errorf("eform with eformid %s doesn't exist", eformid)
	}

	attr, ok, err := cid.GetAttributeValue(ctx.GetStub(), "isVerifier")
	if err != nil {
		fmt.Println("An error getting attribue")
	}
	if !ok {
		fmt.Println("identity does not have this perticular attribute")
	}
	if attr != "true" {
		return txid, fmt.Errorf("User who is invoking transaction is not a verifier")
	}

	var eform Eform
	json.Unmarshal(eformAsBytes, &eform)

	eform.VerifiedBy[akcessid] = expirydate

	_, found := Find(eform.VerificationGrade, verificationGrade)
	if !found {
		eform.VerificationGrade = append(eform.VerificationGrade, verificationGrade)
	}

	eformAsBytes, _ = json.Marshal(eform)

	fmt.Printf("%s: Eform %s of verified by %s\n", txid, eformid, akcessid)
	return txid, ctx.GetStub().PutState(eformid, eformAsBytes)
}

// GetTxForEform get eform details for perticular transaction
// func (d *EformContract) GetTxForEform(ctx contractapi.TransactionContextInterface, eformid string, txid string) (*Eform, error) {
// 	resultIterator, err := ctx.GetStub().GetHistoryForKey(eformid)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultIterator.Close()

// 	for resultIterator.HasNext() {
// 		queryResponse, err := resultIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}
// 		d := new(Eform)
// 		_ = json.Unmarshal(queryResponse.Value, d)
// 		fmt.Println(queryResponse.TxId)
// 		fmt.Println(queryResponse.TxId == txid)
// 		if queryResponse.TxId == txid {
// 			queryResult := Eform{
// 				ObjectType:        "eform",
// 				EformID:        d.EformID,
// 				EformHash:      d.EformHash,
// 				SignatureHash:     d.SignatureHash,
// 				AkcessID:          d.AkcessID,
// 				VerifiedBy:        d.VerifiedBy,
// 				VerificationGrade: d.VerificationGrade,
// 			}
// 			return &queryResult, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("there is not tx with %s for eform %s", txid, eformid)
// }

// GetVerifiersOfEform get verifiers of perticular eform
func (d *EformContract) GetVerifiersOfEform(ctx contractapi.TransactionContextInterface, eformid string) ([]string, error) {
	eformAsBytes, err := ctx.GetStub().GetState(eformid)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if eformAsBytes == nil {
		return nil, fmt.Errorf("eform with eformid %s doesn't exist", eformid)
	}

	var eform Eform
	json.Unmarshal(eformAsBytes, &eform)
	verifiers := make([]string, len(eform.VerifiedBy))
	i := 0
	for v := range eform.VerifiedBy {
		verifiers[i] = v
		i++
	}
	return verifiers, nil
}

// GetSignature get signature by signature hash
func (d *EformContract) GetSignature(ctx contractapi.TransactionContextInterface, signHash string) ([]Eform, error) {
	queryString := fmt.Sprintf(`{
		"selector": {
		   "docType": "eform",
		   "signature": {
			  "$elemMatch": {
				 "signatureHash": "%s"
			  }
		   }
		}
	 }`, signHash)

	fmt.Println(queryString)
	resultIterator, _ := ctx.GetStub().GetQueryResult(queryString)
	defer resultIterator.Close()

	result := []Eform{}

	for resultIterator.HasNext() {
		queryResponse, _ := resultIterator.Next()

		eform := new(Eform)
		_ = json.Unmarshal(queryResponse.Value, eform)
		result = append(result, *eform)
		fmt.Println(result)
	}
	return result, nil
}
