package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"plantilla_api/cmd/utils"
	"plantilla_api/internal/version/domains"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

func (v SecurityRepository) AltaUser(ctx context.Context, altaUser domains.RequestAltaUser) (string, error) {
	var message string

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `SELECT TIPO_CANAL_DIGITAL FROM hysteria.TIPO_CANAL_DIGITAL_DF WHERE TIPO_CANAL_DIGITAL = $1`

	rows, err := tx.QueryContext(ctx, query, altaUser.CanalDigital)

	if err != nil {
		return "", err
	}

	if !rows.Next() {
		return "", fmt.Errorf("canal digital invalido")
	}
	rows.Close()

	query = `SELECT id_persona FROM hysteria.persona WHERE id_persona = $1`

	rows, err = tx.QueryContext(ctx, query, fmt.Sprint(altaUser.IdPersona))

	if err != nil {
		return "", err
	}

	if !rows.Next() {
		insert := "INSERT INTO hysteria.PERSONA (ID_PERSONA,LAST_LOCATION) VALUES ($1,0)"

		_, err = tx.ExecContext(ctx, insert, fmt.Sprint(altaUser.IdPersona))

		if err != nil {
			return "", err
		}

		message = "Persona creada - "
	}

	rows.Close()

	query = `SELECT id_canal_digital_persona FROM hysteria.canal_digital_persona WHERE id_persona = $1 and tipo_canal_digital = $2`

	rows, err = tx.QueryContext(ctx, query, fmt.Sprint(altaUser.IdPersona), altaUser.CanalDigital)

	if err != nil {
		return "", err
	}

	if !rows.Next() {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(altaUser.Password), bcrypt.DefaultCost)

		insert := `INSERT INTO hysteria.CANAL_DIGITAL_PERSONA 
		(ID_PERSONA,TIPO_CANAL_DIGITAL,PASSWORD_ACCESO_HASH,id_mail_persona,id_te_persona,LOGIN_NAME) VALUES ($1,$2,$3,$4,$5,$6)`

		_, err = tx.ExecContext(ctx, insert, altaUser.IdPersona, altaUser.CanalDigital, hashedPassword, altaUser.IdMailPersona,
			altaUser.IdTePersona, altaUser.LoginName)

		if err != nil {
			return "", err
		}

		message += " en su canal digital"
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return message, nil
}

func (v SecurityRepository) checkCredentials(ctx context.Context, credentials domains.Credentials) error {

	query := `SELECT id_persona FROM hysteria.persona WHERE id_persona = $1`

	rows, err := v.dbPost.GetDB().QueryContext(ctx, query, credentials.IdPersona)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("persona no encontrada")
	}
	rows.Close()

	query = `SELECT tipo_canal_digital FROM hysteria.tipo_canal_digital_df WHERE tipo_canal_digital = $1`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.CanalDigital)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("canal digital invalido")
	}
	rows.Close()

	query = `SELECT id_canal_digital_persona FROM hysteria.canal_digital_persona WHERE tipo_canal_digital = $1
		and id_persona = $2 and canal_validado = 'N'`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.CanalDigital, credentials.IdPersona)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("canal digital no validado")
	}
	rows.Close()

	query = `SELECT api_key FROM hysteria.api_key WHERE api_key = $1`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.ApiKey)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("api key desconocida")
	}
	rows.Close()

	return nil
}

func (v SecurityRepository) checkRevokes(ctx context.Context, credentials domains.Credentials) error {

	query := `SELECT id_persona FROM hysteria.persona WHERE id_persona = $1 and acceso_revocado = 'S'`

	rows, err := v.dbPost.GetDB().QueryContext(ctx, query, credentials.IdPersona)

	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("persona revocada")
	}
	rows.Close()

	query = `SELECT tipo_canal_digital FROM hysteria.tipo_canal_digital_df WHERE tipo_canal_digital = $1 and acceso_revocado = 'S'`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.CanalDigital)

	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("canal digital revocado")
	}
	rows.Close()

	query = `SELECT id_canal_digital_persona FROM hysteria.canal_digital_persona WHERE tipo_canal_digital = $1
		and id_persona = $2 and canal_validado = 'N' and acceso_revocado = 'S'`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.CanalDigital, credentials.IdPersona)

	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("acceso revocado persona - canal digital")
	}
	rows.Close()

	query = `SELECT api_key FROM hysteria.api_key WHERE api_key = $1 and fecha_fin_vigencia < current_date`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.ApiKey)

	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("api key expirada")
	}
	rows.Close()

	return nil
}

