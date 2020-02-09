package main

// Json Skeleton
func JsonSkeleton() string {
	return `{"_id": "%s","index": "%s","guid": "%s","isActive": "%s","balance": "%s.%s","website": "https://%s/%s","age": "%s","username": "%s","eyeColor": "%s","name": "%s","gender": "%s","company": "%s","email": "%s","phone": "%s","address": "%s","zipcode": "%s","state": "%s","country": "%s","about": "%s","Machine IP": "%s","job title": "%s","registered": "%s-%s-%sT%s:%s:%s-%s:%s","latitude": "%s.%s","longitude": "%s.%s","tags": ["%s","%s","%s","%s","%s","%s","%s"],"friends": [{ "id": "%s", "name": "%s"},{ "id": "%s", "name": "%s"},{ "id": "%s", "name": "%s"}],"greeting": "%s","favoriteBrand": "%s"}`
}

// XML Skeleton
func XMLSkeleton() string {
	return `<?xml version="1.0" encoding="UTF-8"?><shiporder orderid="%s" xmlns:xsi="http://%s/%s/%s" xsi:noNamespaceSchemaLocation="shiporder.xsd"> <orderperson>%s</orderperson> <shipto> <name>%s</name> <address>%s</address> <city>%s</city> <country>%s</country> <email>%s</email> <phone>%s</phone> </shipto> <item> <title>%s</title> <note>%s</note> <quantity>%s</quantity> <color>%s</color> <price>%s.%s</price> </item> <item> <title>%s</title> <quantity>%s</quantity> <price>%s.%s</price> </item></shiporder>`
}

