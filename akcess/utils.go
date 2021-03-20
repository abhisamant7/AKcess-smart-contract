package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Response chaincode response will be returned in this format
type Response struct {
	TxID    string      `json:"txId"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// UnknownTransactionHandler returns a shim error
// with details of a bad transaction request
func UnknownTransactionHandler(ctx contractapi.TransactionContextInterface) error {
	fcn, args := ctx.GetStub().GetFunctionAndParameters()
	return fmt.Errorf("Invalid function %s passed with args %v", fcn, args)
}

func getCommonName(ctx contractapi.TransactionContextInterface) (*string, error) {
	x509, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return nil, err
	}
	return &x509.Subject.CommonName, nil
}
