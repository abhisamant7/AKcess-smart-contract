package main

import (
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

var logger = flogging.MustGetLogger("akcess")

func main() {
	usercontract := new(UserContract)
	usercontract.UnknownTransaction = UnknownTransactionHandler
	usercontract.Name = "usercontract"

	doccontract := new(DocContract)
	doccontract.UnknownTransaction = UnknownTransactionHandler
	doccontract.Name = "doccontract"

	assetContract := new(DigitalAssetContract)
	assetContract.UnknownTransaction = UnknownTransactionHandler
	assetContract.Name = "adat"

	cc, err := contractapi.NewChaincode(usercontract, doccontract, assetContract)
	cc.DefaultContract = usercontract.GetName()

	if err != nil {
		panic(err.Error())
	}

	if os.Getenv("ISEXTERNAL") == "true" {
		server := &shim.ChaincodeServer{
			CCID:    os.Getenv("CHAINCODE_CCID"),
			Address: os.Getenv("CHAINCODE_ADDRESS"),
			CC:      cc,
			TLSProps: shim.TLSProperties{
				Disabled: true,
			},
		}

		if err := server.Start(); err != nil {
			panic(err.Error())
		}
	} else {
		if err := cc.Start(); err != nil {
			panic(err.Error())
		}
	}
}
