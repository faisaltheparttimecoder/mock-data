package main

// Json Skeleton
func JsonSkeleton() string {
	return `{"_id": "%s","index": "%s","guid": "%s","isActive": "%s","balance": "%s.%s","website": "https://%s/%s","age": "%s","username": "%s","eyeColor": "%s","name": "%s","gender": "%s","company": "%s","email": "%s","phone": "%s","address": "%s","zipcode": "%s","state": "%s","country": "%s","about": "%s","Machine IP": "%s","job title": "%s","registered": "%s-%s-%sT%s:%s:%s-%s:%s","latitude": "%s.%s","longitude": "%s.%s","tags": ["%s","%s","%s","%s","%s","%s","%s"],"friends": [{ "id": "%s", "name": "%s"},{ "id": "%s", "name": "%s"},{ "id": "%s", "name": "%s"}],"greeting": "%s","favoriteBrand": "%s"}`
}

// XML Skeleton
func XMLSkeleton() string {
	return `<?xml version="1.0" encoding="UTF-8"?><shiporder orderid="%s" xmlns:xsi="http://%s/%s/%s" xsi:noNamespaceSchemaLocation="shiporder.xsd"> <orderperson>%s</orderperson> <shipto> <name>%s</name> <address>%s</address> <city>%s</city> <country>%s</country> <email>%s</email> <phone>%s</phone> </shipto> <item> <title>%s</title> <note>%s</note> <quantity>%s</quantity> <color>%s</color> <price>%s.%s</price> </item> <item> <title>%s</title> <quantity>%s</quantity> <price>%s.%s</price> </item></shiporder>`
}

// Drop demo database sql
func dropDemoDatabase() string {
	return `
--
-- DROP ALL THE OBJECTS 
--
DO $$ DECLARE
    r RECORD;
BEGIN
    -- if the schema you operate on is not "current", you will want to
    -- replace current_schema() in query with 'schematodeletetablesfrom'
    -- *and* update the generate 'DROP...' accordingly.
    FOR r IN (SELECT schemaname, tablename FROM pg_tables WHERE schemaname NOT IN ('pg_catalog', 'information_schema')) LOOP
        EXECUTE 'DROP TABLE IF EXISTS "' || quote_ident(r.schemaname) ||'"."'|| quote_ident(r.tablename) || '" CASCADE';
    END LOOP;
END $$;

--
-- DROP PUBLIC SCHEMA
--
DROP SCHEMA if exists public cascade;
CREATE SCHEMA public;
DROP ROLE IF EXISTS mockdata;
CREATE ROLE mockdata;
`
}

