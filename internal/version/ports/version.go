package ports

import (
	"context"
	"plantilla_api/internal/version/domains"
)

// Interface para el manejo de versiones
type SecurityService interface {
	GetVersionAPI(ctx context.Context) (*domains.Version, error)
	AltaUserAPI(ctx context.Context, altaUser domains.RequestAltaUser) (*domains.AltaUserResponse, error)
	GetJWTAPI(ctx context.Context, refreshToken string, accessTokenParam string) (*domains.GetJWTResponse, error)
	ValidateJWTAPI(ctx context.Context, token string) (*domains.CheckJWTResponse, error)
	MiddlewareValidationsAPI(ctx context.Context, middlewareValidation domains.Credentials) error
	RevokePersonaAPI(ctx context.Context, revokePersona domains.RequestRevokePer) error
	RevokeCanalDigitalAPI(ctx context.Context, revokeCanalDigital domains.RequestRevokeCanalDigital) error
	RevokeCanalDigitalPersonaAPI(ctx context.Context, revokeCanalDigPer domains.RequestRevokeCanalDigPer) error
	LoginAPI(ctx context.Context, requestLogin domains.RequestLogin) (*domains.LoginResponse, int, error)
	LogProcedure(ctx context.Context, credentials any, mssgError string, token string, idPersona int) error
	Login2FAAPI(ctx context.Context, login domains.RequestLogin2FA) (int, *domains.LoginResponse, error)
	RecuperacionPasswordAPI(ctx context.Context, recuperacionPassword domains.RequestRecuperacionPassword) error
	ValidarCanalDigitalAPI(ctx context.Context, validarCanalDigital domains.ValidarCanalDigital) error
	CrearCanalDigitalAPI(ctx context.Context, crearCanalDigital domains.CrearCanalDigital, apiKey string) error
	CambioPasswordAPI(ctx context.Context, cambioPassword domains.CambioPassword) error
	ActivarUser2FAAPI(ctx context.Context, activarUser2FA domains.ActivarUser2FA, idPersona int, canalDigital string) error
	Generate2FAQRAPI(ctx context.Context, generate2FAQR domains.Generate2FAQR) (string, error)
	CheckApiKeyExpiradaAPI(ctx context.Context, apiKey string) (bool, error)
}

// Interface para consultas a la DB
type SecurityRepository interface {
	GetVersion(ctx context.Context) (string, error)
	AltaUser(ctx context.Context, altaUser domains.RequestAltaUser) (string, error)
	CheckTokenCreation(ctx context.Context, credentials domains.Credentials) error
	PersistToken(ctx context.Context, credentialsToken domains.CredentialsToken) error
	MiddlewareValidations(ctx context.Context, middlewareValidation domains.Credentials) error
	RevokePersona(ctx context.Context, revokePersona domains.RequestRevokePer) error
	RevokeCanalDigital(ctx context.Context, revokeCanalDigital domains.RequestRevokeCanalDigital) error
	RevokeCanalDigPer(ctx context.Context, revokeCanalDigPer domains.RequestRevokeCanalDigPer) error
	LoginValidations(ctx context.Context, requestLogin domains.RequestLogin) (int, *string, error)
	UpsertAccessToken(ctx context.Context, requestUpsert *domains.UpsertAccessToken) error
	LogProcedure(ctx context.Context, logStruct *domains.LogStruct, mssgError string) error
	CheckLastRefreshToken(ctx context.Context, token string, credentials domains.Credentials) error
	CheckLastAccessToken(ctx context.Context, token string, credentials domains.Credentials) error
	GetAccessTokenDuration(ctx context.Context, ApiKey string) (int, error)
	Login2FA(ctx context.Context, login domains.RequestLogin2FA) (int, int, error)
	CheckAPI2FA(ctx context.Context, idPersona int, apiKey string, canalDigital string) (*string, error)
	UpdCode2FA(ctx context.Context, login domains.RequestLogin, code int) error
	RecuperacionPassword(ctx context.Context, recuperacionPassword domains.RequestRecuperacionPassword) error
	ValidarCanalDigital(ctx context.Context, validarCanalDigital domains.ValidarCanalDigital) error
	CrearCanalDigital(ctx context.Context, crearCanalDigital domains.CrearCanalDigital, apiKey string) error
	CambioPassword(ctx context.Context, cambioPassword domains.CambioPassword) error
	GetEmailByID(ctx context.Context, loginName string) (string, error)
	CambioPasswordByLogin(ctx context.Context, loginName string, newPassword string) error
	ActivarUser2FA(ctx context.Context, activarUser2FA domains.ActivarUser2FA, idPersona int, canalDigital string) error
	Generate2FAQR(ctx context.Context, generate2FAQR domains.Generate2FAQR) (string, error)
	CheckApiKeyExpirada(ctx context.Context, apiKey string) (bool, error)
}

type HysteriaService interface {
	AltaBossAPI(ctx context.Context, AltaBoss domains.RequestAltaBoss) error
}

type HysteriaRepository interface {
	AltaBoss(ctx context.Context, AltaBoss domains.RequestAltaBoss) (string, error)
}