func (v SecurityRepository) CheckTokenCreation(ctx context.Context, credentials domains.Credentials) error {

	if err := v.checkCredentials(ctx, credentials); err != nil {
		return err
	}

	if err := v.checkRevokes(ctx, credentials); err != nil {
		return err
	}

	query := `SELECT id_token,fecha_exp_refresh_token FROM hysteria.TOKEN
		WHERE ID_CANAL_DIGITAL_PERSONA = (SELECT ID_CANAL_DIGITAL_PERSONA FROM hysteria.CANAL_DIGITAL_PERSONA
											WHERE ID_PERSONA = $1 AND TIPO_CANAL_DIGITAL = $2)
		and api_key = $3`

	rows, err := v.dbPost.GetDB().QueryContext(ctx, query, credentials.IdPersona, credentials.CanalDigital, credentials.ApiKey)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("loguee por primera vez")
	}

	rows.Close()

	return nil
}

func (v SecurityRepository) updateAccessToken(ctx context.Context, credentialsToken domains.CredentialsToken, idCanalDigitalPersona int) error {

	var idToken int
	accesTokenDuration, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION"))

	if err != nil {
		return fmt.Errorf("error calculando duracion de refresh token")
	}

	accessExpiresAt := time.Now().Add(time.Minute * time.Duration(accesTokenDuration))

	query := `SELECT id_token FROM hysteria.TOKEN	WHERE ID_CANAL_DIGITAL_PERSONA = $1 AND API_KEY = $2`

	_ = v.dbPost.GetDB().QueryRowContext(ctx, query, idCanalDigitalPersona, credentialsToken.ApiKey).Scan(&idToken)

	if idToken == 0 {
		return fmt.Errorf("registro no encontrado en token - loguee por primera vez")
	}

	insert := `
    INSERT INTO hysteria.hist_token (id_token, api_key, id_canal_digital_persona, access_token, fecha_creacion_token, fecha_exp_access_token, refresh_token, fecha_exp_refresh_token, acceso_revocado)
    SELECT id_token, api_key, id_canal_digital_persona, access_token, fecha_creacion_token, fecha_exp_access_token, refresh_token, fecha_exp_refresh_token, acceso_revocado
    FROM hysteria.token
    WHERE id_token = $1`

	_, err = v.dbPost.GetDB().ExecContext(ctx, insert, idToken)

	if err != nil {
		return err
	}

	update := `update hysteria.token set access_token = $1, fecha_exp_access_token = $2 
		where id_token = $3`

	_, err = v.dbPost.GetDB().ExecContext(ctx, update, credentialsToken.AccessToken, accessExpiresAt, idToken)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) PersistToken(ctx context.Context, credentials domains.CredentialsToken) error {
	query := `SELECT ID_CANAL_DIGITAL_PERSONA FROM hysteria.CANAL_DIGITAL_PERSONA WHERE ID_PERSONA = $1 
		AND TIPO_CANAL_DIGITAL = $2`

	var idCanalDigitalPersona int

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, credentials.IdPersona, credentials.CanalDigital).Scan(&idCanalDigitalPersona)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no se encontró el canal digital en relacion a la persona")
		}
		return err
	}

	if err = v.updateAccessToken(ctx, credentials, idCanalDigitalPersona); err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) tokenExpiration(ctx context.Context, credentials domains.Credentials) error {
	var idToken int
	var fechaExpAccessToken time.Time
	query := `SELECT id_token,fecha_exp_access_token FROM hysteria.TOKEN
		WHERE ID_CANAL_DIGITAL_PERSONA = (SELECT ID_CANAL_DIGITAL_PERSONA FROM hysteria.CANAL_DIGITAL_PERSONA
											WHERE ID_PERSONA = $1 AND TIPO_CANAL_DIGITAL = $2)
		and api_key = $3
		`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, credentials.IdPersona, credentials.CanalDigital,
		credentials.ApiKey).Scan(&idToken, &fechaExpAccessToken)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("persona o canal digital inexistente")
		}
		return err
	}

	if fechaExpAccessToken.After(time.Now()) {
		return fmt.Errorf("inicia sesion de nuevo")
	}

	return nil
}

