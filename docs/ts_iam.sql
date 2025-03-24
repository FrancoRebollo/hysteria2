SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

CREATE TABLE hysteria.acceso_api (
    api_key character varying(60) NOT NULL,
    uuid_api uuid NOT NULL,
    fecha_last_update timestamp without time zone DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) NOT NULL
);


ALTER TABLE hysteria.acceso_api OWNER TO postgres;

CREATE TABLE hysteria.anuncios (
    id integer NOT NULL,
    texto text NOT NULL,
    fecha date NOT NULL
);


ALTER TABLE hysteria.anuncios OWNER TO postgres;

CREATE SEQUENCE hysteria.anuncios_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hysteria.anuncios_id_seq OWNER TO postgres;

ALTER SEQUENCE hysteria.anuncios_id_seq OWNED BY hysteria.anuncios.id;


CREATE TABLE hysteria.api (
    uuid_api uuid NOT NULL,
    api character varying(60) NOT NULL,
    version character varying(12) NOT NULL,
    fecha_last_update timestamp without time zone DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) DEFAULT USER NOT NULL
);


ALTER TABLE hysteria.api OWNER TO postgres;


CREATE TABLE hysteria.api_key (
    api_key character varying(60) DEFAULT gen_random_uuid() NOT NULL,
    app_origen character varying(60) NOT NULL,
    estado character varying(15) DEFAULT 'ACTIVO'::character varying NOT NULL,
    req_2fa character(1) DEFAULT 'N'::bpchar NOT NULL,
    ctd_hs_access_token_valido integer DEFAULT 1 NOT NULL,
    req_usuario_db character(1) DEFAULT 'S'::bpchar NOT NULL,
    fecha_vigencia date DEFAULT CURRENT_DATE NOT NULL,
    fecha_fin_vigencia date,
    ctrl_limite_acceso_tiempo character(1) DEFAULT 'N'::bpchar,
    ctd_accesos_unidad_tiempo integer,
    unidad_tiempo_acceso character varying(15) DEFAULT 'MINUTO'::character varying,
    fecha_last_update timestamp without time zone DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) DEFAULT USER NOT NULL,
    is_super_user character(1) DEFAULT 'N'::bpchar NOT NULL
);


ALTER TABLE hysteria.api_key OWNER TO postgres;

CREATE TABLE hysteria.bosses (
    id_bosses integer NOT NULL,
    nombre text NOT NULL,
    respawn_time integer,
    interval_respawn_time integer,
    unidad_interval_respawn_time text,
    lunes text,
    martes text,
    miercoles text,
    jueves text,
    viernes text,
    sabado text,
    domingo text,
    last_time_of_death time without time zone,
    respawn_fijo boolean,
    CONSTRAINT bosses_unidad_interval_respawn_time_check CHECK ((unidad_interval_respawn_time = ANY (ARRAY['minutes'::text, 'hours'::text, 'days'::text])))
);


ALTER TABLE hysteria.bosses OWNER TO postgres;

CREATE SEQUENCE hysteria.bosses_id_bosses_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hysteria.bosses_id_bosses_seq OWNER TO postgres;
ALTER SEQUENCE hysteria.bosses_id_bosses_seq OWNED BY hysteria.bosses.id_bosses;


CREATE TABLE hysteria.canal_digital_persona (
    id_canal_digital_persona integer DEFAULT nextval('hysteria.canal_digital_persona_seq'::regclass) NOT NULL,
    id_persona integer NOT NULL,
    tipo_canal_digital character varying(25) NOT NULL,
    password_acceso_hash character varying(256) NOT NULL,
    id_mail_persona integer,
    id_te_persona integer,
    login_name character varying(100),
    canal_validado character(1) DEFAULT 'N'::bpchar,
    fecha_validacion_canal date,
    acceso_revocado character(1) DEFAULT 'N'::bpchar NOT NULL,
    req_2fa character(1) DEFAULT 'N'::bpchar NOT NULL,
    fecha_last_update date DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) DEFAULT 'CURRENT_USER'::character varying NOT NULL
);


ALTER TABLE hysteria.canal_digital_persona OWNER TO postgres;

