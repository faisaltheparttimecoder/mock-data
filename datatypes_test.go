package main

import (
	"reflect"
	"testing"
)

// Test: SupportedDataTypes, this function is used when user ask to create random tables with
// random datatypes
func TestSupportedDataTypes(t *testing.T) {
	tests := struct {
		name string
		want []string
	}{
		"all_supported_datatypes", []string{
			"int8,",
			"bigint,",
			"bigserial,",
			"integer,",
			"smallint,",
			"smallserial,",
			"serial,",
			"oid,",
			"real,",
			"double precision,",
			"numeric(4,2),",
			"bit,",
			"varbit(4),",
			"bool,",
			"char(10),",
			"character varying(10),",
			"text,",
			"inet,",
			"macaddr,",
			"cidr,",
			"interval,",
			"date,",
			"time without time zone,",
			"time with time zone,",
			"timestamp with time zone,",
			"timestamp without time zone,",
			"timestamp(6) with time zone,",
			"timestamp(6) without time zone,",
			"money,",
			"json,",
			"jsonb,",
			"xml,",
			"tsquery,",
			"tsvector,",
			"box,",
			"circle,",
			"line,",
			"lseg,",
			"path,",
			"polygon,",
			"point,",
			"bytea,",
			"pg_lsn,",
			"txid_snapshot,",
			"uuid,",
			"smallint[],",
			"int[],",
			"bigint[],",
			"char(10)[],",
			"character varying(10)[],",
			"bit(10)[],",
			"varbit(4)[],",
			"numeric[],",
			"numeric(5,3)[],",
			"real[],",
			"double precision[],",
			"money[],",
			"time without time zone[],",
			"interval[],",
			"date[],",
			"time with time zone[],",
			"timestamp with time zone[],",
			"timestamp without time zone[],",
			"text[],",
			"bool[],",
			"inet[],",
			"macaddr[],",
			"cidr[],",
			"uuid[],",
			"txid_snapshot[],",
			"pg_lsn[],",
			"tsquery[],",
			"tsvector[],",
			"box[],",
			"circle[],",
			"line[],",
			"lseg[],",
			"path[],",
			"polygon[],",
			"point[],",
			"json[],",
			"jsonb[],",
			"xml[],",
		},
	}
	t.Run(tests.name, func(t *testing.T) {
		s := SupportedDataTypes()
		if !reflect.DeepEqual(s, tests.want) {
			t.Errorf("TestSupportedDataTypes() = %v, want %v", s, tests.want)
		}
	})
}

// Test: supportedDataTypesDemoTable, create the table of all the support datatype
func TestSupportedDataTypesDemoTable(t *testing.T) {
	t.Run("create_table_with_all_supported_datatypes", func(t *testing.T) {
		setDatabaseConfigForTest()
		if err := supportedDataTypesDemoTable(); err != nil {
			t.Errorf("TestSupportedDataTypesDemoTable, failed to create table, err: %v", err)
		}
		if _, err := ExecuteDB("SELECT 'public.supported_datatypes'::regclass;"); err != nil {
			t.Errorf("TestSupportedDataTypesDemoTable, had to create a table but its not found in the database")
		}
	})
}
