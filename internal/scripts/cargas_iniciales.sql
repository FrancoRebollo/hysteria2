INSERT INTO hysteria.tipo_canal_digital_df 
( tipo_canal_digital )
VALUES
('APLICACION_MOBILE')
;

INSERT INTO hysteria.tipo_canal_digital_df 
( tipo_canal_digital )
VALUES
('USER_PASSWORD')
;

INSERT INTO hysteria.location 
( id_location )
VALUES
('0')
;

DELETE FROM hysteria.API;

INSERT INTO hysteria.api (uuid_api, api, "version", fecha_last_update, actualizado_por)
VALUES 
    (gen_random_uuid(), 'TURNOS', '1.0', CURRENT_TIMESTAMP, USER),
    (gen_random_uuid(), 'TRAZABILIDAD', '1.0', CURRENT_TIMESTAMP, USER);

insert into hysteria.api_key (api_key,app_origen,fecha_fin_vigencia,ctd_accesos_unidad_tiempo)
values 
('x#16!xF12QsfjyQ2351QSVXZ','HIS',CURRENT_DATE+30,0),
('x#16!xF12QsfjyQ2351ADSGA','HIS',CURRENT_DATE+30,0);
COMMIT;