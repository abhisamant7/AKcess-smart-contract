package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	usercontract := new(UserContract)
	usercontract.UnknownTransaction = UnknownTransactionHandler
	usercontract.Name = "usercontract"

	doccontract := new(DocContract)
	doccontract.UnknownTransaction = UnknownTransactionHandler
	doccontract.Name = "doccontract"

	cc, err := contractapi.NewChaincode(usercontract, doccontract)
	cc.DefaultContract = usercontract.GetName()

	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}
