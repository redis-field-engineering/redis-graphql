package redissearchgraphql

import (
	"testing"
)

// Test Index with - in name
func TestDashIndexName(t *testing.T) {
	var badIndices = []string{"Idx1", "Idx2", "Idx-3"}
	err := checkIndexNames(badIndices)
	if err.Error() != "Index name Idx-3 is not GraphQL compliant with a : or - in name" {
		t.Error("Bad index name Idx-3 should have failed")
	}

}

// Test Index with : in name
func TestColonIndexName(t *testing.T) {
	var badIndices = []string{"Idx1", "Idx2", "Idx:3"}
	err := checkIndexNames(badIndices)
	if err.Error() != "Index name Idx:3 is not GraphQL compliant with a : or - in name" {
		t.Error("Bad index name Idx:3 should have failed")
	}

}

// Test Index with all good names
func TestGoodIndexName(t *testing.T) {
	var badIndices = []string{"Index123", "Idx_2", "AnotherGoodIndex"}
	err := checkIndexNames(badIndices)
	if err != nil {
		t.Error("All index names are correct, this test should have passed")
	}

}