func (v SecurityRepository) MiddlewareValidations(ctx context.Context, credentials domains.Credentials) error {
	if err := v.checkCredentials(ctx, credentials); err != nil {
		return err
	}

	if err := v.checkRevokes(ctx, credentials); err != nil {
		return err
	}

	if err := v.tokenExpiration(ctx, credentials); err != nil {
		return err
	}
	return nil
}

func (v SecurityRepository) RevokePersona(ctx context.Context, revokePersona domains.RequestRevokePer) error {

	update := `update hysteria.persona set acceso_revocado = $1 where id_persona = $2`

	_, err := v.dbPost.GetDB().ExecContext(ctx, update, revokePersona.Revoke, revokePersona.IdPersonaRevoke)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) RevokeCanalDigital(ctx context.Context, revokeCanalDigital domains.RequestRevokeCanalDigital) error {
	update := `update hysteria.tipo_canal_digital_df set acceso_revocado = $1 where tipo_canal_digital = $2`

	_, err := v.dbPost.GetDB().ExecContext(ctx, update, revokeCanalDigital.Revoke, revokeCanalDigital.CanalDigitalRevoke)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) RevokeCanalDigPer(ctx context.Context, revokeCanalDigPer domains.RequestRevokeCanalDigPer) error {
	update := `update hysteria.canal_digital_persona set acceso_revocado = $1 where tipo_canal_digital = $2 
		and id_persona = $3`

	_, err := v.dbPost.GetDB().ExecContext(ctx, update, revokeCanalDigPer.Revoke, revokeCanalDigPer.CanalDigitalRevoke, revokeCanalDigPer.IdPersonaRevoke)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) LoginValidations(ctx context.Context, requestLogin domains.RequestLogin) (int, *string, error) {

	var idPersona int
	var hashedPassword string

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)

	if err != nil {
		return 0, nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `select id_persona,password_acceso_hash from hysteria.canal_digital_persona where tipo_canal_digital = $1 
		and login_name = $2`

	err = tx.QueryRowContext(ctx, query, requestLogin.CanalDigital, requestLogin.Username).Scan(&idPersona, &hashedPassword)

	if idPersona == 0 {
		return 0, nil, fmt.Errorf("usuario o canal digital incorrecto")
	}

	if err != nil {
		return idPersona, nil, err
	}

	if err = utils.ComparePasswordHash(hashedPassword, requestLogin.Password); err != nil {
		return idPersona, nil, fmt.Errorf("contraseña incorrecta")
	}

	credentials := domains.Credentials{
		IdPersona:    idPersona,
		ApiKey:       requestLogin.ApiKey,
		CanalDigital: requestLogin.CanalDigital,
	}

	if err := v.checkCredentials(ctx, credentials); err != nil {
		return idPersona, nil, err
	}

	if err := v.checkRevokes(ctx, credentials); err != nil {
		return idPersona, nil, err
	}

	seed2FA, err := v.CheckAPI2FA(ctx, idPersona, requestLogin.ApiKey, requestLogin.CanalDigital)

	if err != nil {
		return idPersona, nil, err
	}

	if err = tx.Commit(); err != nil {
		return idPersona, seed2FA, err
	}

	return idPersona, seed2FA, nil
}

