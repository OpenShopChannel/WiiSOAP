--
-- PostgreSQL database dump
--

-- Dumped from database version 13.2
-- Dumped by pg_dump version 13.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: owned_titles; Type: TABLE; Schema: public; Owner: wiisoap
--

CREATE TABLE public.owned_titles (
                                     account_id integer NOT NULL,
                                     ticket_id character varying(16) NOT NULL,
                                     title_id character varying(16) NOT NULL,
                                     revocation_date timestamp without time zone
);


ALTER TABLE public.owned_titles OWNER TO wiisoap;

--
-- Name: shop_titles; Type: TABLE; Schema: public; Owner: wiisoap
--

CREATE TABLE public.shop_titles (
                                    title_id character varying(16) NOT NULL,
                                    version integer,
                                    description text
);


ALTER TABLE public.shop_titles OWNER TO wiisoap;

--
-- Name: COLUMN shop_titles.description; Type: COMMENT; Schema: public; Owner: wiisoap
--

COMMENT ON COLUMN public.shop_titles.description IS 'Description of the title.';


--
-- Name: userbase; Type: TABLE; Schema: public; Owner: wiisoap
--

CREATE TABLE public.userbase (
                                 device_id bigint NOT NULL,
                                 device_token character varying(21) NOT NULL,
                                 device_token_hashed character varying(32) NOT NULL,
                                 account_id integer NOT NULL,
                                 region character varying(3),
                                 country character varying(2),
                                 language character varying(2),
                                 serial_number character varying(11),
                                 device_code bigint
);


ALTER TABLE public.userbase OWNER TO wiisoap;

--
-- Name: COLUMN userbase.device_code; Type: COMMENT; Schema: public; Owner: wiisoap
--

COMMENT ON COLUMN public.userbase.device_code IS 'Also known as the console''s friend code.';


--
-- Name: owned_titles owned_titles_pk; Type: CONSTRAINT; Schema: public; Owner: wiisoap
--

ALTER TABLE ONLY public.owned_titles
    ADD CONSTRAINT owned_titles_pk PRIMARY KEY (account_id);


--
-- Name: shop_titles shop_titles_pk; Type: CONSTRAINT; Schema: public; Owner: wiisoap
--

ALTER TABLE ONLY public.shop_titles
    ADD CONSTRAINT shop_titles_pk PRIMARY KEY (title_id);


--
-- Name: userbase userbase_pk; Type: CONSTRAINT; Schema: public; Owner: wiisoap
--

ALTER TABLE ONLY public.userbase
    ADD CONSTRAINT userbase_pk PRIMARY KEY (account_id);


--
-- Name: owned_titles_account_id_uindex; Type: INDEX; Schema: public; Owner: wiisoap
--

CREATE UNIQUE INDEX owned_titles_account_id_uindex ON public.owned_titles USING btree (account_id);


--
-- Name: shop_titles_title_id_uindex; Type: INDEX; Schema: public; Owner: wiisoap
--

CREATE UNIQUE INDEX shop_titles_title_id_uindex ON public.shop_titles USING btree (title_id);


--
-- Name: userbase_account_id_uindex; Type: INDEX; Schema: public; Owner: wiisoap
--

CREATE UNIQUE INDEX userbase_account_id_uindex ON public.userbase USING btree (account_id);


--
-- Name: userbase_device_code_uindex; Type: INDEX; Schema: public; Owner: wiisoap
--

CREATE UNIQUE INDEX userbase_device_code_uindex ON public.userbase USING btree (device_code);


--
-- Name: userbase_device_token_uindex; Type: INDEX; Schema: public; Owner: wiisoap
--

CREATE UNIQUE INDEX userbase_device_token_uindex ON public.userbase USING btree (device_token);


--
-- Name: owned_titles match_shop_title_metadata; Type: FK CONSTRAINT; Schema: public; Owner: wiisoap
--

ALTER TABLE ONLY public.owned_titles
    ADD CONSTRAINT match_shop_title_metadata FOREIGN KEY (title_id) REFERENCES public.shop_titles(title_id);


--
-- Name: owned_titles order_account_ids; Type: FK CONSTRAINT; Schema: public; Owner: wiisoap
--

ALTER TABLE ONLY public.owned_titles
    ADD CONSTRAINT order_account_ids FOREIGN KEY (account_id) REFERENCES public.userbase(account_id);


--
-- PostgreSQL database dump complete
--