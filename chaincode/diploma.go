package main

import (
  "fmt"
  "encoding/json"
  "log"
  "strconv"
  "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
  contractapi.Contract
}

type Asset struct {
  University      string `json:"University"`
  ID              string `json:"ID"`
  StudentName     string `json:"StudentName"`
  Credit          int    `json:"Credit"`
  GradStatus      bool   `json:"GradStatus"`
}

func main() {
  assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
  if err != nil {
    log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
  }

  if err := assetChaincode.Start(); err != nil {
    log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
  }
}



// ----------------------------------------------------------------------------------

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
  assets := []Asset{
    { University:"Genesis1", ID: "1", StudentName:"Genesis Name1", Credit:0, GradStatus:false},
    { University:"Genesis2", ID: "2", StudentName:"Genesis Name2", Credit:0, GradStatus:false},
    { University:"Genesis3", ID: "3", StudentName:"Genesis Name3", Credit:0, GradStatus:false},
  }

  for _, asset := range assets {
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    err = ctx.GetStub().PutState(asset.ID, assetJSON)
    if err != nil {
      return fmt.Errorf("failed to put to world state. %v", err)
    }
  }

  return nil
}

func (cc *SmartContract) RequireConsent(ctx contractapi.TransactionContextInterface, receiver string) (bool, error) {
	//Org := receiver // Specify the recipient organization's name

	// Invoke the consent chaincode on the recipient organization's peer
	response1:= ctx.GetStub().InvokeChaincode("basic", [][]byte{[]byte("checkConsent")},"myChannel")

	consent, err := strconv.ParseBool(string(response1.Payload))
	if err != nil {
		return false, fmt.Errorf("Failed to parse consent value: %v", err)
	}

  print(consent)
  
  return consent,nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, university string, id string, studentName string, credit int, gradStatus bool) error {
  // TODO: Implement simulation logic here
  exists, err := s.AssetExists(ctx, id)
  if err != nil {
    return err
  }
  if exists {
    return fmt.Errorf("the asset %s already exists", id)
  }

  asset := Asset{
    University :    university,
    ID:             id,
    StudentName:    studentName,
    Credit:         credit,
    GradStatus:     gradStatus,
  }
  assetJSON, err := json.Marshal(asset)
  if err != nil {
    return err
  }

  return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
  assetJSON, err := ctx.GetStub().GetState(id)
  if err != nil {
    return nil, fmt.Errorf("failed to read from world state: %v", err)
  }
  if assetJSON == nil {
    return nil, fmt.Errorf("the asset %s does not exist", id)
  }

  var asset Asset
  err = json.Unmarshal(assetJSON, &asset)
  if err != nil {
    return nil, err
  }

  return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, university string, id string, studentName string, credit int, gradStatus bool) error {
  exists, err := s.AssetExists(ctx, id)
  if err != nil {
    return err
  }
  if !exists {
    return fmt.Errorf("the asset %s does not exist", id)
  }

  // overwriting original asset with new asset
  asset := Asset{
    University :    university,
    ID:             id,
    StudentName:    studentName,
    Credit:         credit,
    GradStatus:     gradStatus,
  }
  assetJSON, err := json.Marshal(asset)
  if err != nil {
    return err
  }

  return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
  
  exists, err := s.AssetExists(ctx, id)
  if err != nil {
    return err
  }
  if !exists {
    return fmt.Errorf("the asset %s does not exist", id)
  }

  return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {

  assetJSON, err := ctx.GetStub().GetState(id)
  if err != nil {
    return false, fmt.Errorf("failed to read from world state: %v", err)
  }

  return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newCredits int, newGradStatus bool) error {
  asset, err := s.ReadAsset(ctx, id)
  if err != nil {
    return err
  }

  asset.Credit = newCredits
  asset.GradStatus= newGradStatus

  assetJSON, err := json.Marshal(asset)
  if err != nil {
    return err
  }
  err=ctx.GetStub().PutState(id, assetJSON)
  if err != nil {
		return fmt.Errorf("failed to put state: %v", err)
	}
  return nil 
}


// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
  resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
  if err != nil {
    return nil, err
  }
  defer resultsIterator.Close()

  var assets []*Asset
  for resultsIterator.HasNext() {
    queryResponse, err := resultsIterator.Next()
    if err != nil {
      return nil, err
    }

    var asset Asset
    err = json.Unmarshal(queryResponse.Value, &asset)
    if err != nil {
      return nil, err
    }
    assets = append(assets, &asset)
  }

  return assets, nil
}
