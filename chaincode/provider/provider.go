package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for control the food
type SmartContract struct {
	contractapi.Contract
}

//Mask describes basic details
type Mask struct {
	Type  string `json:"type"`
	Code string `json:"code"`
	Madeby string `json:"madeby"`
	Owner string `json:"owner"`
}

// From Laurent question
type MaskTracking struct {
	Code  string `json:"code"`
	Timestamp string `json:"timestamp"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	return nil
	/*
	masks := []Mask{
		{Type: "FP2", Code: "A123", Madeby: "Spain", Owner: "Provider1"},
		{Type: "FP2", Code: "B123", Madeby: "Germany", Owner: "Provider2"},
		Type: "FP2", Code: "C123", Madeby: "Spain", Owner: "Provider3"}
	}

	for _, mask := range masks{
		err := c.Set(ctx, mask.Type, mask.Code, mask.Madeby, mask.Owner)
		if err != nil {
			return nil
		}
	}
	*/
}

// Set or Update
func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, maskId string, typeM string, madeBy string, owner string, code string) error {


	// validate parameters if we dont want to update

	mask := Mask {
		Type:  typeM,
		Code: code,
		Madeby: madeBy,
		Owner: owner,
	}

	maskAsBytes, err := json.Marshal(mask)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}

	return ctx.GetStub().PutState(maskId, maskAsBytes)
}

/*
func (s *SmartContract) changeMaskOwner(ctx contractapi.TransactionContextInterface, maskId string, newOwner string) error {

	// Get current mask state
	maskObject, err := Query(maskId);

	maskObject.Owner = newOwner;

	maskAsBytes, err := json.Marshal(maskObject)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}

	return ctx.GetStub().PutState(maskId, maskAsBytes)
}
*/

// Use couchdb to query history
/*
func (s *SmartContract) Track(ctx contractapi.TransactionContextInterface, maskId string) (*Mask, error) {

	// Get the history of a mask
	i := 1032049348
    fmt.Println(i)
}
*/

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, maskId string) (*Mask, error) {

	maskAsBytes, err := ctx.GetStub().GetState(maskId)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if maskAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", maskId)
	}

	mask := new(Mask)

	err = json.Unmarshal(maskAsBytes, mask)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error. %s", err.Error())
	}

	return mask, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create provider chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting provider chaincode: %s", err.Error())
	}
}
