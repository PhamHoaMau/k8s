/*
 * SPDX-License-Identifier: Apache-2.0
 */

package iotrecord

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract chaincode that defines
// the business logic for managing commercial
// record
type Contract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (c *Contract) Instantiate(ctx TransactionContextInterface) {
	fmt.Println("Instantiated")
	t := time.Now()

	record1 := IoTRecord{MonitorID: "DANANG_1", Name: "TEMPERATURE", Unit: "^C", DateTime: t.Format("2006-01-02 15:04:05"), Value: 28.5}
	ctx.GetRecordList().AddRecord(&record1)

	record2 := IoTRecord{MonitorID: "DANANG_1", Name: "HUMIDITY", Unit: "%", DateTime: t.Format("2006-01-02 15:04:05"), Value: 60.2}
	ctx.GetRecordList().AddRecord(&record2)

	record3 := IoTRecord{MonitorID: "DANANG_1", Name: "PH", Unit: "%", DateTime: t.Format("2006-01-02 15:04:05"), Value: 2}
	ctx.GetRecordList().AddRecord(&record3)

	record4 := IoTRecord{MonitorID: "DANANG_1", Name: "UV", Unit: "%", DateTime: t.Format("2006-01-02 15:04:05"), Value: 2}
	ctx.GetRecordList().AddRecord(&record4)

	record5 := IoTRecord{MonitorID: "DANANG_1", Name: "Rain", Unit: "%", DateTime: t.Format("2006-01-02 15:04:05"), Value: 2}
	ctx.GetRecordList().AddRecord(&record5)

	record6 := IoTRecord{MonitorID: "QUANGNAM_1", Name: "TEMPERATURE", Unit: "^C", DateTime: t.Format("2006-01-02 15:04:05"), Value: 2}
	ctx.GetRecordList().AddRecord(&record6)

	record7 := IoTRecord{MonitorID: "QUANGNAM_1", Name: "HUMIDITY", Unit: "%", DateTime: t.Format("2006-01-02 15:04:05"), Value: 2}
	ctx.GetRecordList().AddRecord(&record7)

	record8 := IoTRecord{MonitorID: "QUANGNAM_1", Name: "PH", Unit: "%", DateTime: t.Format("2006-01-02 15:04:05"), Value: 2}
	ctx.GetRecordList().AddRecord(&record8)

	record9 := IoTRecord{MonitorID: "QUANGNAM_1", Name: "UV", Unit: "%", DateTime: t.Format("2006-01-02 15:04:05"), Value: 2}
	ctx.GetRecordList().AddRecord(&record9)

	record10 := IoTRecord{MonitorID: "QUANGNAM_1", Name: "Rain", Unit: "%", DateTime: t.Format("2006-01-02 15:04:05"), Value: 2}
	ctx.GetRecordList().AddRecord(&record10)

}

// UpdateRecord to save record
func (c *Contract) UpdateRecord(ctx TransactionContextInterface, monitorID string, name string, unit string, value float32, datetime string) error {
	fmt.Println("SaveRecord")
	record := IoTRecord{MonitorID: strings.ToUpper(monitorID), Name: strings.ToUpper(name), Unit: unit, DateTime: datetime, Value: value}
	err := ctx.GetRecordList().UpdateRecord(&record)
	return err
}

// QueryRecordHistoryByTimeRange query record by time range
func (c *Contract) QueryRecordHistoryByTimeRange(ctx TransactionContextInterface, monitorID string, name string, start string, end string) ([]IoTRecord, error) {
	payload, err := ctx.GetRecordList().GetHistoryForRecordByTimeRange(strings.ToUpper(monitorID), strings.ToUpper(name), start, end)
	if err != nil {
		return nil, err
	}
	var record []IoTRecord
	err = json.Unmarshal(payload, &record)
	fmt.Println("QueryRecordHistoryByTimeRange Result:")

	fmt.Println(record)
	// s := string(payload)
	return record, err
}

// QueryRecordHistory query record by time range
func (c *Contract) QueryRecordHistory(ctx TransactionContextInterface, monitorID string, name string) ([]IoTRecord, error) {
	payload, err := ctx.GetRecordList().GetHistoryForRecord(strings.ToUpper(monitorID), strings.ToUpper(name))
	if err != nil {
		return nil, err
	}
	var record []IoTRecord
	err = json.Unmarshal(payload, &record)
	fmt.Println("QueryRecordHistory Result:")

	fmt.Println(record)
	// s := string(payload)
	return record, err
}
