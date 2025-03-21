CREATE SCHEMA IF NOT EXISTS hysteria;

CREATE TYPE / DOMAIN d_no_si AS ...
check (null);

CREATE TABLE hysteria.version_modelo (
	service_name		varchar(60)  NOT NULL ,
	version_modelo		varchar(60)  NOT NULL ,
	fecha_last_update    timestamp DEFAULT CURRENT_DATE NOT NULL
);


CREATE  TABLE hysteria.api ( 
	uuid_api             uuid  NOT NULL  ,
	api                  varchar(60)  NOT NULL  ,
	"version"            varchar(12)  NOT NULL  ,
	fecha_last_update    timestamp DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por      varchar(30) DEFAULT USER NOT NULL  ,
	CONSTRAINT pk_api PRIMARY KEY ( uuid_api )
 );


CREATE  TABLE hysteria.api_key ( 
	api_key              varchar(60) DEFAULT gen_random_uuid () NOT NULL  ,
	app_origen           varchar(60)  NOT NULL  ,
	estado               varchar(15) DEFAULT 'ACTIVO' NOT NULL  ,
	req_2fa              char(1) DEFAULT 'N' NOT NULL  ,
	ctd_hs_access_token_valido integer DEFAULT 1 NOT NULL  ,
	req_usuario_db       char(1) DEFAULT 'S' NOT NULL  ,
	fecha_vigencia       date DEFAULT CURRENT_DATE NOT NULL  ,
	fecha_fin_vigencia   date    ,
	ctrl_limite_acceso_tiempo char(1) DEFAULT 'N' NULL,
	ctd_accesos_unidad_tiempo integer    ,
	unidad_tiempo_acceso varchar(15) DEFAULT 'MINUTO'   ,
	fecha_last_update    timestamp DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por      varchar(30) DEFAULT USER NOT NULL  ,
	CONSTRAINT pk_api_key PRIMARY KEY ( api_key )
 );

ALTER TABLE hysteria.api_key
ADD COLUMN is_super_user char(1) DEFAULT 'N' NOT NULL;

CREATE  TABLE hysteria.tipo_canal_digital_df ( 
	tipo_canal_digital   varchar(25)  NOT NULL  ,
	acceso_revocado      char(1) DEFAULT 'N' NOT NULL  ,
	fecha_last_update    date DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por      varchar(30) DEFAULT 'CURRENT_USER' NOT NULL  ,
	CONSTRAINT pk_canal_digital_df PRIMARY KEY ( tipo_canal_digital )
 );

CREATE  TABLE hysteria.location ( 
	id_location          varchar(15)  NOT NULL  ,
	fecha_last_update    date DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por      varchar(30) DEFAULT 'CURRENT_USER' NOT NULL  ,
	CONSTRAINT pk_locations PRIMARY KEY ( id_location )
 );

CREATE  TABLE hysteria.acceso_api ( 
	api_key              varchar(60)  NOT NULL  ,
	uuid_api             uuid  NOT NULL  ,
	fecha_last_update    timestamp DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por      varchar(30)  NOT NULL  ,
	CONSTRAINT pk_api_key_api PRIMARY KEY ( uuid_api, api_key ),
	CONSTRAINT unq_acceso_api_api_key UNIQUE ( api_key, uuid_api ) 
 );

CREATE  TABLE hysteria.exc_acceso_endpoint_api ( 
	id_exc_acceso_endpoint_api integer  NOT NULL  ,
	api_key              varchar(60)  NOT NULL  ,
	uuid_api             uuid  NOT NULL  ,
	metodo               varchar(10)  NOT NULL  ,
	recurso              varchar(120)  NOT NULL  ,
	fecha_last_update    date DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por      varchar(30) DEFAULT 'CURRENT_USER'   ,
	CONSTRAINT pk_exc_endpoint_api_key PRIMARY KEY ( id_exc_acceso_endpoint_api )
 );

CREATE  TABLE hysteria.persona ( 
	id_persona            serial primary key  ,
	last_location         varchar(15)  NOT NULL  ,
	acceso_revocado       char(1) DEFAULT 'N' NOT NULL  ,
	fecha_last_update     date DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por       varchar(30) DEFAULT 'CURRENT_USER'   ,
	--CONSTRAINT pk_persona PRIMARY KEY ( id_persona )
 );

CREATE  TABLE hysteria.canal_digital_persona ( 
	id_canal_digital_persona integer  NOT NULL  ,
	id_persona           integer  NOT NULL  ,
	tipo_canal_digital   varchar(25)  NOT NULL  ,
	password_acceso_hash varchar(256)  NOT NULL  ,
	id_mail_persona		 integer    ,
	id_te_persona 		 integer    ,
	login_name           varchar(100),
	canal_validado       char(1) DEFAULT 'N'   ,
	fecha_validacion_canal date    ,
	acceso_revocado      char(1) DEFAULT 'N' NOT NULL  ,
	req_2fa				 char(1) DEFAULT 'N' NOT NULL  ,
	fecha_last_update    date DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por      varchar(30) DEFAULT 'CURRENT_USER' NOT NULL  ,
	CONSTRAINT pk_canal_digital_persona PRIMARY KEY ( id_canal_digital_persona )
 );


CREATE SEQUENCE hysteria.canal_digital_persona_seq;

ALTER TABLE hysteria.canal_digital_persona 
    ALTER COLUMN id_canal_digital_persona SET DEFAULT nextval('hysteria.canal_digital_persona_seq');

