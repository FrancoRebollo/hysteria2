package services

import (
	"context"
	"fmt"
	"os"
	"plantilla_api/cmd/utils"
	"plantilla_api/internal/version/domains"
	"reflect"
	"strconv"
	"time"
)

func (s *SecurityService) AltaUserAPI(ctx context.Context, altaUser domains.RequestAltaUser) (*domains.AltaUserResponse, error) {

	// valido con web service externo

	message, err := s.SecurityRepository.AltaUser(ctx, altaUser)

	if err != nil {
		return nil, err
	}

	altaResponse := domains.AltaUserResponse{
		IdPersona:    altaUser.IdPersona,
		CanalDigital: altaUser.CanalDigital,
		Message:      message,
	}

	return &altaResponse, nil

}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) GetJWTAPI(ctx context.Context, refreshToken string, accessTokenParam string) (*domains.GetJWTResponse, error) {

	expirationTime, err := utils.GetTokenExpiration(refreshToken, "REFRESH")

	if err != nil {
		return nil, err
	}

	if expirationTime != nil && expirationTime.Before(time.Now()) {
		return nil, fmt.Errorf("inicie sesion nuevamente")
	}

	claims, err := utils.GetClaimsFromToken(refreshToken, "REFRESH")

	if err != nil {
		return nil, err
	}

	credentials := domains.Credentials{
		IdPersona:    int(claims["id_persona"].(float64)),
		ApiKey:       claims["api_key"].(string),
		CanalDigital: claims["canal_digital"].(string),
	}

	if err := s.SecurityRepository.CheckTokenCreation(ctx, credentials); err != nil {
		return nil, err
	}

	if err := s.SecurityRepository.CheckLastRefreshToken(ctx, refreshToken, credentials); err != nil {
		return nil, err
	}

	if err := s.SecurityRepository.CheckLastAccessToken(ctx, accessTokenParam, credentials); err != nil {
		return nil, err
	}

	ctdMins, err := s.SecurityRepository.GetAccessTokenDuration(ctx, credentials.ApiKey)

	if err != nil {
		return nil, err
	}

	accessToken, err := utils.JWTCreate(ctdMins, credentials, "ACCESS")

	if err != nil {
		accessToken = "error en creacion"
	}

	credentialsToken := domains.CredentialsToken{
		IdPersona:    credentials.IdPersona,
		ApiKey:       credentials.ApiKey,
		CanalDigital: credentials.CanalDigital,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err := s.SecurityRepository.PersistToken(ctx, credentialsToken); err != nil {
		return nil, err
	}

	resp := &domains.GetJWTResponse{
		AccessToken: accessToken,
	}

	return resp, nil

}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) ValidateJWTAPI(ctx context.Context, token string) (*domains.CheckJWTResponse, error) {

	checkJWTResponse, err := utils.CheckJWTAccessToken(token)

	if err != nil {
		return nil, err
	}

	return checkJWTResponse, nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) MiddlewareValidationsAPI(ctx context.Context, credentials domains.Credentials) error {
	if err := s.SecurityRepository.MiddlewareValidations(ctx, credentials); err != nil {
		return err
	}

	return nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) RevokePersonaAPI(ctx context.Context, revokePersona domains.RequestRevokePer) error {
	if err := s.SecurityRepository.RevokePersona(ctx, revokePersona); err != nil {
		return err
	}

	return nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) RevokeCanalDigitalAPI(ctx context.Context, revokeCanalDigital domains.RequestRevokeCanalDigital) error {
	if err := s.SecurityRepository.RevokeCanalDigital(ctx, revokeCanalDigital); err != nil {
		return err
	}

	return nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) RevokeCanalDigitalPersonaAPI(ctx context.Context, revokeCanalDigPer domains.RequestRevokeCanalDigPer) error {
	if err := s.SecurityRepository.RevokeCanalDigPer(ctx, revokeCanalDigPer); err != nil {
		return err
	}

	return nil
}

