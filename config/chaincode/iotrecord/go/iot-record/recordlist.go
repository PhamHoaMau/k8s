/*
 * SPDX-License-Identifier: Apache-2.0
 */

package iotrecord

import (
	"encoding/json"
	"errors"
	ledgerapi "jwclab/iotrecord/ledger-api"
)

// ListInterface defines functionality needed
// to interact with the world state on behalf
// of a commercial record
type ListInterface interface {
	AddRecord(*IoTRecord) error
	UpdateRecord(*IoTRecord) error
	GetHistoryForRecordByTimeRange(string, string, string, string) ([]byte, error)
	GetHistoryForRecord(string, string) ([]byte, error)
}

type list struct {
	stateList ledgerapi.StateListInterface
}

func (cpl *list) UpdateRecord(record *IoTRecord) error {
	return cpl.stateList.UpdateState(record)
}

func (cpl *list) AddRecord(record *IoTRecord) error {
	return cpl.stateList.AddState(record)
}

func (cpl *list) GetHistoryForRecord(monitorID string, name string) ([]byte, error) {
	response := cpl.stateList.GetStateHistory(CreateIoTRecordKey(monitorID, name))
	if response.GetStatus() == 500 {
		return nil, errors.New(response.GetMessage())
	}
	var record []IoTRecord
	payload := response.GetPayload()
	err := json.Unmarshal(payload, &record)
	if err != nil {
		return nil, err
	}
	return payload, err
}

func (cpl *list) GetHistoryForRecordByTimeRange(monitorID string, name string, start string, end string) ([]byte, error) {
	response := cpl.stateList.GetStateHistory(CreateIoTRecordKey(monitorID, name))
	if response.GetStatus() == 500 {
		return nil, errors.New(response.GetMessage())
	}
	var record []IoTRecord
	payload := response.GetPayload()
	err := json.Unmarshal(payload, &record)
	if err != nil {
		return nil, err
	}

	recordLen := len(record)
	datetime := func(r1, r2 *IoTRecord) bool {
		return r1.DateTime < r2.DateTime
	}
	By(datetime).Sort(record)
	var indexStart, indexEnd int
	for i := 0; i < recordLen-1; i++ {
		if record[i].DateTime < start && record[i+1].DateTime >= start {
			indexStart = i + 1
		}
		if record[i].DateTime <= end && record[i+1].DateTime > end {
			indexEnd = i + 1
		}
	}
	if record[0].DateTime >= start {
		indexStart = 0
	}
	if record[recordLen-1].DateTime <= end {
		indexEnd = recordLen
	}
	if start > end {
		return nil, errors.New("Invalid time range")
	}
	if end < record[0].DateTime || start > record[recordLen-1].DateTime {
		return nil, errors.New("History not found by range")
	}
	return json.Marshal(record[indexStart:indexEnd])
}

// NewList create a new list from context
func newList(ctx TransactionContextInterface) *list {
	stateList := new(ledgerapi.StateList)
	stateList.Ctx = ctx
	stateList.Name = "org.jwclab.iotrecord"
	stateList.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		return Deserialize(bytes, state.(*IoTRecord))
	}

	list := new(list)
	list.stateList = stateList

	return list
}

// IndexOfElement returns index of the record with the DateTime as time
// return -1 if not found
func IndexOfElement(array []IoTRecord, time string) int {
	return BinarySearch(array, 0, len(array)-1, time)
}

// BinarySearch return index of the record with the DateTime as time
// return -1 if not found
func BinarySearch(array []IoTRecord, left int, right int, time string) int {
	if right >= left {
		mid := left + (right-left)/2
		if array[mid].DateTime == time {
			return mid
		}

		if array[mid].DateTime > time {
			return BinarySearch(array, left, mid-1, time)
		}

		return BinarySearch(array, mid+1, right, time)
	}
	return -1
}
