package seguridad

import (
	"log"
	"net/http"
	"plantilla_api/cmd/utils"
	"plantilla_api/internal/version/domains"
	"plantilla_api/internal/version/ports"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SecurityHandler struct {
	ports.SecurityService
	*logrus.Logger
}

func NewSecurityHandler(service ports.SecurityService, loggerInstance *logrus.Logger) *SecurityHandler {
	if loggerInstance == nil {
		log.Fatal("es necesaria una instancia de logueo para que la app incie")
	}
	return &SecurityHandler{service, loggerInstance}
}

func (h *SecurityHandler) GetVersion(c *gin.Context) {
	var version_api domains.Version

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	version_api.NombreApi = "Seguridad"

	c.JSON(200, version_api)
}

func (h *SecurityHandler) AltaUser(c *gin.Context) {

	var altaUser domains.RequestAltaUser

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&altaUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateAltaUser(altaUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	credentialsExt := domains.CredentialsExtended{
		IdPersona:    altaUser.IdPersona,
		ApiKey:       "n/a",
		CanalDigital: altaUser.CanalDigital,
		IpAddress:    c.ClientIP(),
		Endpoint:     c.FullPath(),
	}

	accessToken := c.GetHeader("Authorization")
	accessBear := strings.TrimPrefix(accessToken, "Bearer ")

	alta_user_api, err := h.AltaUserAPI(c, altaUser)
	if err != nil {
		h.LogProcedure(c, credentialsExt, err.Error(), accessBear, 0)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, alta_user_api)
}

func (h *SecurityHandler) GetJWT(c *gin.Context) {

	var requestJWT domains.RequestGetJWT

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&requestJWT); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken := c.GetHeader("Authorization")
	accessBear := strings.TrimPrefix(accessToken, "Bearer ")

	credentialsExt := domains.CredentialsExtended{
		IdPersona:    0,
		ApiKey:       "n/a",
		CanalDigital: "",
		IpAddress:    c.ClientIP(),
		Endpoint:     c.FullPath(),
	}

	get_jwt_api, err := h.GetJWTAPI(c, requestJWT.RefreshToken, accessBear)
	if err != nil {
		h.LogProcedure(c, credentialsExt, err.Error(), requestJWT.RefreshToken, 0)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, get_jwt_api)
}

func (h *SecurityHandler) ValidateJWT(c *gin.Context) {

	var idPersona int

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	accessToken := c.GetHeader("Authorization")
	accessBear := strings.TrimPrefix(accessToken, "Bearer ")

	claims, err := utils.GetClaimsFromToken(accessBear, "ACCESS")

	apiKey, _ := claims["api_key"].(string)

	switch v := claims["id_persona"].(type) {
	case float64:
		idPersona = int(v)
	case int:
		idPersona = v
	case string:
		idPersona, _ = strconv.Atoi(v)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read claims"})
		return
	}

	canalDigital, _ := claims["canal_digital"].(string)

	credentialsExt := domains.CredentialsExtended{
		IdPersona:    idPersona,
		ApiKey:       apiKey,
		CanalDigital: canalDigital,
		IpAddress:    c.ClientIP(),
		Endpoint:     c.FullPath(),
	}

	if err != nil {
		h.LogProcedure(c, credentialsExt, err.Error(), accessBear, 0)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	checkJWTResponse, err := h.ValidateJWTAPI(c, accessBear)

	if err != nil {
		h.LogProcedure(c, credentialsExt, err.Error(), accessBear, 0)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	checkJWTResponse.IdPersona = idPersona

	c.JSON(200, checkJWTResponse)
}

func (h *SecurityHandler) RevokePersona(c *gin.Context) {

	var revokePersona domains.RequestRevokePer
	var message string
	var idPersona int

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	accessToken := c.GetHeader("Authorization")
	accessBear := strings.TrimPrefix(accessToken, "Bearer ")

	claims, err := utils.GetClaimsFromToken(accessBear, "ACCESS")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claims"})
		return
	}

	apiKey, _ := claims["api_key"].(string)

	apiKeyExp, err := h.CheckApiKeyExpiradaAPI(c, apiKey)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	if apiKeyExp {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Api key expirada o revocada"})
		return
	}

	switch v := claims["id_persona"].(type) {
	case float64:
		idPersona = int(v)
	case int:
		idPersona = v
	case string:
		idPersona, _ = strconv.Atoi(v)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read claims"})
		return
	}
	canalDigital, _ := claims["canal_digital"].(string)

	if err := c.BindJSON(&revokePersona); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.RevokePersonaAPI(c, revokePersona)

	if err != nil {

		credentialsExt := domains.CredentialsExtended{
			IdPersona:    idPersona,
			ApiKey:       apiKey,
			CanalDigital: canalDigital,
			IpAddress:    c.ClientIP(),
			Endpoint:     c.FullPath(),
		}

		//Cenitnela
		h.LogProcedure(c, credentialsExt, err.Error(), accessBear, 0)
		c.JSON(500, gin.H{"message": "Error que viene del servicio"})
		return
	}

	if revokePersona.Revoke == "N" {
		message = "persona no revocada"
	}

	if revokePersona.Revoke == "S" {
		message = "persona revocada"
	}

	resp := domains.DefaultResponse{
		Message: message,
	}

	c.JSON(200, resp)
}

func (h *SecurityHandler) RevokeCanalDigital(c *gin.Context) {

	var revokeCanalDigital domains.RequestRevokeCanalDigital
	var message string
	var idPersona int

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	accessToken := c.GetHeader("Authorization")
	accessBear := strings.TrimPrefix(accessToken, "Bearer ")

	claims, err := utils.GetClaimsFromToken(accessBear, "ACCESS")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claims"})
		return
	}

	apiKey, _ := claims["api_key"].(string)

	apiKeyExp, err := h.CheckApiKeyExpiradaAPI(c, apiKey)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	if apiKeyExp {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Api key expirada o revocada"})
		return
	}

	switch v := claims["id_persona"].(type) {
	case float64:
		idPersona = int(v)
	case int:
		idPersona = v
	case string:
		idPersona, _ = strconv.Atoi(v)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read claims"})
		return
	}
	canalDigital, _ := claims["canal_digital"].(string)

	if err := c.BindJSON(&revokeCanalDigital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.RevokeCanalDigitalAPI(c, revokeCanalDigital)
	if err != nil {

		credentialsExt := domains.CredentialsExtended{
			IdPersona:    idPersona,
			ApiKey:       apiKey,
			CanalDigital: canalDigital,
			IpAddress:    c.ClientIP(),
			Endpoint:     c.FullPath(),
		}
		//Cenitnela
		h.LogProcedure(c, credentialsExt, err.Error(), accessBear, 0)
		c.JSON(500, gin.H{"message": "Error que viene del servicio"})
		return
	}

	if revokeCanalDigital.Revoke == "N" {
		message = "canal digital no revocado"
	}

	if revokeCanalDigital.Revoke == "S" {
		message = "canal digital revocado"
	}

	resp := domains.DefaultResponse{
		Message: message,
	}

	c.JSON(200, resp)
}

func (h *SecurityHandler) RevokeCanalDigitalPersona(c *gin.Context) {

	var revokeCanalDigitalPer domains.RequestRevokeCanalDigPer
	var message string
	var idPersona int

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}
	accessToken := c.GetHeader("Authorization")
	accessBear := strings.TrimPrefix(accessToken, "Bearer ")

	claims, err := utils.GetClaimsFromToken(accessBear, "ACCESS")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claims"})
		return
	}

	apiKey, _ := claims["api_key"].(string)

	apiKeyExp, err := h.CheckApiKeyExpiradaAPI(c, apiKey)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	if apiKeyExp {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Api key expirada o revocada"})
		return
	}

	switch v := claims["id_persona"].(type) {
	case float64:
		idPersona = int(v)
	case int:
		idPersona = v
	case string:
		idPersona, _ = strconv.Atoi(v)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read claims"})
		return
	}

	canalDigital, _ := claims["canal_digital"].(string)

	if err := c.BindJSON(&revokeCanalDigitalPer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.RevokeCanalDigitalPersonaAPI(c, revokeCanalDigitalPer)
	if err != nil {

		credentialsExt := domains.CredentialsExtended{
			IdPersona:    idPersona,
			ApiKey:       apiKey,
			CanalDigital: canalDigital,
			IpAddress:    c.ClientIP(),
			Endpoint:     c.FullPath(),
		}

		//Cenitnela
		h.LogProcedure(c, credentialsExt, err.Error(), accessBear, 0)
		c.JSON(500, gin.H{"message": "Error que viene del servicio"})
		return
	}

	if revokeCanalDigitalPer.Revoke == "N" {
		message = "persona no revocada en el canal digital"
	}

	if revokeCanalDigitalPer.Revoke == "S" {
		message = "persona revocada en el canal digital"
	}

	resp := domains.DefaultResponse{
		Message: message,
	}

	c.JSON(200, resp)
}

func (h *SecurityHandler) Login(c *gin.Context) {

	var login domains.RequestLogin
	var loginPre domains.RequestLoginDos

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&loginPre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	login.ApiKey = c.GetHeader("Api-Key")
	login.Username = loginPre.Username
	login.CanalDigital = loginPre.CanalDigital
	login.Password = loginPre.Password

	resp, IdPersona, err := h.LoginAPI(c, login)

	if err != nil {

		credentialsExt := domains.CredentialsExtended{
			IdPersona:    IdPersona,
			ApiKey:       login.ApiKey,
			CanalDigital: login.CanalDigital,
			IpAddress:    c.ClientIP(),
			Endpoint:     c.FullPath(),
		}

		h.LogProcedure(c, credentialsExt, err.Error(), "", IdPersona)
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, resp)
}

func (h *SecurityHandler) Login2FA(c *gin.Context) {

	var login2FA domains.RequestLogin2FADos
	var loginResponse *domains.LoginResponse
	var login2FADos domains.RequestLogin2FA

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&login2FA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	login2FADos.ApiKey = c.GetHeader("Api-Key")
	login2FADos.CanalDigital = login2FA.CanalDigital
	login2FADos.Code = login2FA.Code
	login2FADos.Hash2FA = login2FA.Hash2FA
	login2FADos.Username = login2FA.Username

	_, loginResponse, err := h.Login2FAAPI(c, login2FADos)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, loginResponse)
}

