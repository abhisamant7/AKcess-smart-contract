package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DocContract contract for storing user in blockchain
type DocContract struct {
	contractapi.Contract
}

// CreateDoc creates doc
func (d *DocContract) CreateDoc(ctx contractapi.TransactionContextInterface, documentid string, documenthash []string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, _ := getCommonName(ctx)
	docAsBytes, err := ctx.GetStub().GetState(documentid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching doc from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if docAsBytes != nil {
		response.Message = fmt.Sprintf("Document with id %s already exist", documentid)
		logger.Info(response.Message)
		return response
	}

	doc := Document{
		ObjectType:    "document",
		DocumentID:    documentid,
		DocumentHash:  documenthash,
		Signature:     []Signature{},
		AkcessID:      invoker,
		Verifications: []Verification{},
	}

	newDocAsBytes, _ := json.Marshal(doc)
	err = ctx.GetStub().PutState(documentid, newDocAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while creating doc: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Document with id %s created", documentid)
	logger.Info(response.Message)
	return response
}

// SignDoc signs doc with signature Hash
func (d *DocContract) SignDoc(ctx contractapi.TransactionContextInterface, documentid string, signhash string, signDate string, otpCode string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, _ := getCommonName(ctx)
	docAsBytes, err := ctx.GetStub().GetState(documentid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching doc from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if docAsBytes == nil {
		response.Message = fmt.Sprintf("Document with id %s doesn't exist", documentid)
		logger.Info(response.Message)
		return response
	}

	userAsBytes, err := ctx.GetStub().GetState(invoker)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching user from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if userAsBytes == nil {
		response.Message = fmt.Sprintf("User with id %s doesn't exist", invoker)
		logger.Info(response.Message)
		return response
	}

	signdate, err := time.Parse(time.RFC3339, signDate)
	if err != nil {
		response.Message = fmt.Sprint("Error while parsing date pass date in ISO format")
		logger.Info(response.Message)
		return response
	}

	var doc Document
	json.Unmarshal(docAsBytes, &doc)
	signature := Signature{
		SignatureHash: signhash,
		OTP:           otpCode,
		AkcessID:      invoker,
		TimeStamp:     signdate,
	}
	doc.Signature = append(doc.Signature, signature)
	docAsBytes, _ = json.Marshal(doc)

	err = ctx.GetStub().PutState(documentid, docAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while updating signature in doc: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Document %s signed by %s", documentid, invoker)
	logger.Info(response.Message)
	return response
}

// SendDoc shares document from sender to verifier
func (d *DocContract) SendDoc(ctx contractapi.TransactionContextInterface, sharingid string, receivers []string, documentid string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	sender, _ := getCommonName(ctx)
	senderAsBytes, err := ctx.GetStub().GetState(sender)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching user from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if senderAsBytes == nil {
		response.Message = fmt.Sprintf("User with id %s doesn't exist", documentid)
		logger.Info(response.Message)
		return response
	}
	docAsBytes, err := ctx.GetStub().GetState(documentid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching doc from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if docAsBytes == nil {
		response.Message = fmt.Sprintf("Document with id %s doesn't exist", documentid)
		logger.Info(response.Message)
		return response
	}

	sharedoc := DocumentShare{
		ObjectType: "docshare",
		SharingID:  sharingid,
		Sender:     sender,
		Receivers:  receivers,
		DocumentID: documentid,
	}
	shareSDocAdBytes, _ := json.Marshal(sharedoc)
	err = ctx.GetStub().PutState(sharingid, shareSDocAdBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while sending doc: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Document %s shared from %s to %s", documentid, sender, receivers)
	logger.Info(response.Message)
	return response
}

// VerifyDoc verify the doc
func (d *DocContract) VerifyDoc(ctx contractapi.TransactionContextInterface, documentid string, expiryDate string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, _ := getCommonName(ctx)
	docAsBytes, err := ctx.GetStub().GetState(documentid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching doc from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if docAsBytes != nil {
		response.Message = fmt.Sprintf("Document with id %s already exist", documentid)
		logger.Info(response.Message)
		return response
	}

	expirydate, err := time.Parse(time.RFC3339, expiryDate)
	if err != nil {
		response.Message = fmt.Sprintf("Error while parsing error pass date in ISO format", err.Error())
		logger.Error(response.Message)
		return response
	}

	verifierAsBytes, err := ctx.GetStub().GetState(invoker)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching verifier from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if verifierAsBytes == nil {
		response.Message = fmt.Sprintf("Verifier with id %s doesn't exist", invoker)
		logger.Info(response.Message)
		return response
	}

	var verifier Verifier
	json.Unmarshal(verifierAsBytes, &verifier)

	var doc Document
	json.Unmarshal(docAsBytes, &doc)

	verification := Verification{
		VerifierObj: verifier,
		ExpirtyDate: expirydate,
	}

	verifierList := VerifiersList(doc.Verifications)
	_, found := Find(verifierList, invoker)
	if found {
		for i, v := range doc.Verifications {
			if v.VerifierObj.AkcessID == invoker {
				doc.Verifications[i].ExpirtyDate = expirydate
				break
			}
		}
	} else {
		doc.Verifications = append(doc.Verifications, verification)
	}

	docAsBytes, _ = json.Marshal(doc)
	err = ctx.GetStub().PutState(documentid, docAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while updating verification in doc: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Document %s verified by %s", documentid, invoker)
	logger.Info(response.Message)
	return response
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
func (d *DocContract) GetVerifiersOfDoc(ctx contractapi.TransactionContextInterface, documentid string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	docAsBytes, err := ctx.GetStub().GetState(documentid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching doc from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if docAsBytes == nil {
		response.Message = fmt.Sprintf("Document with id %s doesn't exist", documentid)
		logger.Info(response.Message)
		return response
	}

	var doc Document
	json.Unmarshal(docAsBytes, &doc)

	response.Data = doc.Verifications
	response.Success = true
	response.Message = fmt.Sprintf("Successfully fetched verifications of doc %s", documentid)
	logger.Info(response.Message)
	return response
}

// GetSignature get signature by signature hash
func (d *DocContract) GetSignature(ctx contractapi.TransactionContextInterface, signHash string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

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

	resultIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching signature: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	defer resultIterator.Close()

	result := []Document{}
	for resultIterator.HasNext() {
		queryResponse, _ := resultIterator.Next()

		doc := new(Document)
		_ = json.Unmarshal(queryResponse.Value, doc)
		result = append(result, *doc)
	}

	response.Data = result
	response.Success = true
	response.Message = fmt.Sprintf("Successfully fetched all docs with signature %s", signHash)
	logger.Info(response.Message)
	return response
}
