/*
 * SPDX-License-Identifier: Apache-2.0
 */

package iotrecord

import (
	"encoding/json"
	"fmt"
	"sort"

	ledgerapi "jwclab/iotrecord/ledger-api"
)

// CreateIoTRecordKey creates a key for commercial papers
func CreateIoTRecordKey(monitorID string, name string) string {
	return ledgerapi.MakeKey(monitorID, name)
}

// Used for managing the fact status is private but want it in world state
type iotRecordAlias IoTRecord
type jsonIoTRecord struct {
	*iotRecordAlias
	Class string `json:"class"`
	Key   string `json:"key"`
}

// IoTRecord defines an IOT record
type IoTRecord struct {
	Name      string  `json:"name"`
	Unit      string  `json:"unit"`
	DateTime  string  `json:"dateTime"`
	Value     float32 `json:"value"`
	MonitorID string  `json:"monitorId"`
	class     string  `metadata:"class"`
	key       string  `metadata:"key"`
}

// By ...
type By func(p1, p2 *IoTRecord) bool

// Sort ...
func (by By) Sort(records []IoTRecord) {
	ps := &recordSorter{
		records: records,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// planetSorter joins a By function and a slice of Planets to be sorted.
type recordSorter struct {
	records []IoTRecord
	by      func(p1, p2 *IoTRecord) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *recordSorter) Len() int {
	return len(s.records)
}

// Swap is part of sort.Interface.
func (s *recordSorter) Swap(i, j int) {
	s.records[i], s.records[j] = s.records[j], s.records[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *recordSorter) Less(i, j int) bool {
	return s.by(&s.records[i], &s.records[j])
}

// UnmarshalJSON special handler for managing JSON marshalling
func (cp *IoTRecord) UnmarshalJSON(data []byte) error {
	jcp := jsonIoTRecord{iotRecordAlias: (*iotRecordAlias)(cp)}

	err := json.Unmarshal(data, &jcp)

	if err != nil {
		return err
	}

	return nil
}

// MarshalJSON special handler for managing JSON marshalling
func (cp IoTRecord) MarshalJSON() ([]byte, error) {
	jcp := jsonIoTRecord{iotRecordAlias: (*iotRecordAlias)(&cp), Class: "org.jwclab.iotrecord", Key: ledgerapi.MakeKey(cp.MonitorID, cp.Name)}

	return json.Marshal(&jcp)
}

// GetKey returns values which should be used to form key
func (cp *IoTRecord) GetKey() string {
	return cp.MonitorID + ":" + cp.Name
}

// Serialize formats the commercial paper as JSON bytes
func (cp *IoTRecord) Serialize() ([]byte, error) {
	return json.Marshal(cp)
}

// Deserialize formats the commercial paper from JSON bytes
func Deserialize(bytes []byte, cp *IoTRecord) error {
	err := json.Unmarshal(bytes, cp)

	if err != nil {
		return fmt.Errorf("Error deserializing commercial paper. %s", err.Error())
	}

	return nil
}
