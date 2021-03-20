package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

var logger = flogging.MustGetLogger("example")

// DigitalAssetContract Smart contract for AKcess digital asset token
type DigitalAssetContract struct {
	contractapi.Contract
}

// RegisterAsset register new digital asset
func (da *DigitalAssetContract) RegisterAsset(ctx contractapi.TransactionContextInterface, assetType string, metadata map[string]string, description string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, err := getCommonName(ctx)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting commonName from x509 cert: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	asset := DigitalAsset{
		UniqueAssetID: response.TxID,
		AssetType:     assetType,
		Owner:         *invoker,
		Metadata:      metadata,
		Description:   description,
	}

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		response.Message = fmt.Sprintf("Error while marshling asset: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	err = ctx.GetStub().PutState(asset.UniqueAssetID, assetAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while saving asset in ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Digital asset sucessfully saved with id %s", asset.UniqueAssetID)
	logger.Info(response.Message)
	response.Data = asset
	return response
}

// TransferAsset transfers given asset from invoker to recipient
func (da *DigitalAssetContract) TransferAsset(ctx contractapi.TransactionContextInterface, assetID string, recipient string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, err := getCommonName(ctx)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting commonName from x509 cert: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	assetAsBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting asset from ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}
	if assetAsBytes == nil {
		response.Message = fmt.Sprintf("Digital asset with id %s not found", assetID)
		logger.Info(response.Message)
		return response
	}

	var asset DigitalAsset
	err = json.Unmarshal(assetAsBytes, &asset)
	if err != nil {
		response.Message = fmt.Sprintf("Error while unmarshling asset: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	if asset.Owner != *invoker {
		response.Message = fmt.Sprintf("Digtal asset with id %s not owned by %s", asset.UniqueAssetID, *invoker)
		logger.Error(response.Message)
		return response
	}

	asset.Owner = recipient
	assetAsBytes, err = json.Marshal(asset)
	if err != nil {
		response.Message = fmt.Sprintf("Error while marshling asset: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	err = ctx.GetStub().PutState(asset.UniqueAssetID, assetAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while saving asset in ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Digital asset with id %s successfully updated and owned by %s", asset.UniqueAssetID, recipient)
	logger.Info(response.Message)
	response.Data = asset
	return response
}

// LinkDocument link document to digital asset
func (da *DigitalAssetContract) LinkDocument(ctx contractapi.TransactionContextInterface, assetID string, documentID string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, err := getCommonName(ctx)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting commonName from x509 cert: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	assetAsBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting asset from ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}
	if assetAsBytes == nil {
		response.Message = fmt.Sprintf("Digital asset with id %s not found", assetID)
		logger.Info(response.Message)
		return response
	}

	var asset DigitalAsset
	err = json.Unmarshal(assetAsBytes, &asset)
	if err != nil {
		response.Message = fmt.Sprintf("Error while unmarshling asset: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	if asset.Owner != *invoker {
		response.Message = fmt.Sprintf("Digtal asset with id %s not owned by %s", asset.UniqueAssetID, *invoker)
		logger.Error(response.Message)
		return response
	}

	_, found := Find(asset.LinkedDocs, documentID)
	if found {
		response.Message = fmt.Sprintf("Document %s already linked with asset %s", documentID, assetID)
		logger.Error(response.Message)
		return response
	} else {
		asset.LinkedDocs = append(asset.LinkedDocs, documentID)
	}

	assetAsBytes, err = json.Marshal(asset)
	if err != nil {
		response.Message = fmt.Sprintf("Error while marshling asset: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	err = ctx.GetStub().PutState(asset.UniqueAssetID, assetAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while saving asset in ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Document %s linked with asset %s", documentID, assetID)
	logger.Info(response.Message)
	response.Data = asset
	return response
}

// VerifyAssetOwnership verifiers can verify the ownership of asset holders
func (da *DigitalAssetContract) VerifyAssetOwnership(ctx contractapi.TransactionContextInterface, assetID string, expiryDate string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, err := getCommonName(ctx)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting commonName from x509 cert: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	verifierAsBytes, err := ctx.GetStub().GetState(*invoker)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting verifier from ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}
	if verifierAsBytes == nil {
		response.Message = fmt.Sprintf("Verifier with id %s doesn't exist", *invoker)
		logger.Info(response.Message)
		return response
	}
	var verifier Verifier
	err = json.Unmarshal(verifierAsBytes, &verifier)
	if err != nil {
		response.Message = fmt.Sprintf("Error while unmarshling verifier: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	assetAsBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting asset from ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}
	if assetAsBytes == nil {
		response.Message = fmt.Sprintf("Digital asset with id %s not found", assetID)
		logger.Info(response.Message)
		return response
	}

	var asset DigitalAsset
	err = json.Unmarshal(assetAsBytes, &asset)
	if err != nil {
		response.Message = fmt.Sprintf("Error while unmarshling asset: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	verifierList := VerifiersList(asset.Verifications)
	_, found := Find(verifierList, *invoker)
	if found {
		response.Message = fmt.Sprintf("Digital asset %s already verified by %s", assetID, *invoker)
		logger.Error(response.Message)
		return response
	} else {
		expirydate, err := time.Parse(time.RFC3339, expiryDate)
		if err != nil {
			response.Message = fmt.Sprintf("Error while parsing date, please pass in ISO format: %s", err.Error())
			logger.Error(response.Message)
			return response
		}
		verification := Verification{
			VerifierObj: verifier,
			ExpirtyDate: expirydate,
		}
		asset.Verifications = append(asset.Verifications, verification)
	}

	assetAsBytes, err = json.Marshal(asset)
	if err != nil {
		response.Message = fmt.Sprintf("Error while marshling asset: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	err = ctx.GetStub().PutState(asset.UniqueAssetID, assetAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while saving asset in ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Digital asset %s verified by %s", assetID, *invoker)
	logger.Info(response.Message)
	response.Data = asset
	return response
}

// RemoveVerification verifiers can remove their verification from asset
func (da *DigitalAssetContract) RemoveVerification(ctx contractapi.TransactionContextInterface, assetID string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, err := getCommonName(ctx)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting commonName from x509 cert: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	verifierAsBytes, err := ctx.GetStub().GetState(*invoker)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting verifier from ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}
	if verifierAsBytes == nil {
		response.Message = fmt.Sprintf("Verifier with id %s doesn't exist", *invoker)
		logger.Info(response.Message)
		return response
	}
	var verifier Verifier
	err = json.Unmarshal(verifierAsBytes, &verifier)
	if err != nil {
		response.Message = fmt.Sprintf("Error while unmarshling verifier: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	assetAsBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting asset from ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}
	if assetAsBytes == nil {
		response.Message = fmt.Sprintf("Digital asset with id %s not found", assetID)
		logger.Info(response.Message)
		return response
	}

	var asset DigitalAsset
	err = json.Unmarshal(assetAsBytes, &asset)
	if err != nil {
		response.Message = fmt.Sprintf("Error while unmarshling asset: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	verifierList := VerifiersList(asset.Verifications)
	index, found := Find(verifierList, *invoker)
	if !found {
		response.Message = fmt.Sprintf("Verifier %s didn't make any verification on asset %s", *invoker, *invoker)
		logger.Error(response.Message)
		return response
	} else {
		Remove(asset.Verifications, index)
	}

	assetAsBytes, err = json.Marshal(asset)
	if err != nil {
		response.Message = fmt.Sprintf("Error while marshling asset: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	err = ctx.GetStub().PutState(asset.UniqueAssetID, assetAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while saving asset in ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Verification of verifier %s removed from asset %s", *invoker, assetID)
	logger.Info(response.Message)
	response.Data = asset
	return response

}

// GetDigitalAsset returns asset with all the details it's inked documents and verifications
func (da *DigitalAssetContract) GetDigitalAsset(ctx contractapi.TransactionContextInterface, assetID string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	assetAsBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		response.Message = fmt.Sprintf("Error while getting asset from ledger: %s", err.Error())
		logger.Error(response.Message)
		return response
	}
	if assetAsBytes == nil {
		response.Message = fmt.Sprintf("Digital asset with id %s not found", assetID)
		logger.Info(response.Message)
		return response
	}

	var asset DigitalAsset
	err = json.Unmarshal(assetAsBytes, &asset)
	if err != nil {
		response.Message = fmt.Sprintf("Error while unmarshling asset: %s", err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprint("Successfully fetched asset")
	logger.Info(response.Message)
	response.Data = asset
	return response
}
