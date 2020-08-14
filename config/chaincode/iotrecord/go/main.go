/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"

	iotrecord "jwclab/iotrecord/iot-record"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	contract := new(iotrecord.Contract)
	contract.TransactionContextHandler = new(iotrecord.TransactionContext)
	contract.Name = "org.jwclab.iotrecord"
	contract.Info.Version = "0.0.1"

	chaincode, err := contractapi.NewChaincode(contract)

	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode. %s", err.Error()))
	}

	chaincode.Info.Title = "IoTRecord"
	chaincode.Info.Version = "0.0.1"

	err = chaincode.Start()

	if err != nil {
		panic(fmt.Sprintf("Error starting chaincode. %s", err.Error()))
	}
}
