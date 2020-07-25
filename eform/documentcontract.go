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
func (d *EformContract) CreateEform(ctx contractapi.TransactionContextInterface, efomrid string, eformHash []string, akcessid string) (string, error) {
	eformAsBytes, err := ctx.GetStub().GetState(efomrid)
	txid := ctx.GetStub().GetTxID()

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if eformAsBytes != nil {
		return txid, fmt.Errorf("EformID with %s already exist", efomrid)
	}

	eform := Eform{
		ObjectType:        "efomr",
		EformID:           efomrid,
		EformHash:         eformHash,
		SignatureHash:     []string{},
		AkcessID:          akcessid,
		VerifiedBy:        map[string]time.Time{},
		VerificationGrade: []string{},
	}

	newEformAsBytes, _ := json.Marshal(eform)
	fmt.Printf("Eform with %s id created\n", efomrid)
	return txid, ctx.GetStub().PutState(efomrid, newEformAsBytes)
}

// SignEform signs the eform
func (d *EformContract) SignEform(ctx contractapi.TransactionContextInterface, efomrid string, signhash string, otpCode string, akcessid string) (string, error) {
	eformAsBytes, err := ctx.GetStub().GetState(efomrid)
	txid := ctx.GetStub().GetTxID()

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if eformAsBytes == nil {
		return txid, fmt.Errorf("efomr with efomrid %s doesn't exist", efomrid)
	}

	var eform Eform
	json.Unmarshal(eformAsBytes, &eform)
	eform.OTP = otpCode
	_, found := Find(eform.SignatureHash, signhash)
	if !found {
		eform.SignatureHash = append(eform.SignatureHash, signhash)
	}

	eformAsBytes, _ = json.Marshal(eform)
	return txid, ctx.GetStub().PutState(efomrid, eformAsBytes)
}

// SendEform shares efomr from sender to verifier
func (d *EformContract) SendEform(ctx contractapi.TransactionContextInterface, sharingid string, sender string, verifier string, efomrid string) (string, error) {
	eformAsBytes, err := ctx.GetStub().GetState(efomrid)
	txid := ctx.GetStub().GetTxID()

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if eformAsBytes == nil {
		return txid, fmt.Errorf("Eform with efomrid %s doesn't exist", efomrid)
	}

	shareeform := EformShare{
		ObjectType: "eformshare",
		SharingID:  sharingid,
		Sender:     sender,
		Verifier:   verifier,
		EformID:    efomrid,
	}

	shareEformAsBytes, _ := json.Marshal(shareeform)

	fmt.Printf("Eform %s shared from %s to %s\n", efomrid, sender, verifier)
	return txid, ctx.GetStub().PutState(sharingid, shareEformAsBytes)
}

// VerifyEform verify the eform
func (d *EformContract) VerifyEform(ctx contractapi.TransactionContextInterface, akcessid string, efomrid string, expiryDate string, verificationGrade string) (string, error) {
	eformAsBytes, err := ctx.GetStub().GetState(efomrid)
	txid := ctx.GetStub().GetTxID()

	expirydate, err := time.Parse(time.RFC3339, expiryDate)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if eformAsBytes == nil {
		return txid, fmt.Errorf("efomr with efomrid %s doesn't exist", efomrid)
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

	fmt.Printf("Eform %s of verified by %s\n", efomrid, akcessid)
	return txid, ctx.GetStub().PutState(efomrid, eformAsBytes)
}

// GetTxForEform get efomr details for perticular transaction
// func (d *EformContract) GetTxForEform(ctx contractapi.TransactionContextInterface, efomrid string, txid string) (*Eform, error) {
// 	resultIterator, err := ctx.GetStub().GetHistoryForKey(efomrid)
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
// 				ObjectType:        "efomr",
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
// 	return nil, fmt.Errorf("there is not tx with %s for efomr %s", txid, efomrid)
// }

// GetVerifiersOfEform get verifiers of perticular eform
func (d *EformContract) GetVerifiersOfEform(ctx contractapi.TransactionContextInterface, efomrid string) ([]string, error) {
	eformAsBytes, err := ctx.GetStub().GetState(efomrid)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if eformAsBytes == nil {
		return nil, fmt.Errorf("efomr with efomrid %s doesn't exist", efomrid)
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