// Demo database sql
func demoDatabase() string {
	Debug("Creating a demo database")
	return `

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

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA IF NOT EXISTS public;


ALTER SCHEMA public OWNER TO postgres;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
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
-- Name: mpaa_rating; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE mpaa_rating AS ENUM (
    'G',
    'PG',
    'PG-13',
    'R',
    'NC-17'
);


ALTER TYPE public.mpaa_rating OWNER TO postgres;

--
-- Name: year; Type: DOMAIN; Schema: public; Owner: postgres
--

CREATE DOMAIN year AS integer
	CONSTRAINT year_check CHECK (((VALUE >= 1901) AND (VALUE <= 2155)));


ALTER DOMAIN public.year OWNER TO postgres;

--
-- Name: _group_concat(text, text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION _group_concat(text, text) RETURNS text
    LANGUAGE sql IMMUTABLE
    AS $_$
SELECT CASE
  WHEN $2 IS NULL THEN $1
  WHEN $1 IS NULL THEN $2
  ELSE $1 || ', ' || $2
END
$_$;


ALTER FUNCTION public._group_concat(text, text) OWNER TO postgres;

--
-- Name: film_in_stock(integer, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION film_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer) RETURNS SETOF integer
    LANGUAGE sql
    AS $_$
     SELECT inventory_id
     FROM inventory
     WHERE film_id = $1
     AND store_id = $2
     AND inventory_in_stock(inventory_id);
$_$;


ALTER FUNCTION public.film_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer) OWNER TO postgres;

--
-- Name: film_not_in_stock(integer, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION film_not_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer) RETURNS SETOF integer
    LANGUAGE sql
    AS $_$
    SELECT inventory_id
    FROM inventory
    WHERE film_id = $1
    AND store_id = $2
    AND NOT inventory_in_stock(inventory_id);
$_$;


ALTER FUNCTION public.film_not_in_stock(p_film_id integer, p_store_id integer, OUT p_film_count integer) OWNER TO postgres;

--
-- Name: get_customer_balance(integer, timestamp without time zone); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION get_customer_balance(p_customer_id integer, p_effective_date timestamp without time zone) RETURNS numeric
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


ALTER FUNCTION public.get_customer_balance(p_customer_id integer, p_effective_date timestamp without time zone) OWNER TO postgres;

--
-- Name: inventory_held_by_customer(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION inventory_held_by_customer(p_inventory_id integer) RETURNS integer
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


ALTER FUNCTION public.inventory_held_by_customer(p_inventory_id integer) OWNER TO postgres;

--
-- Name: inventory_in_stock(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION inventory_in_stock(p_inventory_id integer) RETURNS boolean
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


ALTER FUNCTION public.inventory_in_stock(p_inventory_id integer) OWNER TO postgres;

--
-- Name: last_day(timestamp without time zone); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION last_day(timestamp without time zone) RETURNS date
    LANGUAGE sql IMMUTABLE STRICT
    AS $_$
  SELECT CASE
    WHEN EXTRACT(MONTH FROM $1) = 12 THEN
      (((EXTRACT(YEAR FROM $1) + 1) operator(pg_catalog.||) '-01-01')::date - INTERVAL '1 day')::date
    ELSE
      ((EXTRACT(YEAR FROM $1) operator(pg_catalog.||) '-' operator(pg_catalog.||) (EXTRACT(MONTH FROM $1) + 1) operator(pg_catalog.||) '-01')::date - INTERVAL '1 day')::date
    END
$_$;


ALTER FUNCTION public.last_day(timestamp without time zone) OWNER TO postgres;

--
-- Name: last_updated(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION last_updated() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.last_update = CURRENT_TIMESTAMP;
    RETURN NEW;
END $$;


ALTER FUNCTION public.last_updated() OWNER TO postgres;

--
-- Name: customer_customer_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE customer_customer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.customer_customer_id_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: customer; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE customer (
    customer_id integer DEFAULT nextval('customer_customer_id_seq'::regclass) NOT NULL,
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


ALTER TABLE public.customer OWNER TO postgres;

--
-- Name: rewards_report(integer, numeric); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION rewards_report(min_monthly_purchases integer, min_dollar_amount_purchased numeric) RETURNS SETOF customer
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


ALTER FUNCTION public.rewards_report(min_monthly_purchases integer, min_dollar_amount_purchased numeric) OWNER TO postgres;

--
-- Name: group_concat(text); Type: AGGREGATE; Schema: public; Owner: postgres
--

CREATE AGGREGATE group_concat(text) (
    SFUNC = _group_concat,
    STYPE = text
);


ALTER AGGREGATE public.group_concat(text) OWNER TO postgres;

--
-- Name: actor_actor_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE actor_actor_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.actor_actor_id_seq OWNER TO postgres;

--
-- Name: actor; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
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


ALTER TABLE public.actor OWNER TO postgres;

--
-- Name: category_category_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE category_category_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.category_category_id_seq OWNER TO postgres;

--
-- Name: category; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE category (
    category_id integer DEFAULT nextval('category_category_id_seq'::regclass) NOT NULL,
    name character varying(25) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.category OWNER TO postgres;

--
-- Name: film_film_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE film_film_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.film_film_id_seq OWNER TO postgres;

--
-- Name: film; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
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
    last_update timestamp without time zone DEFAULT now() NOT NULL,
    special_features text,
    fulltext tsvector NOT NULL
);


ALTER TABLE public.film OWNER TO postgres;

--
-- Name: film_actor; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE film_actor (
    actor_id smallint NOT NULL,
    film_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.film_actor OWNER TO postgres;

--
-- Name: film_category; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE film_category (
    film_id smallint NOT NULL,
    category_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.film_category OWNER TO postgres;

--
-- Name: actor_info; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW actor_info AS
    SELECT a.actor_id, a.first_name, a.last_name, group_concat(DISTINCT (((c.name)::text || ': '::text) || (SELECT group_concat((f.title)::text) AS group_concat FROM ((film f JOIN film_category fc ON ((f.film_id = fc.film_id))) JOIN film_actor fa ON ((f.film_id = fa.film_id))) WHERE ((fc.category_id = c.category_id) AND (fa.actor_id = a.actor_id)) GROUP BY fa.actor_id))) AS film_info FROM (((actor a LEFT JOIN film_actor fa ON ((a.actor_id = fa.actor_id))) LEFT JOIN film_category fc ON ((fa.film_id = fc.film_id))) LEFT JOIN category c ON ((fc.category_id = c.category_id))) GROUP BY a.actor_id, a.first_name, a.last_name;


ALTER TABLE public.actor_info OWNER TO postgres;

--
-- Name: address_address_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE address_address_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.address_address_id_seq OWNER TO postgres;

--
-- Name: address; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
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


ALTER TABLE public.address OWNER TO postgres;

--
-- Name: city_city_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE city_city_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.city_city_id_seq OWNER TO postgres;

--
-- Name: city; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE city (
    city_id integer DEFAULT nextval('city_city_id_seq'::regclass) NOT NULL,
    city character varying(50) NOT NULL,
    country_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.city OWNER TO postgres;

--
-- Name: country_country_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE country_country_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.country_country_id_seq OWNER TO postgres;

--
-- Name: country; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE country (
    country_id integer DEFAULT nextval('country_country_id_seq'::regclass) NOT NULL,
    country character varying(50) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.country OWNER TO postgres;

--
-- Name: customer_list; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW customer_list AS
    SELECT cu.customer_id AS id, (((cu.first_name)::text || ' '::text) || (cu.last_name)::text) AS name, a.address, a.postal_code AS "zip code", a.phone, city.city, country.country, CASE WHEN cu.activebool THEN 'active'::text ELSE ''::text END AS notes, cu.store_id AS sid FROM (((customer cu JOIN address a ON ((cu.address_id = a.address_id))) JOIN city ON ((a.city_id = city.city_id))) JOIN country ON ((city.country_id = country.country_id)));


ALTER TABLE public.customer_list OWNER TO postgres;

--
-- Name: film_list; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW film_list AS
    SELECT film.film_id AS fid, film.title, film.description, category.name AS category, film.rental_rate AS price, film.length, film.rating, group_concat((((actor.first_name)::text || ' '::text) || (actor.last_name)::text)) AS actors FROM ((((category LEFT JOIN film_category ON ((category.category_id = film_category.category_id))) LEFT JOIN film ON ((film_category.film_id = film.film_id))) JOIN film_actor ON ((film.film_id = film_actor.film_id))) JOIN actor ON ((film_actor.actor_id = actor.actor_id))) GROUP BY film.film_id, film.title, film.description, category.name, film.rental_rate, film.length, film.rating;


ALTER TABLE public.film_list OWNER TO postgres;

--
-- Name: inventory_inventory_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE inventory_inventory_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.inventory_inventory_id_seq OWNER TO postgres;

--
-- Name: inventory; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE inventory (
    inventory_id integer DEFAULT nextval('inventory_inventory_id_seq'::regclass) NOT NULL,
    film_id smallint NOT NULL,
    store_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.inventory OWNER TO postgres;

--
-- Name: language_language_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE language_language_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.language_language_id_seq OWNER TO postgres;

--
-- Name: language; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE language (
    language_id integer DEFAULT nextval('language_language_id_seq'::regclass) NOT NULL,
    name character(20) NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.language OWNER TO postgres;

--
-- Name: nicer_but_slower_film_list; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW nicer_but_slower_film_list AS
    SELECT film.film_id AS fid, film.title, film.description, category.name AS category, film.rental_rate AS price, film.length, film.rating, group_concat((((upper("substring"((actor.first_name)::text, 1, 1)) || lower("substring"((actor.first_name)::text, 2))) || upper("substring"((actor.last_name)::text, 1, 1))) || lower("substring"((actor.last_name)::text, 2)))) AS actors FROM ((((category LEFT JOIN film_category ON ((category.category_id = film_category.category_id))) LEFT JOIN film ON ((film_category.film_id = film.film_id))) JOIN film_actor ON ((film.film_id = film_actor.film_id))) JOIN actor ON ((film_actor.actor_id = actor.actor_id))) GROUP BY film.film_id, film.title, film.description, category.name, film.rental_rate, film.length, film.rating;


ALTER TABLE public.nicer_but_slower_film_list OWNER TO postgres;

--
-- Name: payment_payment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE payment_payment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.payment_payment_id_seq OWNER TO postgres;

--
-- Name: payment; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE payment (
    payment_id integer DEFAULT nextval('payment_payment_id_seq'::regclass) NOT NULL,
    customer_id smallint NOT NULL,
    staff_id smallint NOT NULL,
    rental_id integer NOT NULL,
    amount numeric(5,2) NOT NULL,
    payment_date timestamp without time zone NOT NULL
);


ALTER TABLE public.payment OWNER TO postgres;

--
-- Name: rental_rental_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE rental_rental_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rental_rental_id_seq OWNER TO postgres;

--
-- Name: rental; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
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


ALTER TABLE public.rental OWNER TO postgres;

--
-- Name: sales_by_film_category; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW sales_by_film_category AS
    SELECT c.name AS category, sum(p.amount) AS total_sales FROM (((((payment p JOIN rental r ON ((p.rental_id = r.rental_id))) JOIN inventory i ON ((r.inventory_id = i.inventory_id))) JOIN film f ON ((i.film_id = f.film_id))) JOIN film_category fc ON ((f.film_id = fc.film_id))) JOIN category c ON ((fc.category_id = c.category_id))) GROUP BY c.name ORDER BY sum(p.amount) DESC;


ALTER TABLE public.sales_by_film_category OWNER TO postgres;

--
-- Name: staff_staff_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE staff_staff_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.staff_staff_id_seq OWNER TO postgres;

--
-- Name: staff; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
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


ALTER TABLE public.staff OWNER TO postgres;

--
-- Name: store_store_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE store_store_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.store_store_id_seq OWNER TO postgres;

--
-- Name: store; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE store (
    store_id integer DEFAULT nextval('store_store_id_seq'::regclass) NOT NULL,
    manager_staff_id smallint NOT NULL,
    address_id smallint NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.store OWNER TO postgres;

--
-- Name: sales_by_store; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW sales_by_store AS
    SELECT (((c.city)::text || ','::text) || (cy.country)::text) AS store, (((m.first_name)::text || ' '::text) || (m.last_name)::text) AS manager, sum(p.amount) AS total_sales FROM (((((((payment p JOIN rental r ON ((p.rental_id = r.rental_id))) JOIN inventory i ON ((r.inventory_id = i.inventory_id))) JOIN store s ON ((i.store_id = s.store_id))) JOIN address a ON ((s.address_id = a.address_id))) JOIN city c ON ((a.city_id = c.city_id))) JOIN country cy ON ((c.country_id = cy.country_id))) JOIN staff m ON ((s.manager_staff_id = m.staff_id))) GROUP BY cy.country, c.city, s.store_id, m.first_name, m.last_name ORDER BY cy.country, c.city;


ALTER TABLE public.sales_by_store OWNER TO postgres;

--
-- Name: staff_list; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW staff_list AS
    SELECT s.staff_id AS id, (((s.first_name)::text || ' '::text) || (s.last_name)::text) AS name, a.address, a.postal_code AS "zip code", a.phone, city.city, country.country, s.store_id AS sid FROM (((staff s JOIN address a ON ((s.address_id = a.address_id))) JOIN city ON ((a.city_id = city.city_id))) JOIN country ON ((city.country_id = country.country_id)));


ALTER TABLE public.staff_list OWNER TO postgres;


--
-- Name: actor_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY actor
    ADD CONSTRAINT actor_pkey PRIMARY KEY (actor_id);


ALTER TABLE ONLY actor
    ADD CONSTRAINT actor_ukey UNIQUE (first_name, last_name);

--
-- Name: actor_gender_ckey; Type: CK CONSTRAINT; Schema: public; Owner: actor
--

ALTER TABLE ONLY actor
    ADD CONSTRAINT actor_rate_ckey CHECK (rate > 0);

--
-- Name: address_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY address
    ADD CONSTRAINT address_pkey PRIMARY KEY (address_id);


--
-- Name: category_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY category
    ADD CONSTRAINT category_pkey PRIMARY KEY (category_id);


--
-- Name: city_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY city
    ADD CONSTRAINT city_pkey PRIMARY KEY (city_id);


--
-- Name: country_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY country
    ADD CONSTRAINT country_pkey PRIMARY KEY (country_id);


--
-- Name: customer_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY customer
    ADD CONSTRAINT customer_pkey PRIMARY KEY (customer_id);


--
-- Name: film_actor_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY film_actor
    ADD CONSTRAINT film_actor_pkey PRIMARY KEY (actor_id, film_id);


--
-- Name: film_category_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY film_category
    ADD CONSTRAINT film_category_pkey PRIMARY KEY (film_id, category_id);


--
-- Name: film_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY film
    ADD CONSTRAINT film_pkey PRIMARY KEY (film_id);


--
-- Name: inventory_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY inventory
    ADD CONSTRAINT inventory_pkey PRIMARY KEY (inventory_id);


--
-- Name: language_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY language
    ADD CONSTRAINT language_pkey PRIMARY KEY (language_id);


--
-- Name: payment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY payment
    ADD CONSTRAINT payment_pkey PRIMARY KEY (payment_id);


--
-- Name: rental_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY rental
    ADD CONSTRAINT rental_pkey PRIMARY KEY (rental_id);


--
-- Name: staff_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY staff
    ADD CONSTRAINT staff_pkey PRIMARY KEY (staff_id);


--
-- Name: store_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY store
    ADD CONSTRAINT store_pkey PRIMARY KEY (store_id);


--
-- Name: film_fulltext_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX film_fulltext_idx ON film USING gist (fulltext);


--
-- Name: idx_actor_last_name; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_actor_last_name ON actor USING btree (last_name);


--
-- Name: idx_fk_address_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_fk_address_id ON customer USING btree (address_id);


--
-- Name: idx_fk_city_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_fk_city_id ON address USING btree (city_id);


--
-- Name: idx_fk_country_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_fk_country_id ON city USING btree (country_id);


--
-- Name: idx_fk_customer_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_fk_customer_id ON payment USING btree (customer_id);


--
-- Name: idx_fk_film_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_fk_film_id ON film_actor USING btree (film_id);


--
-- Name: idx_fk_inventory_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_fk_inventory_id ON rental USING btree (inventory_id);


--
-- Name: idx_fk_language_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_fk_language_id ON film USING btree (language_id);


--
-- Name: idx_fk_rental_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_fk_rental_id ON payment USING btree (rental_id);


--
-- Name: idx_fk_staff_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_fk_staff_id ON payment USING btree (staff_id);


--
-- Name: idx_fk_store_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_fk_store_id ON customer USING btree (store_id);


--
-- Name: idx_last_name; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_last_name ON customer USING btree (last_name);


--
-- Name: idx_store_id_film_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_store_id_film_id ON inventory USING btree (store_id, film_id);


--
-- Name: idx_title; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX idx_title ON film USING btree (title);


--
-- Name: idx_unq_manager_staff_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX idx_unq_manager_staff_id ON store USING btree (manager_staff_id);


--
-- Name: idx_unq_rental_rental_date_inventory_id_customer_id; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX idx_unq_rental_rental_date_inventory_id_customer_id ON rental USING btree (rental_date, inventory_id, customer_id);


--
-- Name: film_fulltext_trigger; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER film_fulltext_trigger BEFORE INSERT OR UPDATE ON film FOR EACH ROW EXECUTE PROCEDURE tsvector_update_trigger('fulltext', 'pg_catalog.english', 'title', 'description');


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON actor FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON address FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON category FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON city FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON country FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON customer FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON film FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON film_actor FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON film_category FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON inventory FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON language FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON rental FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON staff FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: last_updated; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER last_updated BEFORE UPDATE ON store FOR EACH ROW EXECUTE PROCEDURE last_updated();


--
-- Name: customer_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY customer
    ADD CONSTRAINT customer_address_id_fkey FOREIGN KEY (address_id) REFERENCES address(address_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: film_actor_actor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY film_actor
    ADD CONSTRAINT film_actor_actor_id_fkey FOREIGN KEY (actor_id) REFERENCES actor(actor_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: film_actor_film_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY film_actor
    ADD CONSTRAINT film_actor_film_id_fkey FOREIGN KEY (film_id) REFERENCES film(film_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: film_category_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY film_category
    ADD CONSTRAINT film_category_category_id_fkey FOREIGN KEY (category_id) REFERENCES category(category_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: film_category_film_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY film_category
    ADD CONSTRAINT film_category_film_id_fkey FOREIGN KEY (film_id) REFERENCES film(film_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: film_language_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY film
    ADD CONSTRAINT film_language_id_fkey FOREIGN KEY (language_id) REFERENCES language(language_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: fk_address_city; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY address
    ADD CONSTRAINT fk_address_city FOREIGN KEY (city_id) REFERENCES city(city_id);


--
-- Name: fk_city; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY city
    ADD CONSTRAINT fk_city FOREIGN KEY (country_id) REFERENCES country(country_id);


--
-- Name: inventory_film_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY inventory
    ADD CONSTRAINT inventory_film_id_fkey FOREIGN KEY (film_id) REFERENCES film(film_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: payment_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY payment
    ADD CONSTRAINT payment_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer(customer_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: payment_rental_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY payment
    ADD CONSTRAINT payment_rental_id_fkey FOREIGN KEY (rental_id) REFERENCES rental(rental_id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: payment_staff_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY payment
    ADD CONSTRAINT payment_staff_id_fkey FOREIGN KEY (staff_id) REFERENCES staff(staff_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: rental_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY rental
    ADD CONSTRAINT rental_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer(customer_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: rental_inventory_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY rental
    ADD CONSTRAINT rental_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES inventory(inventory_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: rental_staff_id_key; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY rental
    ADD CONSTRAINT rental_staff_id_key FOREIGN KEY (staff_id) REFERENCES staff(staff_id);


--
-- Name: staff_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY staff
    ADD CONSTRAINT staff_address_id_fkey FOREIGN KEY (address_id) REFERENCES address(address_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: store_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY store
    ADD CONSTRAINT store_address_id_fkey FOREIGN KEY (address_id) REFERENCES address(address_id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: store_manager_staff_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY store
    ADD CONSTRAINT store_manager_staff_id_fkey FOREIGN KEY (manager_staff_id) REFERENCES staff(staff_id) ON UPDATE CASCADE ON DELETE RESTRICT;

--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;

DROP TABLE IF EXISTS supported_datatypes;
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
  col_smallint_array                          smallint[],
  col_int_array                               int[],
  col_bigint_array                            bigint[],
  col_character_array                         char(10)[],
  col_char_varying_array                      varchar(10)[],
  col_bit_array                               bit(10)[],
  col_varbit_array                            varbit(4)[],
  col_numeric_array                           numeric[],
  col_numeric_range_array                     numeric(5,3)[],
  col_double_precsion_array                   float4[],
  col_real_array                              float8[],
  col_money_array                             money[],
  col_time_array                              time without time zone[],
  col_intreval_array                          interval[],
  col_date_array                              date[],
  col_time_tz_array                           time with time zone[],
  col_timestamp_array                         timestamp with time zone[],
  col_timestamp_tz_array                      timestamp without time zone[],
  col_text_array                              text[],
  col_bool_array                              bool[],
  col_inet_array                              inet[],
  col_macaddr_array                           macaddr[],
  col_cidr_array                              cidr[],
  col_uuid_array                              uuid[],
  col_txid_snapshot_array                     txid_snapshot[],
  col_pg_lsn_array                            pg_lsn[],
  col_tsquery_array                           tsquery[],
  col_tsvector_array                          tsvector[],
  col_box_array                               box[],
  col_circle_array                            circle[],
  col_line_array                              line[],
  col_lseg_array                              lseg[],
  col_path_array                              path[],
  col_polygon_array                           polygon[],
  col_point_array                             point[],
  col_json_array                              json[],
  col_jsonb_array                             jsonb[],
  col_xml_array                               xml[]

);

DROP TABLE IF EXISTS big_table;
CREATE TABLE big_table (
     sample_col_1 smallint NOT NULL,
     sample_col_2 integer NOT NULL,
     sample_col_3 integer NOT NULL,
     sample_col_4 date NOT NULL,
     sample_col_5 date NOT NULL,
     sample_col_6 character varying(1) NOT NULL,
     sample_col_7 timestamp without time zone NOT NULL,
     sample_col_8 date,
     sample_col_9 smallint,
     sample_col_10 smallint,
     sample_col_11 date,
     sample_col_12 numeric(13,2),
     sample_col_13 numeric(13,2),
     sample_col_14 smallint,
     sample_col_15 date,
     sample_col_16 date,
     sample_col_17 smallint,
     sample_col_18 timestamp without time zone NOT NULL,
     sample_col_19 smallint,
     sample_col_20 smallint,
     sample_col_21 numeric(13,2),
     sample_col_22 numeric(5,2),
     sample_col_23 character varying(1),
     sample_col_24  date,
     sample_col_25 time without time zone NOT NULL,
     sample_col_26 character varying(1),
     sample_col_27 numeric(11,2),
     sample_col_28 character varying(1),
     sample_col_29 character varying(4) NOT NULL,
     sample_col_30 date NOT NULL,
     sample_col_31 time without time zone NOT NULL,
     sample_col_32 character varying(8) NOT NULL,
     sample_col_33 character varying(1) NOT NULL,
     sample_col_34 numeric(6,4) NOT NULL,
     sample_col_35 smallint NOT NULL,
     sample_col_36 smallint,
     sample_col_37 character varying(1) NOT NULL,
     sample_col_38 integer NOT NULL,
     sample_col_39 integer NOT NULL,
     sample_col_40 integer NOT NULL,
     sample_col_41 integer NOT NULL,
     sample_col_42 numeric(11,2),
     sample_col_43 character varying(1),
     sample_col_44 integer NOT NULL,
     sample_col_45 integer NOT NULL,
     sample_col_46 integer NOT NULL,
     sample_col_47 character varying(30) NOT NULL,
     sample_col_48 character varying(10) NOT NULL,
     sample_col_49 character varying(8) NOT NULL,
     sample_col_50 character varying(1) NOT NULL,
     sample_col_51 character varying(20) NOT NULL,
     sample_col_52 character varying(45) NOT NULL,
     sample_col_53 character varying(1),
     sample_col_54 character varying(1),
     sample_col_55 integer,
     sample_col_56 character varying(20) NOT NULL,
     sample_col_57 character varying(30) NOT NULL,
     sample_col_58 character varying(10) NOT NULL,
     sample_col_59 character varying(8) NOT NULL,
     sample_col_60 character varying(30) NOT NULL,
     sample_col_61 character varying(10) NOT NULL,
     sample_col_62 character varying(8) NOT NULL,
     sample_col_63 character varying(20) NOT NULL,
     sample_col_64 character varying(10),
     sample_col_65 time without time zone NOT NULL,
     sample_col_66 character varying(3),
     sample_col_67 character varying(1),
     sample_col_68 character varying(1) NOT NULL,
     sample_col_69 integer,
     sample_col_70 numeric(13,2) NOT NULL,
     sample_col_71 character varying(1),
     sample_col_72 character varying(1),
     sample_col_73 smallint,
     sample_col_74 smallint,
     sample_col_75 character varying(1) NOT NULL,
     sample_col_76 integer NOT NULL,
     sample_col_77 integer NOT NULL,
     sample_col_78 character varying(2) NOT NULL,
     sample_col_79 character varying(1),
     sample_col_80 character varying(1) NOT NULL,
     sample_col_81 numeric(11,2) NOT NULL,
     sample_col_82 character varying(1),
     sample_col_83 numeric(5,2) NOT NULL,
     sample_col_84 character varying(6) NOT NULL,
     sample_col_85 character varying(15) NOT NULL,
     sample_col_86 character varying(3) NOT NULL,
     sample_col_87 character varying(1) NOT NULL,
     sample_col_88 smallint NOT NULL,
     sample_col_89 smallint NOT NULL,
     sample_col_90 character varying(1) NOT NULL,
     sample_col_91 character varying(1) NOT NULL,
     sample_col_92 character varying(1),
     sample_col_93 integer,
     sample_col_94 character varying(6) NOT NULL,
     sample_col_95 character varying(1),
     sample_col_96 character varying(1) NOT NULL,
     sample_col_97 character varying(25) NOT NULL,
     sample_col_98 character varying(25) NOT NULL,
     sample_col_99 numeric(11,2),
     sample_col_100 character varying(1),
     sample_col_101 character varying(1),
     sample_col_102 smallint,
     sample_col_103 smallint,
     sample_col_104 smallint,
     sample_col_105 date NOT NULL,
     sample_col_106 date NOT NULL,
     sample_col_107 character varying(1),
     sample_col_108 character varying(1) NOT NULL,
     sample_col_109 character varying(1) NOT NULL,
     sample_col_110 smallint,
     sample_col_111 date,
     sample_col_112 character varying(1) NOT NULL,
     sample_col_113 bigint,
     sample_col_114 character varying(1),
     sample_col_115 character varying(1) NOT NULL,
     sample_col_116 character varying(1) NOT NULL,
     sample_col_117 smallint NOT NULL,
     sample_col_118 smallint NOT NULL,
     sample_col_119 smallint NOT NULL,
     sample_col_120 character varying(1) NOT NULL,
     sample_col_121 character varying(1),
     sample_col_122 integer,
     sample_col_123 integer,
     sample_col_124 smallint NOT NULL,
     sample_col_125 character varying(1) NOT NULL,
     sample_col_126 numeric(5,2) NOT NULL,
     sample_col_127 smallint,
     sample_col_128 character varying(1) NOT NULL,
     sample_col_129 character varying(25) NOT NULL,
     sample_col_130 character varying(1),
     sample_col_131 character varying(11) NOT NULL,
     sample_col_132 character varying(1) NOT NULL,
     sample_col_133 character varying(1),
     sample_col_134 character varying(1),
     sample_col_135 character varying(1) NOT NULL,
     sample_col_136 numeric(10,2) NOT NULL,
     sample_col_137 numeric(10,2) NOT NULL,
     sample_col_138 numeric(13,2),
     sample_col_139 numeric(10,2) NOT NULL,
     sample_col_140 character varying(1) NOT NULL,
     sample_col_141 character varying(4),
     sample_col_142 character varying(1),
     sample_col_143 character varying(1) NOT NULL,
     sample_col_144 smallint,
     sample_col_145 character varying(1),
     sample_col_146 smallint NOT NULL,
     sample_col_147 character varying(1) NOT NULL,
     sample_col_148 character varying(1) NOT NULL,
     sample_col_149 character varying(1),
     sample_col_150 character varying(5) NOT NULL,
     sample_col_151 character varying(5) NOT NULL,
     sample_col_152 character varying(1) NOT NULL,
     sample_col_153 character varying(1) NOT NULL,
     sample_col_154 character varying(1),
     sample_col_155 character varying(1),
     sample_col_156 smallint NOT NULL,
     sample_col_157 date,
     sample_col_158 character varying(2) NOT NULL,
     sample_col_159 character varying(1) NOT NULL,
     sample_col_160 character varying(1),
     sample_col_161 numeric(13,2),
     sample_col_162 smallint NOT NULL,
     sample_col_163 character varying(1),
     sample_col_164 character varying(2) NOT NULL,
     sample_col_165 smallint NOT NULL,
     sample_col_166 character varying(1) NOT NULL,
     sample_col_167 character varying(1) NOT NULL,
     sample_col_168 character varying(4),
     sample_col_169 character varying(2),
     sample_col_170 character varying(1) NOT NULL,
     sample_col_171 character varying(1) NOT NULL,
     sample_col_172 numeric(6,4) NOT NULL,
     sample_col_173 smallint NOT NULL,
     sample_col_174 character varying(8),
     sample_col_175 character varying(7),
     sample_col_176 smallint NOT NULL,
     sample_col_177 character varying(15) NOT NULL,
     sample_col_178 character varying(1) NOT NULL,
     sample_col_179 numeric(5,2) NOT NULL,
     sample_col_180 character varying(4) NOT NULL,
     sample_col_181 smallint NOT NULL,
     sample_col_182 character varying(6) NOT NULL,
     sample_col_183 integer NOT NULL,
     sample_col_184 character varying(1),
     sample_col_185 bigint,
     sample_col_186 numeric(11,2),
     sample_col_187 character varying(1) NOT NULL,
     sample_col_188 character varying(1) NOT NULL,
     sample_col_189 character varying(1) NOT NULL,
     sample_col_190 character varying(12) NOT NULL,
     sample_col_191 character varying(2) NOT NULL,
     sample_col_192 character varying(1) NOT NULL,
     sample_col_193 character varying(1),
     sample_col_194 smallint,
     sample_col_195 character varying(1),
     sample_col_196 smallint NOT NULL,
     sample_col_197 character varying(1),
     sample_col_198 character varying(36),
     sample_col_199 character varying(3),
     sample_col_200 bigint,
     sample_col_201 bigint NOT NULL,
     sample_col_202 bigint NOT NULL,
     sample_col_203 bigint NOT NULL,
     sample_col_204 character varying(1) NOT NULL,
     sample_col_205 character varying(1) NOT NULL,
     sample_col_206 character varying(20),
     sample_col_207 character varying(1) NOT NULL,
     sample_col_208 smallint,
     sample_col_209 character varying(1) NOT NULL,
     sample_col_210 date,
     sample_col_211 smallint,
     sample_col_212 smallint,
     sample_col_213 numeric(13,2),
     sample_col_214 smallint,
     sample_col_215 numeric(13,2),
     sample_col_216 smallint,
     sample_col_217 smallint,
     sample_col_218 character varying(1),
     sample_col_219 numeric(5,2),
     sample_col_220 smallint,
     sample_col_221 numeric(13,2),
     sample_col_222 smallint,
     sample_col_223 smallint,
     sample_col_224 date,
     sample_col_225 smallint NOT NULL,
     sample_col_226 smallint,
     sample_col_227 smallint,
     sample_col_228 integer,
     sample_col_229 date,
     sample_col_230 integer,
     sample_col_231 date,
     sample_col_232 integer,
     sample_col_233 date,
     sample_col_234 date,
     sample_col_235 timestamp without time zone,
     sample_col_236 date,
     sample_col_237 timestamp without time zone,
     sample_col_238 integer,
     sample_col_239 date,
     sample_col_240 integer,
     sample_col_241 date,
     sample_col_242 integer,
     sample_col_243 date,
     sample_col_244 integer,
     sample_col_245 date,
     sample_col_246 integer,
     sample_col_247 date,
     sample_col_248 integer,
     sample_col_249 date,
     sample_col_250 integer,
     sample_col_251 date,
     sample_col_252 integer,
     sample_col_253 date,
     sample_col_254 integer,
     sample_col_255 date
);
--
-- PostgreSQL database dump complete
--

`
}
