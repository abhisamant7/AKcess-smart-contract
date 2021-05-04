package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// UserContract contract for storing user in blockchain
type UserContract struct {
	contractapi.Contract
}

// CreateUser adds a new user to the world state with given details
func (u *UserContract) CreateUser(ctx contractapi.TransactionContextInterface) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, _ := getCommonName(ctx)
	userAsBytes, err := ctx.GetStub().GetState(invoker)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching user from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if userAsBytes != nil {
		response.Message = fmt.Sprintf("AKcessID %s already exist", invoker)
		logger.Info(response.Message)
		return response
	}

	user := User{
		ObjectType:    "user",
		AkcessID:      invoker,
		Verifications: map[string][]Verification{},
	}
	newUserAsBytes, _ := json.Marshal(user)
	err = ctx.GetStub().PutState(invoker, newUserAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while registering user: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("User with AKcessID %s added\n", invoker)
	logger.Info(response.Message)
	return response
}

// CreateVerifier register new verifier in Blockchain
func (u *UserContract) CreateVerifier(ctx contractapi.TransactionContextInterface, verifierName string, VerifierGrade string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, _ := getCommonName(ctx)
	verifierAsBytes, err := ctx.GetStub().GetState(invoker)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching verifier from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if verifierAsBytes != nil {
		response.Message = fmt.Sprintf("AKcessID %s already exist", invoker)
		logger.Info(response.Message)
		return response
	}

	verifier := Verifier{
		ObjectType:    "verifier",
		AkcessID:      invoker,
		VerifierName:  verifierName,
		VerifierGrade: VerifierGrade,
	}
	newVerifierAsBytes, _ := json.Marshal(verifier)
	err = ctx.GetStub().PutState(invoker, newVerifierAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while registering verifier: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Verifier with AKcessID %s added\n", invoker)
	logger.Info(response.Message)
	return response
}

// AddUserProfileVerification add verifcation transaction and field of users profiles is verfiied
func (u *UserContract) AddUserProfileVerification(ctx contractapi.TransactionContextInterface, verifierAKcessID string, userAKcessID string, profileFields []string, expiryDates []string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, _ := getCommonName(ctx)
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

	userAsBytes, err := ctx.GetStub().GetState(userAKcessID)
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
	var user User
	json.Unmarshal(userAsBytes, &user)

	for index, profileField := range profileFields {
		verifierList := VerifiersList(user.Verifications[profileField])
		expirydate, err := time.Parse(time.RFC3339, expiryDates[index])
		if err != nil {
			response.Message = fmt.Sprintf("Error while parsing date pass date in ISO format: %s", err.Error())
			logger.Info(response.Message)
			return response
		}

		_, found := Find(verifierList, verifierAKcessID)
		if found {
			for i, v := range user.Verifications[profileField] {
				if v.VerifierObj.AkcessID == verifierAKcessID {
					user.Verifications[profileField][i].ExpirtyDate = expirydate
					break
				}
			}
		} else {
			verification := Verification{
				VerifierObj: verifier,
				ExpirtyDate: expirydate,
			}
			user.Verifications[profileField] = append(user.Verifications[profileField], verification)
		}
	}

	userAsBytes, _ = json.Marshal(user)
	err = ctx.GetStub().PutState(userAKcessID, userAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while updating user profile verification: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Profile field %s of user %s verified by %s", profileFields, userAKcessID, verifierAKcessID)
	logger.Info(response.Message)
	return response
}

// GetVerifiersOfUserProfile get verifiers of perticular user field
func (u *UserContract) GetVerifiersOfUserProfile(ctx contractapi.TransactionContextInterface, akcessid string, profileField string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	userAsBytes, err := ctx.GetStub().GetState(akcessid)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching user from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if userAsBytes == nil {
		response.Message = fmt.Sprintf("AKcessID %s doesn't exist", akcessid)
		logger.Info(response.Message)
		return response
	}

	var user User
	json.Unmarshal(userAsBytes, &user)

	response.Data = user.Verifications
	response.Success = true
	response.Message = fmt.Sprintf("Succesfully fetched verfiers list of user profile field %s", profileField)
	logger.Info(response.Message)
	return response
}

// GetVerifier get verifier
func (u *UserContract) GetVerifier(ctx contractapi.TransactionContextInterface, akcessid string) (*Verifier, error) {
	verifierAsBytes, err := ctx.GetStub().GetState(akcessid)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if verifierAsBytes == nil {
		return nil, fmt.Errorf("AKcessID %s doesn't exist", akcessid)
	}

	var verifier Verifier
	json.Unmarshal(verifierAsBytes, &verifier)

	return &verifier, nil
}

// DeleteVerification deletes the verification from user profile
func (u *UserContract) DeleteVerification(ctx contractapi.TransactionContextInterface, profileField string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	invoker, _ := getCommonName(ctx)
	userAsBytes, err := ctx.GetStub().GetState(invoker)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching user from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if userAsBytes == nil {
		response.Message = fmt.Sprintf("AKcessID %s doesn't exist", invoker)
		logger.Info(response.Message)
		return response
	}

	var user User
	json.Unmarshal(userAsBytes, &user)
	user.Verifications[profileField] = []Verification{}
	newUserAsBytes, _ := json.Marshal(user)
	err = ctx.GetStub().PutState(invoker, newUserAsBytes)
	if err != nil {
		response.Message = fmt.Sprintf("Error while deleting profile field verification: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Deleted verification of %s's %s profile field", invoker, profileField)
	logger.Info(response.Message)
	return response
}

// DeleteUser deletes the user from Blockchain world state
func (u *UserContract) DeleteUser(ctx contractapi.TransactionContextInterface, key string) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	userAsBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching data from world state: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}
	if userAsBytes == nil {
		response.Message = fmt.Sprintf("Key %s doesn't exist", key)
		logger.Info(response.Message)
		return response
	}

	err = ctx.GetStub().DelState(key)
	if err != nil {
		response.Message = fmt.Sprintf("Error while deleting data: %s" + err.Error())
		logger.Error(response.Message)
		return response
	}

	response.Success = true
	response.Message = fmt.Sprintf("Key %s deleted", key)
	logger.Info(response.Message)
	return response
}

// GetAllVerifiers returns all registered verifiers
func (u *UserContract) GetAllVerifiers(ctx contractapi.TransactionContextInterface) Response {
	response := Response{
		TxID:    ctx.GetStub().GetTxID(),
		Success: false,
		Message: "",
		Data:    nil,
	}

	var richQuery string = `{
		"selector": {
		   "docType": "verifier"
		}
	}`
	resultIterator, err := ctx.GetStub().GetQueryResult(richQuery)
	if err != nil {
		response.Message = fmt.Sprintf("Error while fetching query result: %s", err.Error())
		logger.Error(response.Message)
		return response
	}
	defer resultIterator.Close()

	var result []Verifier
	for resultIterator.HasNext() {
		queryResponse, _ := resultIterator.Next()

		v := new(Verifier)
		_ = json.Unmarshal(queryResponse.Value, v)
		result = append(result, *v)
	}

	response.Success = true
	response.Message = fmt.Sprint("Successfully fetched all verifiers")
	logger.Info(response.Message)
	response.Data = result
	return response
}