CREATE TABLE hysteria.error_log (
    id_error_log integer DEFAULT nextval('hysteria.id_error_log'::regclass) NOT NULL,
    message_error character varying(5000) NOT NULL,
    endpoint character varying(400),
    id_tipo_error integer DEFAULT 0,
    ip_address character varying(50),
    id_persona integer NOT NULL,
    canal_digital character varying(60) NOT NULL,
    api_key character varying(60) NOT NULL,
    id_token integer NOT NULL,
    access_token character varying(500) NOT NULL,
    fecha_last_update date DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) DEFAULT 'CURRENT_USER'::character varying
);


ALTER TABLE hysteria.error_log OWNER TO postgres;

CREATE TABLE hysteria.exc_acceso_endpoint_api (
    id_exc_acceso_endpoint_api integer NOT NULL,
    api_key character varying(60) NOT NULL,
    uuid_api uuid NOT NULL,
    metodo character varying(10) NOT NULL,
    recurso character varying(120) NOT NULL,
    fecha_last_update date DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) DEFAULT 'CURRENT_USER'::character varying
);


ALTER TABLE hysteria.exc_acceso_endpoint_api OWNER TO postgres;

CREATE TABLE hysteria.hist_token (
    id_hist_token integer DEFAULT nextval('hysteria.hist_id_token'::regclass) NOT NULL,
    id_token integer NOT NULL,
    api_key character varying(60) NOT NULL,
    id_canal_digital_persona integer NOT NULL,
    access_token character varying(500) NOT NULL,
    fecha_creacion_token date,
    fecha_exp_access_token date,
    refresh_token character varying(500) NOT NULL,
    fecha_exp_refresh_token date,
    acceso_revocado character(1) DEFAULT 'N'::bpchar NOT NULL,
    fecha_last_update date DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) DEFAULT 'CURRENT_USER'::character varying
);


ALTER TABLE hysteria.hist_token OWNER TO postgres;

CREATE TABLE hysteria.location (
    id_location character varying(15) NOT NULL,
    fecha_last_update date DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) DEFAULT 'CURRENT_USER'::character varying NOT NULL
);


ALTER TABLE hysteria.location OWNER TO postgres;

CREATE TABLE hysteria.persona (
    id_persona integer NOT NULL,
    last_location character varying(15) NOT NULL,
    acceso_revocado character(1) DEFAULT 'N'::bpchar NOT NULL,
    fecha_last_update date DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) DEFAULT 'CURRENT_USER'::character varying
);


ALTER TABLE hysteria.persona OWNER TO postgres;

CREATE SEQUENCE hysteria.persona_id_persona_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE hysteria.persona_id_persona_seq OWNER TO postgres;

ALTER SEQUENCE hysteria.persona_id_persona_seq OWNED BY hysteria.persona.id_persona;

CREATE TABLE hysteria.tipo_canal_digital_df (
    tipo_canal_digital character varying(25) NOT NULL,
    acceso_revocado character(1) DEFAULT 'N'::bpchar NOT NULL,
    fecha_last_update date DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) DEFAULT 'CURRENT_USER'::character varying NOT NULL
);


ALTER TABLE hysteria.tipo_canal_digital_df OWNER TO postgres;

CREATE TABLE hysteria.token (
    id_token integer DEFAULT nextval('hysteria.id_token'::regclass) NOT NULL,
    api_key character varying(60) NOT NULL,
    id_canal_digital_persona integer NOT NULL,
    access_token character varying(500) NOT NULL,
    fecha_creacion_token date,
    fecha_exp_access_token timestamp without time zone,
    refresh_token character varying(500) NOT NULL,
    fecha_exp_refresh_token timestamp without time zone,
    acceso_revocado character(1) DEFAULT 'N'::bpchar NOT NULL,
    fecha_last_update date DEFAULT CURRENT_DATE NOT NULL,
    actualizado_por character varying(30) DEFAULT 'CURRENT_USER'::character varying,
    "2fa_seed" character varying(100),
    last_code_2fa numeric
);


ALTER TABLE hysteria.token OWNER TO postgres;

ALTER TABLE ONLY hysteria.anuncios ALTER COLUMN id SET DEFAULT nextval('hysteria.anuncios_id_seq'::regclass);

ALTER TABLE ONLY hysteria.bosses ALTER COLUMN id_bosses SET DEFAULT nextval('hysteria.bosses_id_bosses_seq'::regclass);

ALTER TABLE ONLY hysteria.persona ALTER COLUMN id_persona SET DEFAULT nextval('hysteria.persona_id_persona_seq'::regclass);

