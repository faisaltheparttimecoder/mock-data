DROP TABLE IF EXISTS supported_datatypes;
CREATE TABLE supported_datatypes (
sample_bigint bigint,
sample_bigserial bigserial,
sample_bit bit(5),
sample_varbit bit varying(7),
sample_boolean bool,
-- sample_box	 	32 bytes	((x1,y1),(x2,y2))	rectangular box in the plane - not allowed in distribution key columns.
-- sample_bytea1	 	1 byte + binary string	sequence of octets	variable-length binary string
sample_character char(20),
sample_varcharacter varchar(50),
sample_cidr	cidr,
-- sample_circle	 	24 bytes	<(x,y),r> (center and radius)	circle in the plane - not allowed in distribution key columns.
sample_date	date,
sample_decimal numeric(5,3),
sample_double float,
sample_inet	inet,
sample_integer int,
-- sample_interval [ (p) ]	 	12 bytes	-178000000 years - 178000000 years	time span
-- sample_lseg	 	32 bytes	((x1,y1),(x2,y2))	line segment in the plane - not allowed in distribution key columns.
sample_macaddr macaddr,
-- sample_money	 	4 bytes	-21474836.48 to +21474836.47	currency amount
-- sample_path1	 	16+16n bytes	[(x1,y1),...]	geometric path in the plane - not allowed in distribution key columns.
-- sample_point	 	16 bytes	(x,y)	geometric point in the plane - not allowed in distribution key columns.
-- sample_polygon	 	40+16n bytes	((x1,y1),...)	closed geometric path in the plane - not allowed in distribution key columns.
sample_real	float4,
sample_serial serial,
sample_smallint	smallint,
sample_text	text,
sample_time time without time zone,
sample_timetz time with time zone,
sample_timestamp timestamp without time zone,
sample_timestamptz timestamp with time zone
-- sample_xml1	 	1 byte + xml size	xml of any length	variable unlimited lengthsample_
);
