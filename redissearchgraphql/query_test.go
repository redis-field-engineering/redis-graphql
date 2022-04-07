package redissearchgraphql

import (
	"reflect"
	"testing"
)

var blankArgs map[string]interface{}
var blankArgsMap map[string]interface{}

// TestBlankQuery : check to make sure we receive a blank query
func TestBlankQuery(t *testing.T) {
	qstring, err := QueryBuilder(blankArgs, blankArgsMap, false)
	if err != nil {
		t.Log("Empty query returns err: ", err)
	}
	if !reflect.DeepEqual(qstring, "") {
		t.Error("Query should be empty but returns: ", qstring)
	}
}

// TestBlankQueryWildcard : check to make sure we receive a blank query
func TestBlankQueryWildcard(t *testing.T) {
	qstring, err := QueryBuilder(blankArgs, blankArgsMap, true)
	if err != nil {
		t.Log("Empty query wildcard returns err: ", err)
	}
	if !reflect.DeepEqual(qstring, "*") {
		t.Error("Query should return '*' empty but returns: ", qstring)
	}
}

// TestAndQuery : check to see if we can build a query with an AND
func TestAndQuery(t *testing.T) {
	fullArgs := map[string]interface{}{
		"_agg_groupby": "field1",
		"field1":       "value1",
		"field2":       "value2",
	}
	qstring, err := QueryBuilder(fullArgs, blankArgsMap, true)
	if err != nil {
		t.Log("And query returns err: ", err)
	}
	if !(reflect.DeepEqual(qstring, "@field1:value1 @field2:value2") ||
		reflect.DeepEqual(qstring, "@field2:value2 @field1:value1")) {
		t.Error("Query should be '@field1:value1 @field2:value2' but returns: ", qstring)
	}
}

// TestOrQuery : check to see if we can build a query with an OR
func TestOrQuery(t *testing.T) {
	fullArgs := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}
	argsMap := map[string]interface{}{
		"ormatch": true,
	}
	qstring, err := QueryBuilder(fullArgs, argsMap, true)
	if err != nil {
		t.Log("Or query returns err: ", err)
	}
	if !(reflect.DeepEqual(qstring, "(@field1:value1)|(@field2:value2)") || reflect.DeepEqual(qstring, "(@field2:value2)|(@field1:value1)")) {
		t.Error("Or Query should return '(A)|(B)' empty but returns: ", qstring)
	}
}

// TestGeoQuery : check to see if we can build a query with a geo
func TestGeoQuery(t *testing.T) {
	fullArgs := map[string]interface{}{
		"geo": map[string]interface{}{"lat": 1.0, "lon": 2.0, "radius": 3.0, "unit": "km"},
	}
	qstring, err := QueryBuilder(fullArgs, blankArgsMap, true)
	if err != nil {
		t.Log("Geo query returns err: ", err)
	}
	if !reflect.DeepEqual(qstring, "@geo: [2.000000,1.000000,3.000000,km]") {
		t.Error("And Query should return '@tags:{tag1|tag2|tag3}' empty but returns: ", qstring)
	}
}

// TestTagQuery : check to see if we can build a query with  tags
func TestTagQuery(t *testing.T) {
	fullArgs := map[string]interface{}{
		"tags": []interface{}{"tag1", "tag2", "tag3"},
	}
	qstring, err := QueryBuilder(fullArgs, blankArgsMap, false)
	if err != nil {
		t.Log("Geo query returns err: ", err)
	}
	if !reflect.DeepEqual(qstring, "@tags:{tag1|tag2|tag3}") {
		t.Error("And Query should return '@tags:{tag1|tag2|tag3}' empty but returns: ", qstring)
	}
}

// TestBTEQuery : check to see if we can build a query with between
func TestBTEQuery(t *testing.T) {
	fullArgs := map[string]interface{}{
		"num_bte": []interface{}{1.0, 3.0},
	}
	qstring, err := QueryBuilder(fullArgs, blankArgsMap, false)
	if err != nil {
		t.Log("BTE query returns err: ", err)
	}
	if !reflect.DeepEqual(qstring, "@num:[1.000000, 3.000000]") {
		t.Error("And Query should return '@num:[1.000000, 3.000000]' empty but returns: ", qstring)
	}
}

// TestNumQuery : check to see if we can build a query with between
func TestNumQuery(t *testing.T) {
	fullArgs := map[string]interface{}{
		"num": 3.0,
	}
	qstring, err := QueryBuilder(fullArgs, blankArgsMap, false)
	if err != nil {
		t.Log("BTE query returns err: ", err)
	}
	if !reflect.DeepEqual(qstring, "@num:[3.000000,3.000000]") {
		t.Error("Numeric Query should return '@num:[3.000000,3.000000]' empty but returns: ", qstring)
	}
}

// TestGTEQuery : check to see if we can build a query with between
func TestGTEQuery(t *testing.T) {
	fullArgs := map[string]interface{}{
		"num_gte": 3.0,
	}
	qstring, err := QueryBuilder(fullArgs, blankArgsMap, false)
	if err != nil {
		t.Log("BTE query returns err: ", err)
	}
	if !reflect.DeepEqual(qstring, "@num:[3.000000,+inf]") {
		t.Error("BET Query should return '@num:[3.000000,+inf]' empty but returns: ", qstring)
	}
}

// TestLTEQuery : check to see if we can build a query with between
func TestLTEQuery(t *testing.T) {
	fullArgs := map[string]interface{}{
		"num_lte": 3.0,
	}
	qstring, err := QueryBuilder(fullArgs, blankArgsMap, false)
	if err != nil {
		t.Log("BTE query returns err: ", err)
	}
	if !reflect.DeepEqual(qstring, "@num:[-inf,3.000000]") {
		t.Error("BET Query should return '@num:[-inf,3.000000]' empty but returns: ", qstring)
	}
}