ALTER TABLE ONLY hysteria.anuncios
    ADD CONSTRAINT anuncios_pkey PRIMARY KEY (id);

ALTER TABLE ONLY hysteria.bosses
    ADD CONSTRAINT bosses_pkey PRIMARY KEY (id_bosses);

ALTER TABLE ONLY hysteria.persona
    ADD CONSTRAINT persona_pkey PRIMARY KEY (id_persona);

ALTER TABLE ONLY hysteria.api
    ADD CONSTRAINT pk_api PRIMARY KEY (uuid_api);

ALTER TABLE ONLY hysteria.api_key
    ADD CONSTRAINT pk_api_key PRIMARY KEY (api_key);

ALTER TABLE ONLY hysteria.acceso_api
    ADD CONSTRAINT pk_api_key_api PRIMARY KEY (uuid_api, api_key);

ALTER TABLE ONLY hysteria.tipo_canal_digital_df
    ADD CONSTRAINT pk_canal_digital_df PRIMARY KEY (tipo_canal_digital);

ALTER TABLE ONLY hysteria.canal_digital_persona
    ADD CONSTRAINT pk_canal_digital_persona PRIMARY KEY (id_canal_digital_persona);

ALTER TABLE ONLY hysteria.error_log
    ADD CONSTRAINT pk_error_log PRIMARY KEY (id_error_log);

ALTER TABLE ONLY hysteria.exc_acceso_endpoint_api
    ADD CONSTRAINT pk_exc_endpoint_api_key PRIMARY KEY (id_exc_acceso_endpoint_api);
	
ALTER TABLE ONLY hysteria.hist_token
    ADD CONSTRAINT pk_hist_token PRIMARY KEY (id_hist_token);

ALTER TABLE ONLY hysteria.location
    ADD CONSTRAINT pk_locations PRIMARY KEY (id_location);

ALTER TABLE ONLY hysteria.token
    ADD CONSTRAINT pk_token_api_key PRIMARY KEY (id_token);

ALTER TABLE ONLY hysteria.bosses
    ADD CONSTRAINT unica UNIQUE (nombre);

ALTER TABLE ONLY hysteria.canal_digital_persona
    ADD CONSTRAINT unique_login_name UNIQUE (login_name);

ALTER TABLE ONLY hysteria.acceso_api
    ADD CONSTRAINT unq_acceso_api_api_key UNIQUE (api_key, uuid_api);

ALTER TABLE ONLY hysteria.acceso_api
    ADD CONSTRAINT fk_api_key_api_api FOREIGN KEY (uuid_api) REFERENCES hysteria.api(uuid_api);
	
ALTER TABLE ONLY hysteria.acceso_api
    ADD CONSTRAINT fk_api_key_api_api_key FOREIGN KEY (api_key) REFERENCES hysteria.api_key(api_key);
	
ALTER TABLE ONLY hysteria.canal_digital_persona
    ADD CONSTRAINT fk_canal_digital_persona_canal_digital_df FOREIGN KEY (tipo_canal_digital) REFERENCES hysteria.tipo_canal_digital_df(tipo_canal_digital);

ALTER TABLE ONLY hysteria.canal_digital_persona
    ADD CONSTRAINT fk_canal_digital_persona_persona FOREIGN KEY (id_persona) REFERENCES hysteria.persona(id_persona);
	

ALTER TABLE ONLY hysteria.exc_acceso_endpoint_api
    ADD CONSTRAINT fk_exc_acceso_endpoint_api FOREIGN KEY (api_key, uuid_api) REFERENCES hysteria.acceso_api(api_key, uuid_api);
	

ALTER TABLE ONLY hysteria.persona
    ADD CONSTRAINT fk_persona_locations FOREIGN KEY (last_location) REFERENCES hysteria.location(id_location);

ALTER TABLE ONLY hysteria.token
    ADD CONSTRAINT fk_token_api_key FOREIGN KEY (api_key) REFERENCES hysteria.api_key(api_key);
	
ALTER TABLE ONLY hysteria.token
    ADD CONSTRAINT fk_token_canal_digital_persona FOREIGN KEY (id_canal_digital_persona) REFERENCES hysteria.canal_digital_persona(id_canal_digital_persona);


-- Completed on 2025-03-20 03:45:54

--
-- PostgreSQL database dump complete
--