func (s *SecurityService) LogProcedure(ctx context.Context, credentialsExt any, mssgError string, token string, idPersona int) error {
	v := reflect.ValueOf(credentialsExt)

	logStruct := &domains.LogStruct{Token: token}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("es necesario un struct para loguear errores")
	}

	if idPersona != 0 {
		logStruct.IdPersona = idPersona
	} else if field := v.FieldByName("IdPersona"); field.IsValid() && field.CanInt() {
		logStruct.IdPersona = int(field.Int())
	}

	if apiKeyField := v.FieldByName("ApiKey"); apiKeyField.IsValid() {
		logStruct.ApiKey = apiKeyField.String()
	}

	if canalDigitalField := v.FieldByName("CanalDigital"); canalDigitalField.IsValid() {
		logStruct.CanalDigital = canalDigitalField.String()
	}

	if ipAddress := v.FieldByName("IpAddress"); ipAddress.IsValid() {
		logStruct.IpAddress = ipAddress.String()
	}

	if endpoint := v.FieldByName("Endpoint"); endpoint.IsValid() {
		logStruct.Endpoint = endpoint.String()
	}

	if err := s.SecurityRepository.LogProcedure(ctx, logStruct, mssgError); err != nil {
		return err
	}

	return nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) LoginAPI(ctx context.Context, requestLogin domains.RequestLogin) (*domains.LoginResponse, int, error) {

	idPersona, seed2FA, err := s.SecurityRepository.LoginValidations(ctx, requestLogin)

	if err != nil {
		return nil, idPersona, err
	}

	if seed2FA != nil {

		encrypted2FA, err := utils.EncryptTwo(requestLogin.Username+":"+requestLogin.Password, *seed2FA)
		if err != nil {
			return nil, idPersona, err
		}

		resp := &domains.LoginResponse{
			Username:     requestLogin.Username,
			Status:       "Ingrese el codigo de seguridad de su aplicacion",
			RefreshToken: "",
			AccessToken:  "",
			Hash2FA:      encrypted2FA,
		}

		return resp, idPersona, nil

	}

	credentials := domains.Credentials{
		IdPersona:    idPersona,
		CanalDigital: requestLogin.CanalDigital,
		ApiKey:       requestLogin.ApiKey,
	}

	ctdMins, err := s.SecurityRepository.GetAccessTokenDuration(ctx, credentials.ApiKey)

	if err != nil {
		return nil, idPersona, err
	}
	//
	accessToken, err := utils.JWTCreate(ctdMins, credentials, "ACCESS")

	if err != nil {
		accessToken = "error en creacion"
	}

	refreshDuration, err := strconv.Atoi(os.Getenv("REF_TOKEN_DURATION"))

	if err != nil {
		return nil, idPersona, err
	}

	refreshToken, err := utils.JWTCreate(refreshDuration, credentials, "REFRESH")

	if err != nil {
		refreshToken = "error en creacion"
	}

	upsertAccessToken := &domains.UpsertAccessToken{
		IdPersona:    credentials.IdPersona,
		CanalDigital: credentials.CanalDigital,
		ApiKey:       credentials.ApiKey,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err := s.SecurityRepository.UpsertAccessToken(ctx, upsertAccessToken); err != nil {
		return nil, idPersona, err
	}

	resp := &domains.LoginResponse{
		Username:     requestLogin.Username,
		Status:       "Logged",
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		Hash2FA:      "",
	}

	return resp, idPersona, nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) Login2FAAPI(ctx context.Context, login domains.RequestLogin2FA) (int, *domains.LoginResponse, error) {
	idPersona, ctdMins, err := s.SecurityRepository.Login2FA(ctx, login)

	if err != nil {
		return idPersona, nil, err
	}

	credentials := domains.Credentials{
		IdPersona:    idPersona,
		ApiKey:       login.ApiKey,
		CanalDigital: login.CanalDigital,
	}

	if ctdMins == 0 {
		return idPersona, nil, fmt.Errorf("no se pudo recuperar la duracion del token de acceso")
	}

	accessToken, err := utils.JWTCreate(ctdMins, credentials, "ACCESS")

	if err != nil {
		accessToken = "error en creacion"
	}

	refreshDuration, err := strconv.Atoi(os.Getenv("REF_TOKEN_DURATION"))

	if err != nil {
		return idPersona, nil, err
	}

	refreshToken, err := utils.JWTCreate(refreshDuration, credentials, "REFRESH")

	if err != nil {
		refreshToken = "error en creacion"
	}

	resp := &domains.LoginResponse{
		Username:     login.Username,
		Status:       "Logged",
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		Hash2FA:      "",
	}

	upsertAccessToken := &domains.UpsertAccessToken{
		IdPersona:    credentials.IdPersona,
		CanalDigital: credentials.CanalDigital,
		ApiKey:       credentials.ApiKey,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err := s.SecurityRepository.UpsertAccessToken(ctx, upsertAccessToken); err != nil {
		return idPersona, resp, err
	}

	return idPersona, resp, nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) RecuperacionPasswordAPI(ctx context.Context, recuperacionPassword domains.RequestRecuperacionPassword) error {

	newPassword, err := utils.GenerateRandomPassword(16)

	if err != nil {
		return fmt.Errorf("no fue posible generar una nueva contraseña")
	}

	if err := s.SecurityRepository.CambioPasswordByLogin(ctx, recuperacionPassword.LoginName, newPassword); err != nil {
		return err
	}

	return nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) ValidarCanalDigitalAPI(ctx context.Context, validarCanalDigital domains.ValidarCanalDigital) error {
	if err := s.SecurityRepository.ValidarCanalDigital(ctx, validarCanalDigital); err != nil {
		return err
	}

	return nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) CrearCanalDigitalAPI(ctx context.Context, crearCanalDigital domains.CrearCanalDigital, apiKey string) error {
	if err := s.SecurityRepository.CrearCanalDigital(ctx, crearCanalDigital, apiKey); err != nil {
		return err
	}

	return nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) CambioPasswordAPI(ctx context.Context, cambioPassword domains.CambioPassword) error {
	if err := s.SecurityRepository.CambioPassword(ctx, cambioPassword); err != nil {
		return err
	}

	return nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) ActivarUser2FAAPI(ctx context.Context, activarUser2FA domains.ActivarUser2FA, idPersona int, canalDigital string) error {
	if err := s.SecurityRepository.ActivarUser2FA(ctx, activarUser2FA, idPersona, canalDigital); err != nil {
		return err
	}

	return nil
}

// MANEJO TRANSACCIONAL HECHO
func (s *SecurityService) Generate2FAQRAPI(ctx context.Context, generate2FAQR domains.Generate2FAQR) (string, error) {

	seed2FA, err := s.SecurityRepository.Generate2FAQR(ctx, generate2FAQR)

	if err != nil {
		return "", err
	}

	return seed2FA, nil
}

func (s *SecurityService) CheckApiKeyExpiradaAPI(ctx context.Context, apiKey string) (bool, error) {

	bool, err := s.SecurityRepository.CheckApiKeyExpirada(ctx, apiKey)

	if err != nil {
		return false, err
	}

	if !bool {
		return false, nil
	}

	return true, nil
}
