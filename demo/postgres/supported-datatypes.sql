-- Supports all datatypes of postgresql (https://www.postgresql.org/docs/9.6/static/datatype.html)
-- We don't support any user custom datatypes.

DROP TABLE supported_datatypes;
CREATE TABLE supported_datatypes (

	-- Integer type
	col_bigint		 		int8,
	col_bigserial	 		bigserial,
	col_integer		 		int,
	col_interval 			interval,
	col_numeric 			numeric(4,2),
	col_real			 	float4,
	col_smallint	 		smallint,
	col_smallserial	 		smallserial,
	col_serial		 		serial,
	col_double_precision	 float8,

	-- Bit String Type
	col_bit 				bit,
	col_varbit 				varbit(4),

	-- Boolean type
	col_boolean		 		bool,

	-- Character type
	col_character 			char(10),
	col_character_varying 	varchar(10),
	col_text	 	 		text,

	-- Network Address Type
	col_inet	 			inet,
	col_macaddr	 			macaddr,
	col_cidr	 			cidr,

	-- Date / Time type
	col_date	 			date,
	col_time 		 		time without time zone,
	col_time_tz 	 		time without time zone,
	col_timestamp 	 		timestamp with time zone,
	col_timestamp_tz 		timestamp without time zone,

	-- Monetary Types
	col_money	 			money,

	-- JSON Type
	col_json	 			json,
  col_jsonb	 			jsonb,

  -- XML Type
   col_xml	 				xml,

  -- Text Search Type
  col_tsquery	 	 		tsquery,
  col_tsvector	 		tsvector,

  -- Geometric Type
  col_box			 		box,
  col_circle	 			circle,
  col_line	 			line,
  col_lseg	 			lseg,
  col_path	 			path,
  col_polygon	 			polygon,
  col_point	 			point,

	-- Bytea / blob type
	col_bytea	 			bytea,

	-- Log Sequence Number
	col_pg_lsn	 			pg_lsn,

	-- txid snapshot
	col_txid_snapshot	 	txid_snapshot,

	-- UUID Type
	col_uuid	 			uuid

);
