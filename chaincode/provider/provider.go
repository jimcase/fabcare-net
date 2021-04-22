package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for control the food
type SmartContract struct {
	contractapi.Contract
}

//Mask describes basic details
type Mask struct {
	Type   string `json:"type"`
	Code   string `json:"code"`
	Madeby string `json:"madeby"`
	Owner  string `json:"owner"`
}

// From Laurent question
type MaskTx struct {
	TxId      string    `json:"tx"`
	PrevOwner string    `json:"prevOwner"`
	NewOwner  string    `json:"newOwner"`
	Timestamp time.Time `json:"timestamp"`
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

// Create a new mask
func (s *SmartContract) CreateMask(ctx contractapi.TransactionContextInterface, maskId string, typeM string, madeBy string, owner string, code string) error {

	// validate parameters if we dont want to update

	exists, err := s.MaskExists(ctx, maskId)

	if err != nil {
		return err
	}

	if !exists {
		mask := Mask{
			Type:   typeM,
			Code:   code,
			Madeby: madeBy,
			Owner:  owner,
		}

		maskAsBytes, err := json.Marshal(mask)
		if err != nil {
			fmt.Printf("Marshal error: %s", err.Error())
			return err
		}
		return ctx.GetStub().PutState(maskId, maskAsBytes)
	} else {
		return fmt.Errorf("The mask %s  already exists", maskId)
	}
}

// Create a new mask
func (s *SmartContract) UpdateMask(ctx contractapi.TransactionContextInterface, maskId string, typeM string, madeBy string, owner string, code string) error {

	// validate parameters if we dont want to update

	exists, err := s.MaskExists(ctx, maskId)

	if err != nil {
		return err
	}

	if exists {
		mask := Mask{
			Type:   typeM,
			Code:   code,
			Madeby: madeBy,
			Owner:  owner,
		}

		maskAsBytes, err := json.Marshal(mask)
		if err != nil {
			fmt.Printf("Marshal error: %s", err.Error())
			return err
		}
		return ctx.GetStub().PutState(maskId, maskAsBytes)
	} else {
		return fmt.Errorf("The mask %s  does not exist", maskId)
	}
}

func (s *SmartContract) GetMask(ctx contractapi.TransactionContextInterface, maskId string) (*Mask, error) {

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

func (s *SmartContract) DeleteMask(ctx contractapi.TransactionContextInterface, maskId string) error {

	exists, err := s.MaskExists(ctx, maskId)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("The mask %s does not exist", maskId)
	}

	return ctx.GetStub().DelState(maskId)
}

// MaskExists returns true when asset with given ID exists in world state
func (s *SmartContract) MaskExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) ChangeMaskOwner(ctx contractapi.TransactionContextInterface, maskId string, newOwner string) (string, error) {
	mask, err := s.GetMask(ctx, maskId)

	if err != nil {
		return "error0", err
	}

	txid := ctx.GetStub().GetTxID()

	maskTx := MaskTx{
		TxId:      txid,
		PrevOwner: mask.Owner,
		NewOwner:  newOwner,
		Timestamp: time.Now(),
	}
	txBytes, err := json.Marshal(maskTx)
	if err != nil {
		return "error2", err
	}

	mask.Owner = newOwner
	maskAsBytes, err := json.Marshal(mask)
	if err != nil {
		return "error1", err
	}

	ctx.GetStub().PutState(maskId, maskAsBytes)
	ctx.GetStub().PutState(txid, txBytes)

	return txid, err
}

func (s *SmartContract) GetAllMasks(ctx contractapi.TransactionContextInterface) ([]*Mask, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all masks in the chaincode namespace.

	// Get iterator object
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var masks []*Mask
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var mask Mask
		err = json.Unmarshal(queryResponse.Value, &mask)
		if err != nil {
			return nil, err
		}
		masks = append(masks, &mask)
	}

	return masks, nil
}

func (s *SmartContract) GetTotalMasks(ctx contractapi.TransactionContextInterface) (int, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return 0, err
	}
	defer resultsIterator.Close()

	var counter int
	counter = 0
	for resultsIterator.HasNext() {
		counter = counter + 1
		_, err := resultsIterator.Next()
		if err != nil {
			return 0, err
		}
	}

	return counter, nil
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