ALTER TABLE hysteria.canal_digital_persona
ADD CONSTRAINT unique_login_name UNIQUE (login_name);

CREATE  TABLE hysteria.token ( 
	id_token             integer  NOT NULL  ,
	api_key              varchar(60)  NOT NULL  ,
	id_canal_digital_persona integer  NOT NULL  ,
	access_token         varchar(128)  NOT NULL  ,
	fecha_creacion_token date    ,
	fecha_exp_access_token date    ,
	refresh_token        varchar(128)  NOT NULL  ,
	fecha_exp_refresh_token date    ,
	acceso_revocado      char(1) DEFAULT 'N' NOT NULL  ,
	fecha_last_update    date DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por      varchar(30) DEFAULT 'CURRENT_USER'   ,
	2fa_seed			varchar(100) ,
	CONSTRAINT pk_token_api_key PRIMARY KEY ( id_token )
 );

ALTER TABLE hysteria.token
ALTER COLUMN fecha_exp_Refresh_token TYPE TIMESTAMP;


ALTER TABLE hysteria.token
ALTER COLUMN fecha_exp_access_token TYPE TIMESTAMP;

ALTER TABLE hysteria.token
ADD COLUMN last_code_2fa NUMERIC;

CREATE SEQUENCE hysteria.id_token;

ALTER TABLE hysteria.token 
    ALTER COLUMN id_token SET DEFAULT nextval('hysteria.id_token');

ALTER TABLE hysteria.token
ALTER COLUMN refresh_token TYPE varchar(500);

ALTER TABLE hysteria.token
ALTER COLUMN access_token TYPE varchar(500);

CREATE  TABLE hysteria.hist_token ( 
	id_hist_token		 integer not null,
	id_token             integer  NOT NULL  ,
	api_key              varchar(60)  NOT NULL  ,
	id_canal_digital_persona integer  NOT NULL  ,
	access_token         varchar(128)  NOT NULL  ,
	fecha_creacion_token date    ,
	fecha_exp_access_token date    ,
	refresh_token        varchar(128)  NOT NULL  ,
	fecha_exp_refresh_token date    ,
	acceso_revocado      char(1) DEFAULT 'N' NOT NULL  ,
	fecha_last_update    date DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por      varchar(30) DEFAULT 'CURRENT_USER'   ,
	CONSTRAINT pk_hist_token PRIMARY KEY ( id_hist_token )
 );

CREATE SEQUENCE hysteria.hist_id_token;

ALTER TABLE hysteria.hist_token 
    ALTER COLUMN id_hist_token SET DEFAULT nextval('hysteria.hist_id_token');

ALTER TABLE hysteria.hist_token
ALTER COLUMN refresh_token TYPE varchar(500);

ALTER TABLE hysteria.hist_token
ALTER COLUMN access_token TYPE varchar(500);

CREATE  TABLE hysteria.error_log ( 
	id_error_log		 	integer not null,
	message_error			varchar(5000) not null ,
	endpoint				varchar(400),
	id_TIPO_ERROR			integer default 0 ,
	ip_address				varchar(50),
	ID_PERSONA             	integer  NOT NULL  ,
	canal_digital           varchar(60)  NOT NULL ,
	api_key              	varchar(60)  NOT NULL ,
	id_token             	integer  NOT NULL  ,
	access_token			varchar(500) not null ,
	fecha_last_update    	date DEFAULT CURRENT_DATE NOT NULL  ,
	actualizado_por      	varchar(30) DEFAULT 'CURRENT_USER'   ,
	CONSTRAINT pk_error_log PRIMARY KEY ( id_error_log )
 );

CREATE SEQUENCE hysteria.id_error_log;

ALTER TABLE hysteria.error_log 
    ALTER COLUMN id_error_log SET DEFAULT nextval('hysteria.id_error_log');

ALTER TABLE hysteria.acceso_api ADD CONSTRAINT fk_api_key_api_api FOREIGN KEY ( uuid_api ) REFERENCES hysteria.api( uuid_api );

ALTER TABLE hysteria.acceso_api ADD CONSTRAINT fk_api_key_api_api_key FOREIGN KEY ( api_key ) REFERENCES hysteria.api_key( api_key );

ALTER TABLE hysteria.canal_digital_persona ADD CONSTRAINT fk_canal_digital_persona_persona FOREIGN KEY ( id_persona ) REFERENCES hysteria.persona( id_persona );

ALTER TABLE hysteria.canal_digital_persona ADD CONSTRAINT fk_canal_digital_persona_canal_digital_df FOREIGN KEY ( tipo_canal_digital ) REFERENCES hysteria.tipo_canal_digital_df( tipo_canal_digital );

ALTER TABLE hysteria.exc_acceso_endpoint_api ADD CONSTRAINT fk_exc_acceso_endpoint_api FOREIGN KEY ( api_key, uuid_api ) REFERENCES hysteria.acceso_api( api_key, uuid_api );

ALTER TABLE hysteria.persona ADD CONSTRAINT fk_persona_locations FOREIGN KEY ( last_location ) REFERENCES hysteria.location( id_location );

ALTER TABLE hysteria.token ADD CONSTRAINT fk_token_api_key FOREIGN KEY ( api_key ) REFERENCES hysteria.api_key( api_key );

ALTER TABLE hysteria.token ADD CONSTRAINT fk_token_canal_digital_persona FOREIGN KEY ( id_canal_digital_persona ) REFERENCES hysteria.canal_digital_persona( id_canal_digital_persona );

