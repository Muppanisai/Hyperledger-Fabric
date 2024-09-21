package main

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Asset struct
type Asset struct {
    DealerID    string  `json:"dealerID"`
    MSISDN      string  `json:"msisdn"`
    MPIN        string  `json:"mpin"`
    Balance     float64 `json:"balance"`
    Status      string  `json:"status"`
    TransAmount float64 `json:"transAmount"`
    TransType   string  `json:"transType"`
    Remarks     string  `json:"remarks"`
}

// SmartContract provides functions for managing assets
type SmartContract struct {
    contractapi.Contract
}

// CreateAsset creates a new asset
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerID string, msisdn string, mpin string, balance float64, status string) error {
    asset := Asset{
        DealerID: dealerID,
        MSISDN:   msisdn,
        MPIN:     mpin,
        Balance:  balance,
        Status:   status,
    }

    assetAsBytes, err := json.Marshal(asset)
    if err != nil {
        return fmt.Errorf("failed to marshal asset: %v", err)
    }

    return ctx.GetStub().PutState(msisdn, assetAsBytes)
}

// UpdateBalance updates the balance of the asset
func (s *SmartContract) UpdateBalance(ctx contractapi.TransactionContextInterface, msisdn string, transAmount float64, transType string, remarks string) error {
    assetAsBytes, err := ctx.GetStub().GetState(msisdn)
    if err != nil {
        return fmt.Errorf("failed to read asset: %v", err)
    }
    if assetAsBytes == nil {
        return fmt.Errorf("asset not found: %s", msisdn)
    }

    asset := new(Asset)
    err = json.Unmarshal(assetAsBytes, asset)
    if err != nil {
        return fmt.Errorf("failed to unmarshal asset: %v", err)
    }

    // Update balance based on transaction type
    if transType == "debit" {
        asset.Balance -= transAmount
    } else if transType == "credit" {
        asset.Balance += transAmount
    } else {
        return fmt.Errorf("invalid transaction type: %s", transType)
    }

    asset.TransAmount = transAmount
    asset.TransType = transType
    asset.Remarks = remarks

    updatedAssetAsBytes, err := json.Marshal(asset)
    if err != nil {
        return fmt.Errorf("failed to marshal updated asset: %v", err)
    }

    return ctx.GetStub().PutState(msisdn, updatedAssetAsBytes)
}

// QueryAsset returns the asset stored in the world state with the given MSISDN
func (s *SmartContract) QueryAsset(ctx contractapi.TransactionContextInterface, msisdn string) (*Asset, error) {
    assetAsBytes, err := ctx.GetStub().GetState(msisdn)
    if err != nil {
        return nil, fmt.Errorf("failed to read asset: %v", err)
    }
    if assetAsBytes == nil {
        return nil, fmt.Errorf("asset not found: %s", msisdn)
    }

    asset := new(Asset)
    err = json.Unmarshal(assetAsBytes, asset)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
    }

    return asset, nil
}

// GetAssetHistory returns the transaction history of a given asset (MSISDN)
func (s *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, msisdn string) ([]*Asset, error) {
    resultsIterator, err := ctx.GetStub().GetHistoryForKey(msisdn)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve asset history: %v", err)
    }
    defer resultsIterator.Close()

    var history []*Asset
    for resultsIterator.HasNext() {
        record, err := resultsIterator.Next()
        if err != nil {
            return nil, fmt.Errorf("failed to read history record: %v", err)
        }

        var asset Asset
        json.Unmarshal(record.Value, &asset)
        history = append(history, &asset)
    }

    return history, nil
}

// main function to start the chaincode
func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        fmt.Printf("Error creating chaincode: %v", err)
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting chaincode: %v", err)
    }
}
