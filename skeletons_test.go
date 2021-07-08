package main

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

// Test: JsonSkeleton, check if it sends in a valid JSON template
func TestJsonSkeleton(t *testing.T) {
	var js json.RawMessage
	t.Run("test_validity_of_json_skeleton", func(t *testing.T) {
		if got := JsonSkeleton(); json.Unmarshal([]byte(got), &js) != nil {
			t.Errorf("TestJsonSkeleton = %v, want valid JSON format", got)
		}
	})
}

// Test: XMLSkeleton, check if it sends in a valid xml template
func TestXMLSkeleton(t *testing.T) {
	t.Run("test_validity_of_xml_skeleton", func(t *testing.T) {
		if got := XMLSkeleton(); xml.Unmarshal([]byte(got), new(interface{})) != nil {
			t.Errorf("TestXMLSkeleton = %v, want valid xml format", got)
		}
	})
}

// Test: demoDatabasePostgres, it should create a database without any error
func TestDemoDatabasePostgres(t *testing.T) {
	setDatabaseConfigForTest()
	t.Run("test_creation_of_demo_postgres_database", func(t *testing.T) {
		if _, err := ExecuteDB(demoDatabasePostgres()); err != nil {
			t.Errorf("TestDemoDatabasePostgres should create a demo database, but got err: %v", err)
		}
	})
}

// Test: demoDatabaseGreenplum, it should create a database without any error
func TestDemoDatabaseGreenplum(t *testing.T) {
	// Skipping this step
}