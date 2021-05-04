package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/util"
)

// EformContract contract for storing user in blockchain
type EformContract struct {
	contractapi.Contract
}

// CreateEform creates eform
func (d *EformContract) CreateEform(ctx contractapi.TransactionContextInterface, eformid string, eformHash []string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, _ := getCommonName(ctx)
	eformAsBytes, err := ctx.GetStub().GetState(eformid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching doc from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if eformAsBytes != nil {
		response.Message = fmt.Sprintf("Eform with id %s already exist", eformid)
		logger.Info(response.Message)
		return response
	}

	eform := Eform{
		ObjectType:    "eform",
		EformID:       eformid,
		EformHash:     eformHash,
		Signature:     []Signature{},
		AkcessID:      invoker,
		Verifications: []Verification{},
	}

	newEformAsBytes, _ := json.Marshal(eform)
	err = ctx.GetStub().PutState(eformid, newEformAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while creating eform: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Eform with id %s created", eformid)
	logger.Info(response.Message)
	return response
}

// SignEform signs the eform
func (d *EformContract) SignEform(ctx contractapi.TransactionContextInterface, eformid string, signhash string, signDate string, otpCode string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, _ := getCommonName(ctx)
	eformAsBytes, err := ctx.GetStub().GetState(eformid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching eform from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if eformAsBytes == nil {
		response.Message = fmt.Sprintf("Eform with id %s doesn't exist", eformid)
		logger.Info(response.Message)
		return response
	}

	signdate, err := time.Parse(time.RFC3339, signDate)
	if err != nil {
		response.Message = fmt.Sprintf("Error while parsing date pass date in ISO format: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	var eform Eform
	json.Unmarshal(eformAsBytes, &eform)
	signature := Signature{
		SignatureHash: signhash,
		OTP:           otpCode,
		AkcessID:      invoker,
		TimeStamp:     signdate,
	}
	eform.Signature = append(eform.Signature, signature)
	eformAsBytes, _ = json.Marshal(eform)
	err = ctx.GetStub().PutState(eformid, eformAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while signing eform: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Eform %s signed by %s", eformid, invoker)
	logger.Info(response.Message)
	return response
}

// SendEform shares eform from sender to verifier
func (d *EformContract) SendEform(ctx contractapi.TransactionContextInterface, sharingid string, receivers []string, eformid string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	sender, _ := getCommonName(ctx)
	eformAsBytes, err := ctx.GetStub().GetState(eformid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching eform from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if eformAsBytes == nil {
		response.Message = fmt.Sprintf("Eform with id %s doesn't exist", eformid)
		logger.Info(response.Message)
		return response
	}
	userAsBytes, err := ctx.GetStub().GetState(sender)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching user from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if userAsBytes == nil {
		response.Message = fmt.Sprintf("User with id %s doesn't exist", sender)
		logger.Info(response.Message)
		return response
	}

	shareeform := EformShare{
		ObjectType: "eformshare",
		SharingID:  sharingid,
		Sender:     sender,
		Receivers:  receivers,
		EformID:    eformid,
	}
	shareEformAsBytes, _ := json.Marshal(shareeform)
	err = ctx.GetStub().PutState(sharingid, shareEformAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while sharing eform: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Eform %s shared from %s to %s", eformid, sender, receivers)
	logger.Info(response.Message)
	return response
}

// VerifyEform verify the eform
func (d *EformContract) VerifyEform(ctx contractapi.TransactionContextInterface, eformid string, expiryDate string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, _ := getCommonName(ctx)
	eformAsBytes, err := ctx.GetStub().GetState(eformid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching eform from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if eformAsBytes == nil {
		response.Message = fmt.Sprintf("Eform with id %s doesn't exist", eformid)
		logger.Info(response.Message)
		return response
	}

	expirydate, err := time.Parse(time.RFC3339, expiryDate)
	if err != nil {
		response.Message = fmt.Sprintf("Error while parsing date pass date in ISO format: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	invokeArgs := util.ToChaincodeArgs("GetVerifier", invoker)
	verifierAsBytes := ctx.GetStub().InvokeChaincode("akcess", invokeArgs, "akcessglobal")
	if verifierAsBytes.Payload == nil {
		response.Message = fmt.Sprintf("Verifier %s not registered on global channel", invoker)
		logger.Info(response.Message)
		return response
	}
	var verifier Verifier
	json.Unmarshal(verifierAsBytes.Payload, &verifier)

	var eform Eform
	json.Unmarshal(eformAsBytes, &eform)

	verification := Verification{
		VerifierObj: verifier,
		ExpirtyDate: expirydate,
	}

	verifierList := VerifiersList(eform.Verifications)
	_, found := Find(verifierList, invoker)
	if found {
		for i, v := range eform.Verifications {
			if v.VerifierObj.AkcessID == invoker {
				eform.Verifications[i].ExpirtyDate = expirydate
				break
			}
		}
	} else {
		eform.Verifications = append(eform.Verifications, verification)
	}

	eformAsBytes, _ = json.Marshal(eform)
	err = ctx.GetStub().PutState(eformid, eformAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while verifying eform: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Eform %s verified by %s", eformid, invoker)
	logger.Info(response.Message)
	return response
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
func (d *EformContract) GetVerifiersOfEform(ctx contractapi.TransactionContextInterface, eformid string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	eformAsBytes, err := ctx.GetStub().GetState(eformid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching eform from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if eformAsBytes == nil {
		response.Message = fmt.Sprintf("Eform with id %s doesn't exist", eformid)
		logger.Info(response.Message)
		return response
	}

	var eform Eform
	json.Unmarshal(eformAsBytes, &eform)

	response.Data = eform.Verifications
	response.Success = true
	response.Message = fmt.Sprintf("Successfully fetched verifications of eform %s", eformid)
	logger.Info(response.Message)
	return response
}

// GetSignature get signature by signature hash
func (d *EformContract) GetSignature(ctx contractapi.TransactionContextInterface, signHash string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

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

	resultIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		response.Message = fmt.Sprintf("Error while searching eform by signature: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	defer resultIterator.Close()

	result := []Eform{}
	for resultIterator.HasNext() {
		queryResponse, _ := resultIterator.Next()

		eform := new(Eform)
		_ = json.Unmarshal(queryResponse.Value, eform)
		result = append(result, *eform)
		fmt.Println(result)
	}

	response.Data = result
	response.Success = true
	response.Message = fmt.Sprintf("Successfully fetched all eform with signature %s", signHash)
	logger.Info(response.Message)
	return response
}
