package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DocContract contract for storing user in blockchain
type DocContract struct {
	contractapi.Contract
}

// CreateDoc creates doc
func (d *DocContract) CreateDoc(ctx contractapi.TransactionContextInterface, documentid string, documenthash []string, akcessid string) (string, error) {
	docAsBytes, err := ctx.GetStub().GetState(documentid)
	txid := ctx.GetStub().GetTxID()

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if docAsBytes != nil {
		return txid, fmt.Errorf("DocumentID with %s already exist", documentid)
	}

	doc := Document{
		ObjectType:        "document",
		DocumentID:        documentid,
		DocumentHash:      documenthash,
		Signature:         []Signature{},
		AkcessID:          akcessid,
		VerifiedBy:        map[string]time.Time{},
		VerificationGrade: []string{},
	}

	newDocAsBytes, _ := json.Marshal(doc)
	fmt.Printf("%s: Document with %s id created\n", txid, documentid)
	return txid, ctx.GetStub().PutState(documentid, newDocAsBytes)
}

// SignDoc signs doc with signature Hash
func (d *DocContract) SignDoc(ctx contractapi.TransactionContextInterface, documentid string, signhash string, signDate string, otpCode string, akcessid string) (string, error) {
	docAsBytes, err := ctx.GetStub().GetState(documentid)
	userAsBytes, err := ctx.GetStub().GetState(akcessid)
	txid := ctx.GetStub().GetTxID()

	signdate, err := time.Parse(time.RFC3339, signDate)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if docAsBytes == nil {
		return txid, fmt.Errorf("document with documentid %s doesn't exist", documentid)
	}

	if userAsBytes == nil {
		return txid, fmt.Errorf("AKcessId %s doesn't exist", akcessid)
	}
	var doc Document
	json.Unmarshal(docAsBytes, &doc)

	signature := Signature{
		SignatureHash: signhash,
		OTP:           otpCode,
		AkcessID:      akcessid,
		TimeStamp:     signdate,
	}

	doc.Signature = append(doc.Signature, signature)

	docAsBytes, _ = json.Marshal(doc)
	fmt.Printf("%s: Document %s signed by %s", txid, documentid, signhash)
	return txid, ctx.GetStub().PutState(documentid, docAsBytes)
}

// SendDoc shares document from sender to verifier
func (d *DocContract) SendDoc(ctx contractapi.TransactionContextInterface, sharingid string, sender string, verifier string, documentid string) (string, error) {
	senderAsBytes, err := ctx.GetStub().GetState(sender)
	verifierAsBytes, err := ctx.GetStub().GetState(verifier)
	docAsBytes, err := ctx.GetStub().GetState(documentid)
	txid := ctx.GetStub().GetTxID()

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if senderAsBytes == nil {
		return txid, fmt.Errorf("AKcessId %s doesn't exist", sender)
	}
	if verifierAsBytes == nil {
		return txid, fmt.Errorf("AKcessId %s doesn't exist", verifier)
	}
	if docAsBytes == nil {
		return txid, fmt.Errorf("Document with documentid %s doesn't exist", documentid)
	}

	sharedoc := DocumentShare{
		ObjectType: "docshare",
		SharingID:  sharingid,
		Sender:     sender,
		Verifier:   verifier,
		DocumentID: documentid,
	}

	shareSDocAdBytes, _ := json.Marshal(sharedoc)

	fmt.Printf("%s: Document %s shared from %s to %s\n", txid, documentid, sender, verifier)
	return txid, ctx.GetStub().PutState(sharingid, shareSDocAdBytes)
}

// VerifyDoc verify the doc
func (d *DocContract) VerifyDoc(ctx contractapi.TransactionContextInterface, akcessid string, documentid string, expiryDate string, verificationGrade string) (string, error) {
	docAsBytes, err := ctx.GetStub().GetState(documentid)
	txid := ctx.GetStub().GetTxID()

	expirydate, err := time.Parse(time.RFC3339, expiryDate)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return txid, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if docAsBytes == nil {
		return txid, fmt.Errorf("document with documentid %s doesn't exist", documentid)
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

	var doc Document
	json.Unmarshal(docAsBytes, &doc)

	doc.VerifiedBy[akcessid] = expirydate

	_, found := Find(doc.VerificationGrade, verificationGrade)
	if !found {
		doc.VerificationGrade = append(doc.VerificationGrade, verificationGrade)
	}

	docAsBytes, _ = json.Marshal(doc)

	fmt.Printf("%s: Document %s of verified by %s\n", txid, documentid, akcessid)
	return txid, ctx.GetStub().PutState(documentid, docAsBytes)
}

// GetTxForDoc get document details for perticular transaction
// func (d *DocContract) GetTxForDoc(ctx contractapi.TransactionContextInterface, documentid string, txid string) (*Document, error) {
// 	resultIterator, err := ctx.GetStub().GetHistoryForKey(documentid)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultIterator.Close()

// 	for resultIterator.HasNext() {
// 		queryResponse, err := resultIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}
// 		d := new(Document)
// 		_ = json.Unmarshal(queryResponse.Value, d)
// 		fmt.Println(queryResponse.TxId)
// 		fmt.Println(queryResponse.TxId == txid)
// 		if queryResponse.TxId == txid {
// 			queryResult := Document{
// 				ObjectType:        "document",
// 				DocumentID:        d.DocumentID,
// 				DocumentHash:      d.DocumentHash,
// 				SignatureHash:     d.SignatureHash,
// 				AkcessID:          d.AkcessID,
// 				VerifiedBy:        d.VerifiedBy,
// 				VerificationGrade: d.VerificationGrade,
// 			}
// 			return &queryResult, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("there is not tx with %s for document %s", txid, documentid)
// }

// GetVerifiersOfDoc get verifiers of perticular doc
func (d *DocContract) GetVerifiersOfDoc(ctx contractapi.TransactionContextInterface, documentid string) ([]string, error) {
	docAsBytes, err := ctx.GetStub().GetState(documentid)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if docAsBytes == nil {
		return nil, fmt.Errorf("document with documentid %s doesn't exist", documentid)
	}

	var doc Document
	json.Unmarshal(docAsBytes, &doc)
	verifiers := make([]string, len(doc.VerifiedBy))
	i := 0
	for v := range doc.VerifiedBy {
		verifiers[i] = v
		i++
	}
	return verifiers, nil
}

// GetSignature get signature by signature hash
func (d *DocContract) GetSignature(ctx contractapi.TransactionContextInterface, signHash string) ([]Document, error) {
	queryString := fmt.Sprintf(`{
		"selector": {
		   "docType": "document",
		   "signature": {
			  "$elemMatch": {
				 "signatureHash": "%s"
			  }
		   }
		}
	 }`, signHash)

	resultIterator, _ := ctx.GetStub().GetQueryResult(queryString)
	defer resultIterator.Close()

	result := []Document{}

	for resultIterator.HasNext() {
		queryResponse, _ := resultIterator.Next()

		doc := new(Document)
		_ = json.Unmarshal(queryResponse.Value, doc)
		result = append(result, *doc)
	}
	return result, nil
}