func (v SecurityRepository) UpsertAccessToken(ctx context.Context, requestUpsert *domains.UpsertAccessToken) error {
	var idCanalDigitalPersona int

	expAccessToken, err := utils.GetTokenExpiration(requestUpsert.AccessToken, "ACCESS")

	if err != nil {
		return err
	}

	expRefreshToken, err := utils.GetTokenExpiration(requestUpsert.RefreshToken, "REFRESH")

	if err != nil {
		return err
	}

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT ID_CANAL_DIGITAL_PERSONA FROM hysteria.CANAL_DIGITAL_PERSONA
		WHERE ID_PERSONA = $1 AND TIPO_CANAL_DIGITAL = $2`

	err = tx.QueryRowContext(ctx, query, requestUpsert.IdPersona, requestUpsert.CanalDigital).Scan(&idCanalDigitalPersona)

	if err != nil {
		return err
	}

	query = `SELECT id_token FROM hysteria.TOKEN WHERE ID_CANAL_DIGITAL_PERSONA = $1 AND API_KEY = $2`

	rows, err := tx.QueryContext(ctx, query, idCanalDigitalPersona, requestUpsert.ApiKey)

	if err != nil {
		return err
	}

	if !rows.Next() {

		insert := `INSERT INTO hysteria.token 
		(api_key,id_canal_digital_persona,access_token,fecha_exp_access_token,refresh_token,fecha_Exp_refresh_token) VALUES ($1,$2,$3,$4,$5,$6)`

		_, err = tx.ExecContext(ctx, insert, requestUpsert.ApiKey, idCanalDigitalPersona, requestUpsert.AccessToken,
			expAccessToken, requestUpsert.RefreshToken, expRefreshToken)

		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		return nil
	}
	rows.Close()

	insert := `
    INSERT INTO hysteria.hist_token (id_token, api_key, id_canal_digital_persona, access_token, fecha_creacion_token, fecha_exp_access_token, refresh_token, fecha_exp_refresh_token, acceso_revocado)
    SELECT id_token, api_key, id_canal_digital_persona, access_token, fecha_creacion_token, fecha_exp_access_token, refresh_token, fecha_exp_refresh_token, acceso_revocado
    FROM hysteria.token
    WHERE id_canal_digital_persona = $1
	and api_key = $2`

	_, err = tx.ExecContext(ctx, insert, idCanalDigitalPersona, requestUpsert.ApiKey)

	if err != nil {
		return err
	}

	update := `update hysteria.token	set access_token = $1, fecha_creacion_token = $2, fecha_exp_access_token = $3, refresh_token = $4
		,fecha_exp_refresh_token = $5 
		where id_canal_digital_persona = $6 
		and api_key = $7`

	_, err = tx.ExecContext(ctx, update, requestUpsert.AccessToken, time.Now(), expAccessToken, requestUpsert.RefreshToken,
		expRefreshToken, idCanalDigitalPersona, requestUpsert.ApiKey)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) LogProcedure(ctx context.Context, logStruct *domains.LogStruct, mssgError string) error {
	var idToken int
	logStruct.MssgError = mssgError

	query := `SELECT id_token FROM hysteria.token
		WHERE api_key = $1 
		AND id_canal_digital_persona = (select id_canal_digital_persona 
											from hysteria.canal_digital_persona
											where id_persona = $2
											and tipo_canal_digital = $3)`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, logStruct.ApiKey, logStruct.IdPersona, logStruct.CanalDigital).Scan(&idToken)

	if err != nil {
		idToken = 0
	}

	logStruct.IdToken = idToken

	insert := `INSERT INTO hysteria.error_log (message_error, id_TIPO_ERROR, ID_PERSONA, canal_digital, api_key, id_token, access_token
		, ip_address, endpoint)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	_, err = v.dbPost.GetDB().ExecContext(ctx, insert, logStruct.MssgError, logStruct.IdTipoError, logStruct.IdPersona, logStruct.CanalDigital,
		logStruct.ApiKey, logStruct.IdToken, logStruct.Token, logStruct.IpAddress, logStruct.Endpoint)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) UpdCode2FA(ctx context.Context, login domains.RequestLogin, code int) error {

	update := `update hysteria.token
		set last_code_2fa = $1
		WHERE api_key = $2
		AND id_canal_digital_persona = (select id_canal_digital_persona 
											from hysteria.canal_digital_persona
											where login_name = $3
											and tipo_canal_digital = $4)`

	_, err := v.dbPost.GetDB().ExecContext(ctx, update, code, login.ApiKey, login.Username, login.CanalDigital)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) CheckLastRefreshToken(ctx context.Context, token string, credentials domains.Credentials) error {

	var idToken int

	query := `SELECT id_token FROM hysteria.token
		WHERE api_key = $1 
		AND id_canal_digital_persona = (select id_canal_digital_persona 
											from hysteria.canal_digital_persona
											where id_persona = $2
											and tipo_canal_digital = $3)
		and refresh_token = $4`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, credentials.ApiKey, credentials.IdPersona, credentials.CanalDigital,
		token).Scan(&idToken)

	if err != nil {
		return fmt.Errorf("token de refresco desconocido")
	}

	return nil
}