// Demo database sql
func demoDatabasePostgres() string {
	Debug("Creating a demo database for postgres")
	return `
--
-- NOTE:
--
-- File paths need to be edited. Search for $$PATH$$ and
-- replace it with the path to the directory containing
-- the extracted data files.
--
--
-- PostgreSQL database dump
--

-- Dumped from database version 11.3
-- Dumped by pg_dump version 11.2
SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

--
-- DROP ALL THE OBJECTS 
--
ALTER TABLE IF EXISTS ONLY public.store DROP CONSTRAINT store_manager_staff_id_fkey;
ALTER TABLE IF EXISTS ONLY public.store DROP CONSTRAINT store_address_id_fkey;
ALTER TABLE IF EXISTS ONLY public.staff DROP CONSTRAINT staff_address_id_fkey;
ALTER TABLE IF EXISTS ONLY public.rental DROP CONSTRAINT rental_staff_id_key;
ALTER TABLE IF EXISTS ONLY public.rental DROP CONSTRAINT rental_inventory_id_fkey;
ALTER TABLE IF EXISTS ONLY public.rental DROP CONSTRAINT rental_customer_id_fkey;
ALTER TABLE IF EXISTS ONLY public.payment DROP CONSTRAINT payment_staff_id_fkey;
ALTER TABLE IF EXISTS ONLY public.payment DROP CONSTRAINT payment_rental_id_fkey;
ALTER TABLE IF EXISTS ONLY public.payment DROP CONSTRAINT payment_customer_id_fkey;
ALTER TABLE IF EXISTS ONLY public.inventory DROP CONSTRAINT inventory_film_id_fkey;
ALTER TABLE IF EXISTS ONLY public.city DROP CONSTRAINT fk_city;
ALTER TABLE IF EXISTS ONLY public.address DROP CONSTRAINT fk_address_city;
ALTER TABLE IF EXISTS ONLY public.film DROP CONSTRAINT film_language_id_fkey;
ALTER TABLE IF EXISTS ONLY public.film_category DROP CONSTRAINT film_category_film_id_fkey;
ALTER TABLE IF EXISTS ONLY public.film_category DROP CONSTRAINT film_category_category_id_fkey;
ALTER TABLE IF EXISTS ONLY public.film_actor DROP CONSTRAINT film_actor_film_id_fkey;
ALTER TABLE IF EXISTS ONLY public.film_actor DROP CONSTRAINT film_actor_actor_id_fkey;
ALTER TABLE IF EXISTS ONLY public.customer DROP CONSTRAINT customer_address_id_fkey;
DROP TRIGGER IF EXISTS last_updated ON public.store;
DROP TRIGGER IF EXISTS last_updated ON public.staff;
DROP TRIGGER IF EXISTS last_updated ON public.rental;
DROP TRIGGER IF EXISTS last_updated ON public.language;
DROP TRIGGER IF EXISTS last_updated ON public.inventory;
DROP TRIGGER IF EXISTS last_updated ON public.film_category;
DROP TRIGGER IF EXISTS last_updated ON public.film_actor;
DROP TRIGGER IF EXISTS last_updated ON public.film;
DROP TRIGGER IF EXISTS last_updated ON public.customer;
DROP TRIGGER IF EXISTS last_updated ON public.country;
DROP TRIGGER IF EXISTS last_updated ON public.city;
DROP TRIGGER IF EXISTS last_updated ON public.category;
DROP TRIGGER IF EXISTS last_updated ON public.address;
DROP TRIGGER IF EXISTS last_updated ON public.actor;
DROP TRIGGER IF EXISTS film_fulltext_trigger ON public.film;
DROP INDEX IF EXISTS public.idx_unq_rental_rental_date_inventory_id_customer_id;
DROP INDEX IF EXISTS public.idx_unq_manager_staff_id;
DROP INDEX IF EXISTS public.idx_title;
DROP INDEX IF EXISTS public.idx_store_id_film_id;
DROP INDEX IF EXISTS public.idx_last_name;
DROP INDEX IF EXISTS public.idx_fk_store_id;
DROP INDEX IF EXISTS public.idx_fk_staff_id;
DROP INDEX IF EXISTS public.idx_fk_rental_id;
DROP INDEX IF EXISTS public.idx_fk_language_id;
DROP INDEX IF EXISTS public.idx_fk_inventory_id;
DROP INDEX IF EXISTS public.idx_fk_film_id;
DROP INDEX IF EXISTS public.idx_fk_customer_id;
DROP INDEX IF EXISTS public.idx_fk_country_id;
DROP INDEX IF EXISTS public.idx_fk_city_id;
DROP INDEX IF EXISTS public.idx_fk_address_id;
DROP INDEX IF EXISTS public.idx_actor_last_name;
DROP INDEX IF EXISTS public.film_fulltext_idx;
ALTER TABLE IF EXISTS ONLY public.store DROP CONSTRAINT store_pkey;
ALTER TABLE IF EXISTS ONLY public.staff DROP CONSTRAINT staff_pkey;
ALTER TABLE IF EXISTS ONLY public.rental DROP CONSTRAINT rental_pkey;
ALTER TABLE IF EXISTS ONLY public.payment DROP CONSTRAINT payment_pkey;
ALTER TABLE IF EXISTS ONLY public.language DROP CONSTRAINT language_pkey;
ALTER TABLE IF EXISTS ONLY public.inventory DROP CONSTRAINT inventory_pkey;
ALTER TABLE IF EXISTS ONLY public.film DROP CONSTRAINT film_pkey;
ALTER TABLE IF EXISTS ONLY public.film_category DROP CONSTRAINT film_category_pkey;
ALTER TABLE IF EXISTS ONLY public.film_actor DROP CONSTRAINT film_actor_pkey;
ALTER TABLE IF EXISTS ONLY public.customer DROP CONSTRAINT customer_pkey;
ALTER TABLE IF EXISTS ONLY public.country DROP CONSTRAINT country_pkey;
ALTER TABLE IF EXISTS ONLY public.city DROP CONSTRAINT city_pkey;
ALTER TABLE IF EXISTS ONLY public.category DROP CONSTRAINT category_pkey;
ALTER TABLE IF EXISTS ONLY public.address DROP CONSTRAINT address_pkey;
ALTER TABLE IF EXISTS ONLY public.actor DROP CONSTRAINT actor_pkey;
DROP VIEW IF EXISTS public.staff_list;
DROP VIEW IF EXISTS public.sales_by_store;
DROP TABLE IF EXISTS public.store;
DROP SEQUENCE IF EXISTS public.store_store_id_seq;
DROP TABLE IF EXISTS public.staff;
DROP SEQUENCE IF EXISTS public.staff_staff_id_seq;
DROP VIEW IF EXISTS public.sales_by_film_category;
DROP TABLE IF EXISTS public.rental;
DROP SEQUENCE IF EXISTS public.rental_rental_id_seq;
DROP TABLE IF EXISTS public.payment;
DROP SEQUENCE IF EXISTS public.payment_payment_id_seq;
DROP VIEW IF EXISTS public.nicer_but_slower_film_list;
DROP TABLE IF EXISTS public.language;
DROP SEQUENCE IF EXISTS public.language_language_id_seq;
DROP TABLE IF EXISTS public.inventory;
DROP SEQUENCE IF EXISTS public.inventory_inventory_id_seq;
DROP VIEW IF EXISTS public.film_list;
DROP VIEW IF EXISTS public.customer_list;
DROP TABLE IF EXISTS public.country;
DROP SEQUENCE IF EXISTS public.country_country_id_seq;
DROP TABLE IF EXISTS public.city;
DROP SEQUENCE IF EXISTS public.city_city_id_seq;
DROP TABLE IF EXISTS public.address;
DROP SEQUENCE IF EXISTS public.address_address_id_seq;
DROP VIEW IF EXISTS public.actor_info;
DROP TABLE IF EXISTS public.film_category;
DROP TABLE IF EXISTS public.film_actor;
DROP TABLE IF EXISTS public.film;
DROP SEQUENCE IF EXISTS public.film_film_id_seq;
DROP TABLE IF EXISTS public.category;
DROP SEQUENCE IF EXISTS public.category_category_id_seq;
DROP TABLE IF EXISTS public.actor;
DROP SEQUENCE IF EXISTS public.actor_actor_id_seq;
DROP AGGREGATE IF EXISTS public.group_concat(text);
DROP FUNCTION IF EXISTS public.rewards_report(min_monthly_purchases integer, min_dollar_amount_purchased numeric);
DROP TABLE IF EXISTS public.customer;
DROP SEQUENCE IF EXISTS public.customer_customer_id_seq;
DROP FUNCTION IF EXISTS public.last_updated();
DROP FUNCTION IF EXISTS public.last_day(timestamp without time zone);
DROP FUNCTION IF EXISTS public.inventory_in_stock(p_inventory_id integer);
DROP FUNCTION IF EXISTS public.inventory_held_by_customer(p_inventory_id integer);
DROP FUNCTION IF EXISTS public.get_customer_balance(p_customer_id integer, p_effective_date timestamp without time zone);
DROP FUNCTION IF EXISTS public.film_not_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer);
DROP FUNCTION IF EXISTS public.film_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer);
DROP FUNCTION IF EXISTS public._group_concat(text, text);
DROP DOMAIN IF EXISTS public.year;
DROP TYPE IF EXISTS public.mpaa_rating;
DROP EXTENSION IF EXISTS plpgsql;
DROP OWNED BY mockdata CASCADE;
DROP ROLE IF EXISTS mockdata;
CREATE ROLE mockdata;
DROP TYPE IF EXISTS rating;
CREATE TYPE rating as ENUM ('good', 'ok', 'bad');

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: mockdata
--
CREATE SCHEMA IF NOT EXISTS public;
ALTER SCHEMA public OWNER TO mockdata;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: mockdata
--
COMMENT ON SCHEMA public IS 'Standard public schema';
--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--
COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';

SET search_path = public, pg_catalog;

--
-- Name: mpaa_rating; Type: TYPE; Schema: public; Owner: mockdata
--

CREATE TYPE public.mpaa_rating AS ENUM (
    'G',
    'PG',
    'PG-13',
    'R',
    'NC-17'
);


ALTER TYPE public.mpaa_rating OWNER TO mockdata;

--
-- Name: year; Type: DOMAIN; Schema: public; Owner: mockdata
--

CREATE DOMAIN public.year AS integer
	CONSTRAINT year_check CHECK (((VALUE >= 1901) AND (VALUE <= 2155)));


ALTER DOMAIN public.year OWNER TO mockdata;

--
-- Name: _group_concat(text, text); Type: FUNCTION; Schema: public; Owner: mockdata
--

CREATE FUNCTION public._group_concat(text, text) RETURNS text
    LANGUAGE sql IMMUTABLE
    AS $_$
SELECT CASE
  WHEN $2 IS NULL THEN $1
  WHEN $1 IS NULL THEN $2
  ELSE $1 || ', ' || $2
END
$_$;


ALTER FUNCTION public._group_concat(text, text) OWNER TO mockdata;

--
-- Name: film_in_stock(integer, integer); Type: FUNCTION; Schema: public; Owner: mockdata
--

CREATE FUNCTION public.film_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer) RETURNS SETOF integer
    LANGUAGE sql
    AS $_$
     SELECT inventory_id
     FROM inventory
     WHERE film_id = $1
     AND store_id = $2
     AND inventory_in_stock(inventory_id);
$_$;


ALTER FUNCTION public.film_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer) OWNER TO mockdata;

--
-- Name: film_not_in_stock(integer, integer); Type: FUNCTION; Schema: public; Owner: mockdata
--

CREATE FUNCTION public.film_not_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer) RETURNS SETOF integer
    LANGUAGE sql
    AS $_$
    SELECT inventory_id
    FROM inventory
    WHERE film_id = $1
    AND store_id = $2
    AND NOT inventory_in_stock(inventory_id);
$_$;


ALTER FUNCTION public.film_not_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer) OWNER TO mockdata;

--
-- Name: get_customer_balance(integer, timestamp without time zone); Type: FUNCTION; Schema: public; Owner: mockdata
--

CREATE FUNCTION public.get_customer_balance(p_customer_id integer, p_effective_date timestamp without time zone) RETURNS numeric
    LANGUAGE plpgsql
    AS $$
       --#OK, WE NEED TO CALCULATE THE CURRENT BALANCE GIVEN A CUSTOMER_ID AND A DATE
       --#THAT WE WANT THE BALANCE TO BE EFFECTIVE FOR. THE BALANCE IS:
       --#   1) RENTAL FEES FOR ALL PREVIOUS RENTALS
       --#   2) ONE DOLLAR FOR EVERY DAY THE PREVIOUS RENTALS ARE OVERDUE
       --#   3) IF A FILM IS MORE THAN RENTAL_DURATION * 2 OVERDUE, CHARGE THE REPLACEMENT_COST
       --#   4) SUBTRACT ALL PAYMENTS MADE BEFORE THE DATE SPECIFIED
DECLARE
    v_rentfees DECIMAL(5,2); --#FEES PAID TO RENT THE VIDEOS INITIALLY
    v_overfees INTEGER;      --#LATE FEES FOR PRIOR RENTALS
    v_payments DECIMAL(5,2); --#SUM OF PAYMENTS MADE PREVIOUSLY
BEGIN
    SELECT COALESCE(SUM(film.rental_rate),0) INTO v_rentfees
    FROM film, inventory, rental
    WHERE film.film_id = inventory.film_id
      AND inventory.inventory_id = rental.inventory_id
      AND rental.rental_date <= p_effective_date
      AND rental.customer_id = p_customer_id;

    SELECT COALESCE(SUM(IF((rental.return_date - rental.rental_date) > (film.rental_duration * '1 day'::interval),
        ((rental.return_date - rental.rental_date) - (film.rental_duration * '1 day'::interval)),0)),0) INTO v_overfees
    FROM rental, inventory, film
    WHERE film.film_id = inventory.film_id
      AND inventory.inventory_id = rental.inventory_id
      AND rental.rental_date <= p_effective_date
      AND rental.customer_id = p_customer_id;

    SELECT COALESCE(SUM(payment.amount),0) INTO v_payments
    FROM payment
    WHERE payment.payment_date <= p_effective_date
    AND payment.customer_id = p_customer_id;

    RETURN v_rentfees + v_overfees - v_payments;
END
$$;


ALTER FUNCTION public.get_customer_balance(p_customer_id integer, p_effective_date timestamp without time zone) OWNER TO mockdata;

--
-- Name: inventory_held_by_customer(integer); Type: FUNCTION; Schema: public; Owner: mockdata
--

CREATE FUNCTION public.inventory_held_by_customer(p_inventory_id integer) RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_customer_id INTEGER;
BEGIN

  SELECT customer_id INTO v_customer_id
  FROM rental
  WHERE return_date IS NULL
  AND inventory_id = p_inventory_id;

  RETURN v_customer_id;
END $$;


ALTER FUNCTION public.inventory_held_by_customer(p_inventory_id integer) OWNER TO mockdata;

--
-- Name: inventory_in_stock(integer); Type: FUNCTION; Schema: public; Owner: mockdata
--

CREATE FUNCTION public.inventory_in_stock(p_inventory_id integer) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_rentals INTEGER;
    v_out     INTEGER;
BEGIN
    -- AN ITEM IS IN-STOCK IF THERE ARE EITHER NO ROWS IN THE rental TABLE
    -- FOR THE ITEM OR ALL ROWS HAVE return_date POPULATED

    SELECT count(*) INTO v_rentals
    FROM rental
    WHERE inventory_id = p_inventory_id;

    IF v_rentals = 0 THEN
      RETURN TRUE;
    END IF;

    SELECT COUNT(rental_id) INTO v_out
    FROM inventory LEFT JOIN rental USING(inventory_id)
    WHERE inventory.inventory_id = p_inventory_id
    AND rental.return_date IS NULL;

    IF v_out > 0 THEN
      RETURN FALSE;
    ELSE
      RETURN TRUE;
    END IF;
END $$;


ALTER FUNCTION public.inventory_in_stock(p_inventory_id integer) OWNER TO mockdata;

--
-- Name: last_day(timestamp without time zone); Type: FUNCTION; Schema: public; Owner: mockdata
--

CREATE FUNCTION public.last_day(timestamp without time zone) RETURNS date
    LANGUAGE sql IMMUTABLE STRICT
    AS $_$
  SELECT CASE
    WHEN EXTRACT(MONTH FROM $1) = 12 THEN
      (((EXTRACT(YEAR FROM $1) + 1) operator(pg_catalog.||) '-01-01')::date - INTERVAL '1 day')::date
    ELSE
      ((EXTRACT(YEAR FROM $1) operator(pg_catalog.||) '-' operator(pg_catalog.||) (EXTRACT(MONTH FROM $1) + 1) operator(pg_catalog.||) '-01')::date - INTERVAL '1 day')::date
    END
$_$;


ALTER FUNCTION public.last_day(timestamp without time zone) OWNER TO mockdata;

--
-- Name: last_updated(); Type: FUNCTION; Schema: public; Owner: mockdata
--

CREATE FUNCTION public.last_updated() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.last_update = CURRENT_TIMESTAMP;
    RETURN NEW;
END $$;


ALTER FUNCTION public.last_updated() OWNER TO mockdata;

--
-- Name: customer_customer_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.customer_customer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.customer_customer_id_seq OWNER TO mockdata;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: customer; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.customer (
    customer_id integer DEFAULT nextval('public.customer_customer_id_seq'::regclass) NOT NULL,
    store_id smallint NOT NULL,
    first_name character varying(45) NOT NULL,
    last_name character varying(45) NOT NULL,
    email character varying(50),
    address_id smallint NOT NULL,
    activebool boolean DEFAULT true NOT NULL,
    create_date date DEFAULT ('now'::text)::date NOT NULL,
    last_update timestamp without time zone DEFAULT now(),
    active integer
);


ALTER TABLE public.customer OWNER TO mockdata;

--
-- Name: rewards_report(integer, numeric); Type: FUNCTION; Schema: public; Owner: mockdata
--

CREATE FUNCTION public.rewards_report(min_monthly_purchases integer, min_dollar_amount_purchased numeric) RETURNS SETOF public.customer
    LANGUAGE plpgsql SECURITY DEFINER
    AS $_$
DECLARE
    last_month_start DATE;
    last_month_end DATE;
rr RECORD;
tmpSQL TEXT;
BEGIN

    /* Some sanity checks... */
    IF min_monthly_purchases = 0 THEN
        RAISE EXCEPTION 'Minimum monthly purchases parameter must be > 0';
    END IF;
    IF min_dollar_amount_purchased = 0.00 THEN
        RAISE EXCEPTION 'Minimum monthly dollar amount purchased parameter must be > $0.00';
    END IF;

    last_month_start := CURRENT_DATE - '3 month'::interval;
    last_month_start := to_date((extract(YEAR FROM last_month_start) || '-' || extract(MONTH FROM last_month_start) || '-01'),'YYYY-MM-DD');
    last_month_end := LAST_DAY(last_month_start);

    /*
    Create a temporary storage area for Customer IDs.
    */
    CREATE TEMPORARY TABLE tmpCustomer (customer_id INTEGER NOT NULL PRIMARY KEY);

    /*
    Find all customers meeting the monthly purchase requirements
    */

    tmpSQL := 'INSERT INTO tmpCustomer (customer_id)
        SELECT p.customer_id
        FROM payment AS p
        WHERE DATE(p.payment_date) BETWEEN '||quote_literal(last_month_start) ||' AND '|| quote_literal(last_month_end) || '
        GROUP BY customer_id
        HAVING SUM(p.amount) > '|| min_dollar_amount_purchased || '
        AND COUNT(customer_id) > ' ||min_monthly_purchases ;

    EXECUTE tmpSQL;

    /*
    Output ALL customer information of matching rewardees.
    Customize output as needed.
    */
    FOR rr IN EXECUTE 'SELECT c.* FROM tmpCustomer AS t INNER JOIN customer AS c ON t.customer_id = c.customer_id' LOOP
        RETURN NEXT rr;
    END LOOP;

    /* Clean up */
    tmpSQL := 'DROP TABLE tmpCustomer';
    EXECUTE tmpSQL;

RETURN;
END
$_$;


ALTER FUNCTION public.rewards_report(min_monthly_purchases integer, min_dollar_amount_purchased numeric) OWNER TO mockdata;

--
-- Name: group_concat(text); Type: AGGREGATE; Schema: public; Owner: mockdata
--

CREATE AGGREGATE public.group_concat(text) (
    SFUNC = public._group_concat,
    STYPE = text
);


ALTER AGGREGATE public.group_concat(text) OWNER TO mockdata;

--
-- Name: actor_actor_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.actor_actor_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.actor_actor_id_seq OWNER TO mockdata;

--
-- Name: actor; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.actor (
    actor_id integer DEFAULT nextval('public.actor_actor_id_seq'::regclass) NOT NULL,
    first_name character varying(45) NOT NULL,
    last_name character varying(45) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);

AlTER TABLE  public.actor ADD CONSTRAINT actor_id_check CHECK ( actor_id > 0);
ALTER TABLE public.actor OWNER TO mockdata;

--
-- Name: category_category_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.category_category_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.category_category_id_seq OWNER TO mockdata;

--
-- Name: category; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.category (
    category_id integer DEFAULT nextval('public.category_category_id_seq'::regclass) NOT NULL,
    name character varying(25) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE public.category ADD CONSTRAINT category_id UNIQUE(category_id);
ALTER TABLE public.category OWNER TO mockdata;

--
-- Name: film_film_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.film_film_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.film_film_id_seq OWNER TO mockdata;

--
-- Name: film; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.film (
    film_id integer DEFAULT nextval('public.film_film_id_seq'::regclass) NOT NULL,
    title character varying(255) NOT NULL,
    description text,
    release_year date,
    language_id smallint NOT NULL,
    rental_duration smallint DEFAULT 3 NOT NULL,
    rental_rate numeric(4,2) DEFAULT 4.99 NOT NULL,
    length smallint,
    replacement_cost numeric(5,2) DEFAULT 19.99 NOT NULL,
    rating char(1) DEFAULT 'G',
    user_rating rating,
    last_update timestamp without time zone DEFAULT now() NOT NULL,
    special_features text[],
    fulltext tsvector NOT NULL
);


ALTER TABLE public.film OWNER TO mockdata;

--
-- Name: film_actor; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.film_actor (
    actor_id smallint NOT NULL,
    film_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.film_actor OWNER TO mockdata;

--
-- Name: film_category; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.film_category (
    film_id smallint NOT NULL,
    category_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.film_category OWNER TO mockdata;

--
-- Name: actor_info; Type: VIEW; Schema: public; Owner: mockdata
--

CREATE VIEW public.actor_info AS
 SELECT a.actor_id,
    a.first_name,
    a.last_name,
    public.group_concat(DISTINCT (((c.name)::text || ': '::text) || ( SELECT public.group_concat((f.title)::text) AS group_concat
           FROM ((public.film f
             JOIN public.film_category fc_1 ON ((f.film_id = fc_1.film_id)))
             JOIN public.film_actor fa_1 ON ((f.film_id = fa_1.film_id)))
          WHERE ((fc_1.category_id = c.category_id) AND (fa_1.actor_id = a.actor_id))
          GROUP BY fa_1.actor_id))) AS film_info
   FROM (((public.actor a
     LEFT JOIN public.film_actor fa ON ((a.actor_id = fa.actor_id)))
     LEFT JOIN public.film_category fc ON ((fa.film_id = fc.film_id)))
     LEFT JOIN public.category c ON ((fc.category_id = c.category_id)))
  GROUP BY a.actor_id, a.first_name, a.last_name;


ALTER TABLE public.actor_info OWNER TO mockdata;

--
-- Name: address_address_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.address_address_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.address_address_id_seq OWNER TO mockdata;

--
-- Name: address; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.address (
    address_id integer DEFAULT nextval('public.address_address_id_seq'::regclass) NOT NULL,
    address character varying(50) NOT NULL,
    address2 character varying(50),
    district character varying(20) NOT NULL,
    city_id smallint NOT NULL,
    postal_code character varying(10),
    phone character varying(20) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.address OWNER TO mockdata;

--
-- Name: city_city_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.city_city_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.city_city_id_seq OWNER TO mockdata;

--
-- Name: city; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.city (
    city_id integer DEFAULT nextval('public.city_city_id_seq'::regclass) NOT NULL,
    city character varying(50) NOT NULL,
    country_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.city OWNER TO mockdata;

--
-- Name: country_country_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.country_country_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.country_country_id_seq OWNER TO mockdata;

--
-- Name: country; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.country (
    country_id integer DEFAULT nextval('public.country_country_id_seq'::regclass) NOT NULL,
    country character varying(50) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.country OWNER TO mockdata;

--
-- Name: customer_list; Type: VIEW; Schema: public; Owner: mockdata
--

CREATE VIEW public.customer_list AS
 SELECT cu.customer_id AS id,
    (((cu.first_name)::text || ' '::text) || (cu.last_name)::text) AS name,
    a.address,
    a.postal_code AS "zip code",
    a.phone,
    city.city,
    country.country,
        CASE
            WHEN cu.activebool THEN 'active'::text
            ELSE ''::text
        END AS notes,
    cu.store_id AS sid
   FROM (((public.customer cu
     JOIN public.address a ON ((cu.address_id = a.address_id)))
     JOIN public.city ON ((a.city_id = city.city_id)))
     JOIN public.country ON ((city.country_id = country.country_id)));


ALTER TABLE public.customer_list OWNER TO mockdata;

--
-- Name: film_list; Type: VIEW; Schema: public; Owner: mockdata
--

CREATE VIEW public.film_list AS
 SELECT film.film_id AS fid,
    film.title,
    film.description,
    category.name AS category,
    film.rental_rate AS price,
    film.length,
    film.rating,
    public.group_concat((((actor.first_name)::text || ' '::text) || (actor.last_name)::text)) AS actors
   FROM ((((public.category
     LEFT JOIN public.film_category ON ((category.category_id = film_category.category_id)))
     LEFT JOIN public.film ON ((film_category.film_id = film.film_id)))
     JOIN public.film_actor ON ((film.film_id = film_actor.film_id)))
     JOIN public.actor ON ((film_actor.actor_id = actor.actor_id)))
  GROUP BY film.film_id, film.title, film.description, category.name, film.rental_rate, film.length, film.rating;


ALTER TABLE public.film_list OWNER TO mockdata;

--
-- Name: inventory_inventory_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.inventory_inventory_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.inventory_inventory_id_seq OWNER TO mockdata;

--
-- Name: inventory; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.inventory (
    inventory_id integer DEFAULT nextval('public.inventory_inventory_id_seq'::regclass) NOT NULL,
    film_id smallint NOT NULL,
    store_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.inventory OWNER TO mockdata;

--
-- Name: language_language_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.language_language_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.language_language_id_seq OWNER TO mockdata;

--
-- Name: language; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.language (
    language_id integer DEFAULT nextval('public.language_language_id_seq'::regclass) NOT NULL,
    name character(20) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.language OWNER TO mockdata;

--
-- Name: nicer_but_slower_film_list; Type: VIEW; Schema: public; Owner: mockdata
--

CREATE VIEW public.nicer_but_slower_film_list AS
 SELECT film.film_id AS fid,
    film.title,
    film.description,
    category.name AS category,
    film.rental_rate AS price,
    film.length,
    film.rating,
    public.group_concat((((upper("substring"((actor.first_name)::text, 1, 1)) || lower("substring"((actor.first_name)::text, 2))) || upper("substring"((actor.last_name)::text, 1, 1))) || lower("substring"((actor.last_name)::text, 2)))) AS actors
   FROM ((((public.category
     LEFT JOIN public.film_category ON ((category.category_id = film_category.category_id)))
     LEFT JOIN public.film ON ((film_category.film_id = film.film_id)))
     JOIN public.film_actor ON ((film.film_id = film_actor.film_id)))
     JOIN public.actor ON ((film_actor.actor_id = actor.actor_id)))
  GROUP BY film.film_id, film.title, film.description, category.name, film.rental_rate, film.length, film.rating;


ALTER TABLE public.nicer_but_slower_film_list OWNER TO mockdata;

--
-- Name: payment_payment_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.payment_payment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.payment_payment_id_seq OWNER TO mockdata;

--
-- Name: payment; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.payment (
    payment_id integer DEFAULT nextval('public.payment_payment_id_seq'::regclass) NOT NULL,
    customer_id smallint NOT NULL,
    staff_id smallint NOT NULL,
    rental_id integer NOT NULL,
    amount numeric(5,2) NOT NULL,
    payment_date timestamp without time zone NOT NULL
);


ALTER TABLE public.payment OWNER TO mockdata;

--
-- Name: rental_rental_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.rental_rental_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rental_rental_id_seq OWNER TO mockdata;

--
-- Name: rental; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.rental (
    rental_id integer DEFAULT nextval('public.rental_rental_id_seq'::regclass) NOT NULL,
    rental_date timestamp without time zone NOT NULL,
    inventory_id integer NOT NULL,
    customer_id smallint NOT NULL,
    return_date timestamp without time zone,
    staff_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.rental OWNER TO mockdata;

--
-- Name: sales_by_film_category; Type: VIEW; Schema: public; Owner: mockdata
--

CREATE VIEW public.sales_by_film_category AS
 SELECT c.name AS category,
    sum(p.amount) AS total_sales
   FROM (((((public.payment p
     JOIN public.rental r ON ((p.rental_id = r.rental_id)))
     JOIN public.inventory i ON ((r.inventory_id = i.inventory_id)))
     JOIN public.film f ON ((i.film_id = f.film_id)))
     JOIN public.film_category fc ON ((f.film_id = fc.film_id)))
     JOIN public.category c ON ((fc.category_id = c.category_id)))
  GROUP BY c.name
  ORDER BY (sum(p.amount)) DESC;


ALTER TABLE public.sales_by_film_category OWNER TO mockdata;

--
-- Name: staff_staff_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.staff_staff_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.staff_staff_id_seq OWNER TO mockdata;

--
-- Name: staff; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.staff (
    staff_id integer DEFAULT nextval('public.staff_staff_id_seq'::regclass) NOT NULL,
    first_name character varying(45) NOT NULL,
    last_name character varying(45) NOT NULL,
    address_id smallint NOT NULL,
    email character varying(50),
    store_id smallint NOT NULL,
    active boolean DEFAULT true NOT NULL,
    username character varying(16) NOT NULL,
    password character varying(40),
    last_update timestamp without time zone DEFAULT now() NOT NULL,
    picture bytea
);


ALTER TABLE public.staff OWNER TO mockdata;

--
-- Name: store_store_id_seq; Type: SEQUENCE; Schema: public; Owner: mockdata
--

CREATE SEQUENCE public.store_store_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.store_store_id_seq OWNER TO mockdata;

--
-- Name: store; Type: TABLE; Schema: public; Owner: mockdata
--

CREATE TABLE public.store (
    store_id integer DEFAULT nextval('public.store_store_id_seq'::regclass) NOT NULL,
    manager_staff_id smallint NOT NULL,
    address_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.store OWNER TO mockdata;

--
-- Name: sales_by_store; Type: VIEW; Schema: public; Owner: mockdata
--

CREATE VIEW public.sales_by_store AS
 SELECT (((c.city)::text || ','::text) || (cy.country)::text) AS store,
    (((m.first_name)::text || ' '::text) || (m.last_name)::text) AS manager,
    sum(p.amount) AS total_sales
   FROM (((((((public.payment p
     JOIN public.rental r ON ((p.rental_id = r.rental_id)))
     JOIN public.inventory i ON ((r.inventory_id = i.inventory_id)))
     JOIN public.store s ON ((i.store_id = s.store_id)))
     JOIN public.address a ON ((s.address_id = a.address_id)))
     JOIN public.city c ON ((a.city_id = c.city_id)))
     JOIN public.country cy ON ((c.country_id = cy.country_id)))
     JOIN public.staff m ON ((s.manager_staff_id = m.staff_id)))
  GROUP BY cy.country, c.city, s.store_id, m.first_name, m.last_name
  ORDER BY cy.country, c.city;


ALTER TABLE public.sales_by_store OWNER TO mockdata;

--
-- Name: staff_list; Type: VIEW; Schema: public; Owner: mockdata
--

CREATE VIEW public.staff_list AS
 SELECT s.staff_id AS id,
    (((s.first_name)::text || ' '::text) || (s.last_name)::text) AS name,
    a.address,
    a.postal_code AS "zip code",
    a.phone,
    city.city,
    country.country,
    s.store_id AS sid
   FROM (((public.staff s
     JOIN public.address a ON ((s.address_id = a.address_id)))
     JOIN public.city ON ((a.city_id = city.city_id)))
     JOIN public.country ON ((city.country_id = country.country_id)));


ALTER TABLE public.staff_list OWNER TO mockdata;

--
-- Name: actor_actor_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.actor_actor_id_seq', 200, true);


--
-- Name: address_address_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.address_address_id_seq', 605, true);


--
-- Name: category_category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.category_category_id_seq', 16, true);


--
-- Name: city_city_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.city_city_id_seq', 600, true);


--
-- Name: country_country_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.country_country_id_seq', 109, true);


--
-- Name: customer_customer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.customer_customer_id_seq', 599, true);


--
-- Name: film_film_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.film_film_id_seq', 1000, true);


--
-- Name: inventory_inventory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.inventory_inventory_id_seq', 4581, true);


--
-- Name: language_language_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.language_language_id_seq', 6, true);


--
-- Name: payment_payment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.payment_payment_id_seq', 32098, true);


--
-- Name: rental_rental_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.rental_rental_id_seq', 16049, true);


--
-- Name: staff_staff_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.staff_staff_id_seq', 2, true);


--
-- Name: store_store_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mockdata
--

SELECT pg_catalog.setval('public.store_store_id_seq', 2, true);


--
-- Name: actor actor_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.actor
    ADD CONSTRAINT actor_pkey PRIMARY KEY (actor_id);


--
-- Name: address address_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_pkey PRIMARY KEY (address_id);


--
-- Name: category category_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_pkey PRIMARY KEY (category_id);


--
-- Name: city city_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT city_pkey PRIMARY KEY (city_id);


--
-- Name: country country_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_pkey PRIMARY KEY (country_id);


--
-- Name: customer customer_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_pkey PRIMARY KEY (customer_id);


--
-- Name: film_actor film_actor_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.film_actor
    ADD CONSTRAINT film_actor_pkey PRIMARY KEY (actor_id, film_id);


--
-- Name: film_category film_category_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.film_category
    ADD CONSTRAINT film_category_pkey PRIMARY KEY (film_id, category_id);


--
-- Name: film film_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.film
    ADD CONSTRAINT film_pkey PRIMARY KEY (film_id);


--
-- Name: inventory inventory_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.inventory
    ADD CONSTRAINT inventory_pkey PRIMARY KEY (inventory_id);


--
-- Name: language language_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.language
    ADD CONSTRAINT language_pkey PRIMARY KEY (language_id);


--
-- Name: payment payment_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_pkey PRIMARY KEY (payment_id);


--
-- Name: rental rental_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.rental
    ADD CONSTRAINT rental_pkey PRIMARY KEY (rental_id);


--
-- Name: staff staff_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.staff
    ADD CONSTRAINT staff_pkey PRIMARY KEY (staff_id);


--
-- Name: store store_pkey; Type: CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.store
    ADD CONSTRAINT store_pkey PRIMARY KEY (store_id);


--
-- Name: film_fulltext_idx; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX film_fulltext_idx ON public.film USING gist (fulltext);


--
-- Name: idx_actor_last_name; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_actor_last_name ON public.actor USING btree (last_name);


--
-- Name: idx_fk_address_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_fk_address_id ON public.customer USING btree (address_id);


--
-- Name: idx_fk_city_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_fk_city_id ON public.address USING btree (city_id);


--
-- Name: idx_fk_country_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_fk_country_id ON public.city USING btree (country_id);


--
-- Name: idx_fk_customer_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_fk_customer_id ON public.payment USING btree (customer_id);


--
-- Name: idx_fk_film_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_fk_film_id ON public.film_actor USING btree (film_id);


--
-- Name: idx_fk_inventory_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_fk_inventory_id ON public.rental USING btree (inventory_id);


--
-- Name: idx_fk_language_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_fk_language_id ON public.film USING btree (language_id);


--
-- Name: idx_fk_rental_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_fk_rental_id ON public.payment USING btree (rental_id);


--
-- Name: idx_fk_staff_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_fk_staff_id ON public.payment USING btree (staff_id);


--
-- Name: idx_fk_store_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_fk_store_id ON public.customer USING btree (store_id);


--
-- Name: idx_last_name; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_last_name ON public.customer USING btree (last_name);


--
-- Name: idx_store_id_film_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_store_id_film_id ON public.inventory USING btree (store_id, film_id);


--
-- Name: idx_title; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE INDEX idx_title ON public.film USING btree (title);


--
-- Name: idx_unq_manager_staff_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE UNIQUE INDEX idx_unq_manager_staff_id ON public.store USING btree (manager_staff_id);


--
-- Name: idx_unq_rental_rental_date_inventory_id_customer_id; Type: INDEX; Schema: public; Owner: mockdata
--

CREATE UNIQUE INDEX idx_unq_rental_rental_date_inventory_id_customer_id ON public.rental USING btree (rental_date, inventory_id, customer_id);


--
-- Name: film film_fulltext_trigger; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER film_fulltext_trigger BEFORE INSERT OR UPDATE ON public.film FOR EACH ROW EXECUTE PROCEDURE tsvector_update_trigger('fulltext', 'pg_catalog.english', 'title', 'description');


--
-- Name: actor last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.actor FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: address last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.address FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: category last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.category FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: city last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.city FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: country last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.country FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: customer last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.customer FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: film last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.film FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: film_actor last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.film_actor FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: film_category last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.film_category FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: inventory last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.inventory FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: language last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.language FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: rental last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.rental FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: staff last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.staff FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: store last_updated; Type: TRIGGER; Schema: public; Owner: mockdata
--

CREATE TRIGGER last_updated BEFORE UPDATE ON public.store FOR EACH ROW EXECUTE PROCEDURE public.last_updated();


--
-- Name: customer customer_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(address_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: film_actor film_actor_actor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.film_actor
    ADD CONSTRAINT film_actor_actor_id_fkey FOREIGN KEY (actor_id) REFERENCES public.actor(actor_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: film_actor film_actor_film_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.film_actor
    ADD CONSTRAINT film_actor_film_id_fkey FOREIGN KEY (film_id) REFERENCES public.film(film_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: film_category film_category_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.film_category
    ADD CONSTRAINT film_category_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.category(category_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: film_category film_category_film_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.film_category
    ADD CONSTRAINT film_category_film_id_fkey FOREIGN KEY (film_id) REFERENCES public.film(film_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: film film_language_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.film
    ADD CONSTRAINT film_language_id_fkey FOREIGN KEY (language_id) REFERENCES public.language(language_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: address fk_address_city; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT fk_address_city FOREIGN KEY (city_id) REFERENCES public.city(city_id);


--
-- Name: city fk_city; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT fk_city FOREIGN KEY (country_id) REFERENCES public.country(country_id);


--
-- Name: inventory inventory_film_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.inventory
    ADD CONSTRAINT inventory_film_id_fkey FOREIGN KEY (film_id) REFERENCES public.film(film_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: payment payment_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customer(customer_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: payment payment_rental_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_rental_id_fkey FOREIGN KEY (rental_id) REFERENCES public.rental(rental_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: payment payment_staff_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_staff_id_fkey FOREIGN KEY (staff_id) REFERENCES public.staff(staff_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: rental rental_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.rental
    ADD CONSTRAINT rental_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customer(customer_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: rental rental_inventory_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.rental
    ADD CONSTRAINT rental_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES public.inventory(inventory_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: rental rental_staff_id_key; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.rental
    ADD CONSTRAINT rental_staff_id_key FOREIGN KEY (staff_id) REFERENCES public.staff(staff_id);


--
-- Name: staff staff_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.staff
    ADD CONSTRAINT staff_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(address_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: store store_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.store
    ADD CONSTRAINT store_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(address_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: store store_manager_staff_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mockdata
--

ALTER TABLE ONLY public.store
    ADD CONSTRAINT store_manager_staff_id_fkey FOREIGN KEY (manager_staff_id) REFERENCES public.staff(staff_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--
`
}

