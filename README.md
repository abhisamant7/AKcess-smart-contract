# AKcess-smart-contract
AKcess document verification

## Smart contract with its arguments

## User smart contract

#### Create User
peer chaincode invoke -n akcess -C myc -c '{"Args":["CreateUser","AK1"]}'

#### Create User
peer chaincode invoke -n akcess -C myc -c '{"Args":["CreateVerifier","V1","verifier1","A"]}'

#### Add user profile verification
peer chaincode invoke -n akcess -C myc -c '{"Args":["AddUserProfileVerification","V1","AK1","mobile","2020-07-23T16:44:46.503Z"]}'

#### Get verifiers of user profile
peer chaincode invoke -n akcess -C myc -c '{"Args":["GetVerifiersOfUserProfile","AK1","mobile"]}'

#### Delete verfication from user profile
peer chaincode invoke -n akcess -C myc -c '{"Args":["DeleteVerification","AK1","mobile","V1"]}'


## Document smart contract

#### Create doc
peer chaincode invoke -n akcess -C myc -c '{"Args":["doccontract:CreateDoc","AK1","doc1","[\"skdjaskdjas\",\"kskdaksjdsk\"]"]}'

#### Sign doc
peer chaincode invoke -n akcess -C myc -c '{"Args":["doccontract:SignDoc","AK1","doc1","sjhsdjaksdjak","2020-07-23T16:44:46.503Z","123456"]}'

#### Share doc
peer chaincode invoke -n akcess -C myc -c '{"Args":["doccontract:SendDoc","AK1","uuid","AK4","asd"]}'

#### Verify doc
peer chaincode invoke -n akcess -C myc -c '{"Args":["doccontract:VerifyDoc","V1","soc1","2020-07-23T16:55:32.735Z"]}'

#### Get verifiers of doc
peer chaincode invoke -n akcess -C myc -c '{"Args":["doccontract:GetVerifiersOfDoc","doc1"]}'

#### Get signature
peer chaincode invoke -n akcess -C myc -c '{"Args":["doccontract:GetSignature","sjhsdjaksdjak"]}'


## Eform smart contract

#### Create eform
peer chaincode invoke -n eform -C myc -c '{"Args":["CreateEform","AK1","eform1","[\"skdjaskdjas\",\"kskdaksjdsk\"]"]}'

#### Sign eform
peer chaincode invoke -n eform -C myc -c '{"Args":["SignEform","AK1","eform1","sjhsdjaksdja","2020-07-23T16:44:46.503Z","123456"]}'

#### Share eform
peer chaincode invoke -n eform -C myc -c '{"Args":["SendEform","AK1","uuid_","AK4","eform1"]}'

#### Verify eform
peer chaincode invoke -n eform -C myc -c '{"Args":["VerifyEform","V1","eform1","2020-07-23T16:55:32.735Z"]}'

#### Get verifiers of eform
peer chaincode invoke -n eform -C myc -c '{"Args":["GetVerifiersOfEform","eform1"]}'

#### Get signature of eform
peer chaincode invoke -n eform -C myc -c '{"Args":["GetSignature","sjhsdjaksdja"]}'