func (v SecurityRepository) CheckLastAccessToken(ctx context.Context, token string, credentials domains.Credentials) error {

	var idToken int

	query := `SELECT id_token FROM hysteria.token
		WHERE api_key = $1 
		AND id_canal_digital_persona = (select id_canal_digital_persona 
											from hysteria.canal_digital_persona
											where id_persona = $2
											and tipo_canal_digital = $3)
		and access_token = $4`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, credentials.ApiKey, credentials.IdPersona, credentials.CanalDigital,
		token).Scan(&idToken)

	if err != nil {
		return fmt.Errorf("el token de acceso no coincide con el ultimo registrado")
	}

	return nil
}

func (v SecurityRepository) GetAccessTokenDuration(ctx context.Context, apiKey string) (int, error) {

	var ctdHoras int

	query := `SELECT ctd_hs_access_token_valido FROM hysteria.api_key
		WHERE api_key = $1`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, apiKey).Scan(&ctdHoras)

	if err != nil {
		return 0, err
	}

	return ctdHoras * 60, nil
}

func (v SecurityRepository) Login2FA(ctx context.Context, login domains.RequestLogin2FA) (int, int, error) {
	var seed2FA string
	var idPersona int

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	ctdMins, err := v.GetAccessTokenDuration(ctx, login.ApiKey)

	if err != nil {
		return 0, 0, err
	}

	query := `select id_persona from hysteria.canal_digital_persona where login_name = $1`

	err = tx.QueryRowContext(ctx, query, login.Username).Scan(&idPersona)

	if err != nil {
		return 0, ctdMins, err
	}

	query = `SELECT "2fa_seed" from hysteria.token where api_key = $1
		and id_canal_digital_persona = (select id_canal_digital_persona from hysteria.canal_digital_persona where login_name = $2)`

	err = tx.QueryRowContext(ctx, query, login.ApiKey, login.Username).Scan(&seed2FA)

	if err != nil {
		return idPersona, ctdMins, err
	}

	if seed2FA == "" {
		return idPersona, ctdMins, fmt.Errorf("no se pudo recuperar semilla para el doble factor")
	}

	decrypted2FA, err := utils.DecryptTwo(login.Hash2FA, seed2FA)

	if err != nil {
		return idPersona, ctdMins, err
	}

	parts := strings.Split(decrypted2FA, ":")

	if len(parts) != 2 {
		return idPersona, ctdMins, fmt.Errorf("error en desencriptacion doble factor")
	}

	username := parts[0]

	if username != login.Username {
		return idPersona, ctdMins, fmt.Errorf("inconsistencia entre hash de doble factor y usuario")
	}

	isValid, err := utils.ValidateCredentialsAndTOTP(login.Code, seed2FA)

	if err != nil {
		return idPersona, ctdMins, fmt.Errorf("%s", err.Error())
	}

	if !isValid {
		return idPersona, ctdMins, fmt.Errorf("el codigo ingresado es incorrecto")
	}

	return idPersona, ctdMins, nil
}