// Greenplum demo database
func demoDatabaseGreenplum() string {
	Debug("Creating a demo database for greenplum")
	return `
-- Sample database download from
-- http://www.gpadminqltutorial.com/gpadminql-sample-database/
--
-- NOTE:
--
-- File paths need to be edited. Search for $$PATH$$ and
-- replace it with the path to the directory containing
-- the extracted data files.
--
--
-- gpadminQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

--
-- DROP ALL THE OBJECTS 
--
ALTER TABLE IF EXISTS ONLY public.store DROP CONSTRAINT store_manager_staff_id_fkey;
ALTER TABLE IF EXISTS ONLY public.store DROP CONSTRAINT store_address_id_fkey;
ALTER TABLE IF EXISTS ONLY public.staff DROP CONSTRAINT staff_address_id_fkey;
ALTER TABLE IF EXISTS ONLY public.rental DROP CONSTRAINT rental_staff_id_key;
ALTER TABLE IF EXISTS ONLY public.rental DROP CONSTRAINT rental_inventory_id_fkey;
ALTER TABLE IF EXISTS ONLY public.rental DROP CONSTRAINT rental_customer_id_fkey;
ALTER TABLE IF EXISTS ONLY public.payment DROP CONSTRAINT payment_staff_id_fkey;
ALTER TABLE IF EXISTS ONLY public.payment DROP CONSTRAINT payment_rental_id_fkey;
ALTER TABLE IF EXISTS ONLY public.payment DROP CONSTRAINT payment_customer_id_fkey;
ALTER TABLE IF EXISTS ONLY public.inventory DROP CONSTRAINT inventory_film_id_fkey;
ALTER TABLE IF EXISTS ONLY public.city DROP CONSTRAINT fk_city;
ALTER TABLE IF EXISTS ONLY public.address DROP CONSTRAINT fk_address_city;
ALTER TABLE IF EXISTS ONLY public.film DROP CONSTRAINT film_language_id_fkey;
ALTER TABLE IF EXISTS ONLY public.film_category DROP CONSTRAINT film_category_film_id_fkey;
ALTER TABLE IF EXISTS ONLY public.film_category DROP CONSTRAINT film_category_category_id_fkey;
ALTER TABLE IF EXISTS ONLY public.film_actor DROP CONSTRAINT film_actor_film_id_fkey;
ALTER TABLE IF EXISTS ONLY public.film_actor DROP CONSTRAINT film_actor_actor_id_fkey;
ALTER TABLE IF EXISTS ONLY public.customer DROP CONSTRAINT customer_address_id_fkey;
DROP TRIGGER IF EXISTS last_updated ON public.store;
DROP TRIGGER IF EXISTS last_updated ON public.staff;
DROP TRIGGER IF EXISTS last_updated ON public.rental;
DROP TRIGGER IF EXISTS last_updated ON public.language;
DROP TRIGGER IF EXISTS last_updated ON public.inventory;
DROP TRIGGER IF EXISTS last_updated ON public.film_category;
DROP TRIGGER IF EXISTS last_updated ON public.film_actor;
DROP TRIGGER IF EXISTS last_updated ON public.film;
DROP TRIGGER IF EXISTS last_updated ON public.customer;
DROP TRIGGER IF EXISTS last_updated ON public.country;
DROP TRIGGER IF EXISTS last_updated ON public.city;
DROP TRIGGER IF EXISTS last_updated ON public.category;
DROP TRIGGER IF EXISTS last_updated ON public.address;
DROP TRIGGER IF EXISTS last_updated ON public.actor;
DROP TRIGGER IF EXISTS film_fulltext_trigger ON public.film;
DROP INDEX IF EXISTS public.idx_unq_rental_rental_date_inventory_id_customer_id;
DROP INDEX IF EXISTS public.idx_unq_manager_staff_id;
DROP INDEX IF EXISTS public.idx_title;
DROP INDEX IF EXISTS public.idx_store_id_film_id;
DROP INDEX IF EXISTS public.idx_last_name;
DROP INDEX IF EXISTS public.idx_fk_store_id;
DROP INDEX IF EXISTS public.idx_fk_staff_id;
DROP INDEX IF EXISTS public.idx_fk_rental_id;
DROP INDEX IF EXISTS public.idx_fk_language_id;
DROP INDEX IF EXISTS public.idx_fk_inventory_id;
DROP INDEX IF EXISTS public.idx_fk_film_id;
DROP INDEX IF EXISTS public.idx_fk_customer_id;
DROP INDEX IF EXISTS public.idx_fk_country_id;
DROP INDEX IF EXISTS public.idx_fk_city_id;
DROP INDEX IF EXISTS public.idx_fk_address_id;
DROP INDEX IF EXISTS public.idx_actor_last_name;
DROP INDEX IF EXISTS public.film_fulltext_idx;
ALTER TABLE IF EXISTS ONLY public.store DROP CONSTRAINT store_pkey;
ALTER TABLE IF EXISTS ONLY public.staff DROP CONSTRAINT staff_pkey;
ALTER TABLE IF EXISTS ONLY public.rental DROP CONSTRAINT rental_pkey;
ALTER TABLE IF EXISTS ONLY public.payment DROP CONSTRAINT payment_pkey;
ALTER TABLE IF EXISTS ONLY public.language DROP CONSTRAINT language_pkey;
ALTER TABLE IF EXISTS ONLY public.inventory DROP CONSTRAINT inventory_pkey;
ALTER TABLE IF EXISTS ONLY public.film DROP CONSTRAINT film_pkey;
ALTER TABLE IF EXISTS ONLY public.film_category DROP CONSTRAINT film_category_pkey;
ALTER TABLE IF EXISTS ONLY public.film_actor DROP CONSTRAINT film_actor_pkey;
ALTER TABLE IF EXISTS ONLY public.customer DROP CONSTRAINT customer_pkey;
ALTER TABLE IF EXISTS ONLY public.country DROP CONSTRAINT country_pkey;
ALTER TABLE IF EXISTS ONLY public.city DROP CONSTRAINT city_pkey;
ALTER TABLE IF EXISTS ONLY public.category DROP CONSTRAINT category_pkey;
ALTER TABLE IF EXISTS ONLY public.address DROP CONSTRAINT address_pkey;
ALTER TABLE IF EXISTS ONLY public.actor DROP CONSTRAINT actor_pkey;
DROP VIEW IF EXISTS public.staff_list;
DROP VIEW IF EXISTS public.sales_by_store;
DROP TABLE IF EXISTS public.store;
DROP SEQUENCE IF EXISTS public.store_store_id_seq;
DROP TABLE IF EXISTS public.staff;
DROP SEQUENCE IF EXISTS public.staff_staff_id_seq;
DROP VIEW IF EXISTS public.sales_by_film_category;
DROP TABLE IF EXISTS public.rental;
DROP SEQUENCE IF EXISTS public.rental_rental_id_seq;
DROP TABLE IF EXISTS public.payment;
DROP SEQUENCE IF EXISTS public.payment_payment_id_seq;
DROP TABLE IF EXISTS public.language;
DROP SEQUENCE IF EXISTS public.language_language_id_seq;
DROP TABLE IF EXISTS public.inventory;
DROP SEQUENCE IF EXISTS public.inventory_inventory_id_seq;
DROP VIEW IF EXISTS public.film_list;
DROP VIEW IF EXISTS public.customer_list;
DROP TABLE IF EXISTS public.country;
DROP SEQUENCE IF EXISTS public.country_country_id_seq;
DROP TABLE IF EXISTS public.city;
DROP SEQUENCE IF EXISTS public.city_city_id_seq;
DROP TABLE IF EXISTS public.address;
DROP SEQUENCE IF EXISTS public.address_address_id_seq;
DROP VIEW IF EXISTS public.actor_info;
DROP TABLE IF EXISTS public.film_category;
DROP TABLE IF EXISTS public.film_actor;
DROP TABLE IF EXISTS public.film;
DROP SEQUENCE IF EXISTS public.film_film_id_seq;
DROP TABLE IF EXISTS public.category;
DROP SEQUENCE IF EXISTS public.category_category_id_seq;
DROP TABLE IF EXISTS public.actor;
DROP SEQUENCE IF EXISTS public.actor_actor_id_seq;
DROP FUNCTION IF EXISTS public.rewards_report(min_monthly_purchases integer, min_dollar_amount_purchased numeric);
DROP TABLE IF EXISTS public.customer;
DROP SEQUENCE IF EXISTS public.customer_customer_id_seq;
DROP FUNCTION IF EXISTS public.last_updated();
DROP FUNCTION IF EXISTS public.last_day(timestamp without time zone);
DROP FUNCTION IF EXISTS public.inventory_in_stock(p_inventory_id integer);
DROP FUNCTION IF EXISTS public.inventory_held_by_customer(p_inventory_id integer);
DROP FUNCTION IF EXISTS public.get_customer_balance(p_customer_id integer, p_effective_date timestamp without time zone);
DROP FUNCTION IF EXISTS public.film_not_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer);
DROP FUNCTION IF EXISTS public.film_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer);
DROP DOMAIN IF EXISTS public.year;
DROP TYPE IF EXISTS public.mpaa_rating;
DROP TYPE IF EXISTS rating;
CREATE TYPE rating as ENUM ('good', 'ok', 'bad');

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: gpadmin
--

CREATE SCHEMA IF NOT EXISTS public;


ALTER SCHEMA public OWNER TO gpadmin;

--
-- Name: actor_actor_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE actor_actor_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.actor_actor_id_seq OWNER TO gpadmin;

--
-- Name: actor; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE actor (
    actor_id integer DEFAULT nextval('actor_actor_id_seq'::regclass) NOT NULL,
    first_name character varying(15) NOT NULL,
    last_name character varying(15) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL,
    email char(20) unique,
    gender char(1),
    rate integer
);


ALTER TABLE public.actor OWNER TO gpadmin;

--
-- Name: category_category_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE category_category_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.category_category_id_seq OWNER TO gpadmin;

--
-- Name: category; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE category (
    category_id integer DEFAULT nextval('category_category_id_seq'::regclass) NOT NULL,
    name character varying(25) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.category OWNER TO gpadmin;

--
-- Name: film_film_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE film_film_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.film_film_id_seq OWNER TO gpadmin;

--
-- Name: film; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE film (
    film_id integer DEFAULT nextval('film_film_id_seq'::regclass) NOT NULL,
    title character varying(255) NOT NULL,
    description text,
    release_year date,
    language_id smallint NOT NULL,
    rental_duration smallint DEFAULT 3 NOT NULL,
    rental_rate numeric(4,2) DEFAULT 4.99 NOT NULL,
    length smallint,
    replacement_cost numeric(5,2) DEFAULT 19.99 NOT NULL,
    rating char(1) DEFAULT 'G',
	user_rating rating,
    last_update timestamp without time zone DEFAULT now() NOT NULL,
    special_features text,
    fulltext tsvector NOT NULL
);


ALTER TABLE public.film OWNER TO gpadmin;

--
-- Name: film_actor; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE film_actor (
    actor_id smallint NOT NULL,
    film_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.film_actor OWNER TO gpadmin;

--
-- Name: film_category; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE film_category (
    film_id smallint NOT NULL,
    category_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.film_category OWNER TO gpadmin;

--
-- Name: address_address_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE address_address_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.address_address_id_seq OWNER TO gpadmin;

--
-- Name: address; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE address (
    address_id integer DEFAULT nextval('address_address_id_seq'::regclass) NOT NULL,
    address character varying(50) NOT NULL,
    address2 character varying(50),
    district character varying(20) NOT NULL,
    city_id smallint NOT NULL,
    postal_code character varying(10),
    phone character varying(20) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.address OWNER TO gpadmin;

--
-- Name: city_city_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE city_city_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.city_city_id_seq OWNER TO gpadmin;

--
-- Name: city; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE city (
    city_id integer DEFAULT nextval('city_city_id_seq'::regclass) NOT NULL,
    city character varying(50) NOT NULL,
    country_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.city OWNER TO gpadmin;

--
-- Name: country_country_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE country_country_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.country_country_id_seq OWNER TO gpadmin;

--
-- Name: country; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE country (
    country_id integer DEFAULT nextval('country_country_id_seq'::regclass) NOT NULL,
    country character varying(50) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.country OWNER TO gpadmin;

--
-- Name: inventory_inventory_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE inventory_inventory_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.inventory_inventory_id_seq OWNER TO gpadmin;

--
-- Name: inventory; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE inventory (
    inventory_id integer DEFAULT nextval('inventory_inventory_id_seq'::regclass) NOT NULL,
    film_id smallint NOT NULL,
    store_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.inventory OWNER TO gpadmin;

--
-- Name: language_language_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE language_language_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.language_language_id_seq OWNER TO gpadmin;

--
-- Name: language; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE language (
    language_id integer DEFAULT nextval('language_language_id_seq'::regclass) NOT NULL,
    name character(20) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.language OWNER TO gpadmin;

--
-- Name: payment_payment_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE payment_payment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.payment_payment_id_seq OWNER TO gpadmin;

--
-- Name: payment; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE payment (
    payment_id integer DEFAULT nextval('payment_payment_id_seq'::regclass) NOT NULL,
    customer_id smallint NOT NULL,
    staff_id smallint NOT NULL,
    rental_id integer NOT NULL,
    amount numeric(5,2) NOT NULL,
    payment_date timestamp without time zone NOT NULL
);


ALTER TABLE public.payment OWNER TO gpadmin;

--
-- Name: rental_rental_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE rental_rental_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rental_rental_id_seq OWNER TO gpadmin;

--
-- Name: rental; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE rental (
    rental_id integer DEFAULT nextval('rental_rental_id_seq'::regclass) NOT NULL,
    rental_date timestamp without time zone NOT NULL,
    inventory_id integer NOT NULL,
    customer_id smallint NOT NULL,
    return_date timestamp without time zone,
    staff_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.rental OWNER TO gpadmin;

--
-- Name: sales_by_film_category; Type: VIEW; Schema: public; Owner: gpadmin
--

CREATE VIEW sales_by_film_category AS
    SELECT c.name AS category, sum(p.amount) AS total_sales FROM (((((payment p JOIN rental r ON ((p.rental_id = r.rental_id))) JOIN inventory i ON ((r.inventory_id = i.inventory_id))) JOIN film f ON ((i.film_id = f.film_id))) JOIN film_category fc ON ((f.film_id = fc.film_id))) JOIN category c ON ((fc.category_id = c.category_id))) GROUP BY c.name ORDER BY sum(p.amount) DESC;


ALTER TABLE public.sales_by_film_category OWNER TO gpadmin;

--
-- Name: staff_staff_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE staff_staff_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.staff_staff_id_seq OWNER TO gpadmin;

--
-- Name: staff; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE staff (
    staff_id integer DEFAULT nextval('staff_staff_id_seq'::regclass) NOT NULL,
    first_name character varying(45) NOT NULL,
    last_name character varying(45) NOT NULL,
    address_id smallint NOT NULL,
    email character varying(50),
    store_id smallint NOT NULL,
    active boolean DEFAULT true NOT NULL,
    username character varying(16) NOT NULL,
    password character varying(40),
    last_update timestamp without time zone DEFAULT now() NOT NULL,
    picture bytea
);


ALTER TABLE public.staff OWNER TO gpadmin;

--
-- Name: store_store_id_seq; Type: SEQUENCE; Schema: public; Owner: gpadmin
--

CREATE SEQUENCE store_store_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.store_store_id_seq OWNER TO gpadmin;

--
-- Name: store; Type: TABLE; Schema: public; Owner: gpadmin; Tablespace: 
--

CREATE TABLE store (
    store_id integer DEFAULT nextval('store_store_id_seq'::regclass) NOT NULL,
    manager_staff_id smallint NOT NULL,
    address_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.store OWNER TO gpadmin;

--
-- Name: public; Type: ACL; Schema: -; Owner: gpadmin
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM gpadmin;
GRANT ALL ON SCHEMA public TO gpadmin;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- gpadminQL database dump complete
--
`
}