func (h *SecurityHandler) RecuperacionPassword(c *gin.Context) {

	var recuperacionPassword domains.RequestRecuperacionPasswordDos
	var recuperacionPasswordDos domains.RequestRecuperacionPassword

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&recuperacionPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recuperacionPasswordDos.LoginName = recuperacionPassword.LoginName
	recuperacionPasswordDos.ApiKey = c.GetHeader("Api-Key")

	apiKeyExpirada, err := h.CheckApiKeyExpiradaAPI(c, recuperacionPasswordDos.ApiKey)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if apiKeyExpirada {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "api key expirada o revocada"})
		return
	}

	if err := h.RecuperacionPasswordAPI(c, recuperacionPasswordDos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := domains.DefaultResponse{
		Message: "Una nueva contraseña fue enviada, verifique su correo electronico",
	}

	c.JSON(200, resp)
}

func (h *SecurityHandler) ValidarCanalDigital(c *gin.Context) {
	var validarCanalDigital domains.ValidarCanalDigital

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&validarCanalDigital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ValidarCanalDigitalAPI(c, validarCanalDigital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := domains.DefaultResponse{
		Message: "Canal digital validado",
	}

	c.JSON(200, resp)
}

func (h *SecurityHandler) CrearCanalDigital(c *gin.Context) {
	var crearCanalDigital domains.CrearCanalDigital

	apiKey := c.GetHeader("Api-Key")

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&crearCanalDigital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.CrearCanalDigitalAPI(c, crearCanalDigital, apiKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := domains.DefaultResponse{
		Message: "Canal digital creado",
	}

	c.JSON(200, resp)
}

func (h *SecurityHandler) CambioPassword(c *gin.Context) {
	var cambioPassword domains.CambioPassword

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&cambioPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.CambioPasswordAPI(c, cambioPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := domains.DefaultResponse{
		Message: "Password modificada exitosamente",
	}

	c.JSON(200, resp)
}

func (h *SecurityHandler) ActivarUser2FA(c *gin.Context) {
	var activarUser2FA domains.ActivarUser2FA
	var idPersona int
	var message string

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&activarUser2FA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken := c.GetHeader("Authorization")
	accessBear := strings.TrimPrefix(accessToken, "Bearer ")

	claims, err := utils.GetClaimsFromToken(accessBear, "ACCESS")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claims"})
		return
	}

	switch v := claims["id_persona"].(type) {
	case float64:
		idPersona = int(v)
	case int:
		idPersona = v
	case string:
		idPersona, _ = strconv.Atoi(v)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read claims"})
		return
	}

	canalDigital, _ := claims["canal_digital"].(string)

	if err := h.ActivarUser2FAAPI(c, activarUser2FA, idPersona, canalDigital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if activarUser2FA.Activo == "S" {
		message = "Doble factor activado exitosamente a nivel de usuario"
	}

	if activarUser2FA.Activo == "N" {
		message = "Doble factor desactivado exitosamente a nivel de usuario"
	}

	resp := domains.DefaultResponse{
		Message: message,
	}

	c.JSON(200, resp)
}

func (h *SecurityHandler) Generate2FAQR(c *gin.Context) {
	var generate2FAQR domains.Generate2FAQRDos
	var seed2FA string
	var generate2FAQRDos domains.Generate2FAQR

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&generate2FAQR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken := c.GetHeader("Authorization")
	accessBear := strings.TrimPrefix(accessToken, "Bearer ")

	claims, err := utils.GetClaimsFromToken(accessBear, "ACCESS")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	apiKey, _ := claims["api_key"].(string)

	generate2FAQRDos.AccessToken = accessBear
	generate2FAQRDos.ApiKey = apiKey
	generate2FAQRDos.Username = generate2FAQR.Username

	seed2FA, err = h.Generate2FAQRAPI(c, generate2FAQRDos)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urlCode, base64Code, err := utils.GenerateQRCode(generate2FAQRDos.Username, seed2FA)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := domains.QRBase64UrlResponse{
		Url:        urlCode,
		CodeBase64: base64Code,
		ManualSeed: seed2FA,
	}

	c.JSON(200, resp)
}