func (v SecurityRepository) CheckAPI2FA(ctx context.Context, idPersona int, apiKey string, canalDigital string) (*string, error) {
	var reqApiKey string
	var reqUser string
	var seed2FAPointer *string
	var username string
	var seed2FAString string

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT req_2fa from hysteria.api_key where api_key = $1 `

	err = tx.QueryRowContext(ctx, query, apiKey).Scan(&reqApiKey)

	if err != nil {
		return nil, err
	}

	query = `SELECT req_2fa,login_name from hysteria.canal_digital_persona where id_persona = $1 and tipo_canal_digital = $2 `

	err = tx.QueryRowContext(ctx, query, idPersona, canalDigital).Scan(&reqUser, &username)

	if err != nil {
		return nil, err
	}

	if reqUser == "N" && reqApiKey == "N" {
		return nil, nil
	}

	query = `SELECT coalesce("2fa_seed",'NO TIENE') from hysteria.token where api_key = $1 and id_canal_digital_persona =
		(select id_canal_digital_persona from hysteria.canal_digital_persona where id_persona = $2 and tipo_canal_digital = $3)`

	err = tx.QueryRowContext(ctx, query, apiKey, idPersona, canalDigital).Scan(&seed2FAString)

	if err != nil {
		return nil, err
	}

	if seed2FAString == "NO TIENE" {

		seed2FA, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Thinksoft-autenticacion",
			AccountName: username,
		})

		if err != nil {
			return nil, err
		}

		update := `update hysteria.token set "2fa_seed" = $1 where api_key = $1 and id_canal_digital_persona =
		(select id_canal_digital_persona from hysteria.canal_digital_persona where id_persona = $2 and tipo_canal_digital = $3)`

		_, err = tx.ExecContext(ctx, update, seed2FA.Secret(), idPersona, canalDigital)

		if err != nil {
			return nil, err
		}

		seed2FAString = seed2FA.Secret()

		seed2FAPointer = &seed2FAString
	}

	if err = tx.Commit(); err != nil {
		return seed2FAPointer, err
	}

	seed2FAPointer = &seed2FAString

	return seed2FAPointer, nil
}

func (v SecurityRepository) RecuperacionPassword(ctx context.Context, recuperacionPassword domains.RequestRecuperacionPassword) error {

	return nil
}

func (v SecurityRepository) ValidarCanalDigital(ctx context.Context, validarCanalDigital domains.ValidarCanalDigital) error {

	update := `update hysteria.canal_digital_persona set canal_validado = 'S' where id_persona = $1 and tipo_canal_digital = $2`

	_, err := v.dbPost.GetDB().ExecContext(ctx, update, validarCanalDigital.IdPersona, validarCanalDigital.CanalDigital)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) checkSuperUser(ctx context.Context, apiKey string) (error, bool) {
	var isSuperUser string

	query := `select is_super_user from hysteria.api_key where api_key = $1`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, apiKey).Scan(&isSuperUser)

	if err != nil {
		return err, false
	}

	if isSuperUser == "N" {
		return nil, false
	}

	return nil, true
}

func (v SecurityRepository) CrearCanalDigital(ctx context.Context, crearCanalDigital domains.CrearCanalDigital, apiKey string) error {

	err, isSuperUser := v.checkSuperUser(ctx, apiKey)

	if err != nil {
		return err
	}

	if !isSuperUser {
		return fmt.Errorf("no posee los permisos necesarios para esta operacion")
	}

	insert := `insert into hysteria.tipo_canal_digital_df (tipo_canal_digital) values ($1)`

	_, err = v.dbPost.GetDB().ExecContext(ctx, insert, crearCanalDigital.CanalDigital)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) CambioPassword(ctx context.Context, cambioPassword domains.CambioPassword) error {

	var hashedActualPassword string

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	query := `select password_acceso_hash from hysteria.canal_digital_persona where id_persona = $1
		and tipo_canal_digital = $2`

	err = tx.QueryRowContext(ctx, query, cambioPassword.IdPersona, cambioPassword.CanalDigital).Scan(&hashedActualPassword)

	if err != nil {
		return err
	}

	if err := utils.ComparePasswordHash(hashedActualPassword, cambioPassword.ActualPassword); err != nil {
		return fmt.Errorf("su actual contraseña es incorrecta")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cambioPassword.NuevaPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	update := `update hysteria.canal_digital_persona set password_acceso_hash = $1 where id_persona = $2 and tipo_canal_digital = $3`

	_, err = tx.ExecContext(ctx, update, hashedPassword, cambioPassword.IdPersona, cambioPassword.CanalDigital)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) CambioPasswordByLogin(ctx context.Context, loginName string, newPassword string) error {

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	mailPersona, err := v.GetEmailByID(ctx, loginName)

	if err != nil {
		return err
	}

	if mailPersona == "" {
		return fmt.Errorf("no fue posible recuperar el correo para la recuperacion")
	}

	body := fmt.Sprintf("Hola, tu nueva contraseña es: %s", newPassword)

	if err = utils.SendEmail(mailPersona, "Cambio de contraseña", body); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	update := `update hysteria.canal_digital_persona set password_acceso_hash = $1 where login_name = $2`

	_, err = tx.ExecContext(ctx, update, hashedPassword, loginName)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) GetEmailByID(ctx context.Context, loginName string) (string, error) {

	var mailPersona string
	var idMailPersona int

	query := `SELECT id_mail_persona FROM hysteria.canal_digital_persona
		WHERE login_name = $1`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, loginName).Scan(&idMailPersona)

	if err != nil {
		return "", err
	}

	query = `SELECT mail FROM ts.mail_persona WHERE id_mail_persona = :1`

	err = v.dbOracl.GetDB().QueryRowContext(ctx, query, idMailPersona).Scan(&mailPersona)

	if err != nil {
		return "", err
	}

	return mailPersona, nil
}

func (v SecurityRepository) ActivarUser2FA(ctx context.Context, activarUser2FA domains.ActivarUser2FA, idPersona int, canalDigital string) error {

	update := `update hysteria.canal_digital_persona set req_2fa = $1 where id_persona = $2 and tipo_canal_digital = $3`

	_, err := v.dbPost.GetDB().ExecContext(ctx, update, activarUser2FA.Activo, idPersona, canalDigital)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) CheckApiKeyExpirada(ctx context.Context, apiKey string) (bool, error) {
	var fecha time.Time
	var apiKeyRevocada string

	query := `SELECT fecha_fin_vigencia FROM hysteria.api_key WHERE api_key = $1`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, apiKey).Scan(&fecha)

	if err != nil {
		return false, err
	}

	if fecha.Before(time.Now()) {
		return true, nil
	}

	query = `SELECT estado FROM hysteria.api_key WHERE api_key = $1`

	err = v.dbPost.GetDB().QueryRowContext(ctx, query, apiKey).Scan(&apiKeyRevocada)

	if err != nil {
		return false, err
	}

	if apiKeyRevocada == "INACTIVO" {
		return true, nil
	}

	return false, nil
}

func (v SecurityRepository) Generate2FAQR(ctx context.Context, generate2FAQR domains.Generate2FAQR) (string, error) {
	var seed2FAString string
	var idPersona int
	var canalDigital string
	var idCanalDigitalPersona int

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)

	if err != nil {
		return "", err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	if generate2FAQR.AccessToken != "" {

		claims, err := utils.GetClaimsFromToken(generate2FAQR.AccessToken, "ACCESS")

		if err != nil {
			return "", err
		}

		switch v := claims["id_persona"].(type) {
		case float64:
			idPersona = int(v) // jwt.MapClaims a menudo representa números como float64
		case int:
			idPersona = v
		case string:
			idPersona, _ = strconv.Atoi(v) // Si llega como string, convertirlo a int
		default:
			fmt.Println("Tipo inesperado para id_persona:", v)
		}

		canalDigital, _ = claims["canal_digital"].(string)

		query := `SELECT id_canal_digital_persona FROM hysteria.canal_digital_persona WHERE id_persona = $1
			and tipo_canal_digital = $2`

		err = tx.QueryRowContext(ctx, query, idPersona, canalDigital).Scan(&idCanalDigitalPersona)

		if err != nil {
			return "", err
		}
	}

	if idPersona == 0 && canalDigital == "" {
		query := `SELECT id_canal_digital_persona FROM hysteria.canal_digital_persona WHERE login_name = $1`

		err := tx.QueryRowContext(ctx, query, generate2FAQR.Username).Scan(&idCanalDigitalPersona)

		if err != nil {
			return "", err
		}
	}

	query := `SELECT COALESCE("2fa_seed", 'NO TIENE') FROM hysteria.token WHERE api_key = $1 and id_canal_digital_persona = $2`

	err = tx.QueryRowContext(ctx, query, generate2FAQR.ApiKey, idCanalDigitalPersona).Scan(&seed2FAString)

	if err != nil {
		return "", err
	}

	if seed2FAString == "NO TIENE" {

		seed2FA, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Thinksoft-autenticacion",
			AccountName: generate2FAQR.Username,
		})

		if err != nil {
			return "", err
		}

		update := `update hysteria.token set "2fa_seed" = $1 where api_key = $2 and id_canal_digital_persona = $3`

		_, err = tx.ExecContext(ctx, update, seed2FA.Secret(), generate2FAQR.ApiKey, idCanalDigitalPersona)

		if err != nil {
			return "", err
		}

		seed2FAString = seed2FA.Secret()

	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return seed2FAString, nil
}
