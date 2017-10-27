-- Supports all datatypes of postgresql (https://www.postgresql.org/docs/9.6/static/datatype.html)
-- We don't support any user custom datatypes.

DROP TABLE supported_datatypes;
CREATE TABLE supported_datatypes (

  -- Integer type
  col_bigint                                  int8,
  col_bigserial                               bigserial,
  col_integer                                 int,
  col_smallint                                smallint,
  col_smallserial                             smallserial,
  col_serial                                  serial,

  -- Float type
  col_real                                    float4,
  col_double_precision                        float8,
  col_numeric                                 numeric(4,2),

  -- Bit String Type
  col_bit                                     bit,
  col_varbit                                  varbit(4),

  -- Boolean type
  col_boolean                                 bool,

  -- Character type
  col_character                               char(10),
  col_character_varying                       varchar(10),
  col_text                                    text,

  -- Network Address Type
  col_inet                                    inet,
  col_macaddr                                 macaddr,
  col_cidr                                    cidr,

  -- Date / Time type
  col_interval                                interval,
  col_date                                    date,
  col_time                                    time without time zone,
  col_time_tz                                 time with time zone,
  col_timestamp                               timestamp with time zone,
  col_timestamp_tz                            timestamp without time zone,

  -- Monetary Types
  col_money                                   money,

  -- JSON Type
  col_json                                    json,
  col_jsonb                                   jsonb,

  -- XML Type
  col_xml                                     xml,

  -- Text Search Type
  col_tsquery                                 tsquery,
  col_tsvector                                tsvector,

  -- Geometric Type
  col_box                                     box,
  col_circle                                  circle,
  col_line                                    line,
  col_lseg                                    lseg,
  col_path                                    path,
  col_polygon                                 polygon,
  col_point                                   point,

  -- Bytea / blob type
  col_bytea                                   bytea,

  -- Log Sequence Number
  col_pg_lsn                                  pg_lsn,

  -- txid snapshot
  col_txid_snapshot                           txid_snapshot,

  -- UUID Type
  col_uuid                                    uuid,

  -- Array Datatypes
  col_smallint_array                         smallint[],
  col_int_array                              int[],
  col_bigint_array                           bigint[],
  col_character_array                        char(10)[],
  col_char_varying_array                     varchar(10)[],
  col_bit_array                              bit(10)[],
  col_varbit_array                           varbit(4)[],
  col_numeric_array                          numeric[],
  col_numeric_range_array                    numeric(5,3)[],
  col_double_precsion_array                  float4[],
  col_real_array                             float8[],
  col_money_array                            money[],
  col_time_array                             time without time zone[],
  col_intreval_array                         interval[],
  col_date_array                             date[],
  col_time_tz_array                          time with time zone[],
  col_timestamp_array                        timestamp with time zone[],
  col_timestamp_tz_array                     timestamp without time zone[],
  col_text_array                             text[],
  col_bool_array                             bool[],
  col_inet_array                             inet[],
  col_macaddr_array                          macaddr[],
  col_cidr_array                             cidr[],
  col_uuid_array                             uuid[],
  col_txid_snapshot_array                    txid_snapshot[],
  col_pg_lsn_array                           pg_lsn[],
  col_tsquery_array                          tsquery[],
  col_tsvector_array                         tsvector[]

);
