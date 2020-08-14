/*
 * SPDX-License-Identifier: Apache-2.0
 */

package ledgerapi

import (
	"strings"
)

// SplitKey splits a key on colon
func SplitKey(key string) []string {
	return strings.Split(key, ":")
}

// MakeKey joins key parts using colon
func MakeKey(keyParts ...string) string {
	return strings.Join(keyParts, ":")
}

// StateInterface interface states must implement
// for use in a list
type StateInterface interface {
	// GetKey return components that combine to form the key
	GetKey() string
	Serialize() ([]byte, error)
}
