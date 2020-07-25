package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	eformcontract := new(EformContract)
	eformcontract.UnknownTransaction = UnknownTransactionHandler
	eformcontract.Name = "eformcontract"

	cc, err := contractapi.NewChaincode(eformcontract)
	cc.DefaultContract = eformcontract.GetName()

	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}
