package domains

import (
	"time"
)

type CredentialsExtended struct {
	IdPersona    int
	ApiKey       string
	CanalDigital string
	IpAddress    string
	Endpoint     string
}

type CredentialsToken struct {
	IdPersona    int
	ApiKey       string
	CanalDigital string
	AccessToken  string
	RefreshToken string
}

type DefaultResponse struct {
	Message string `json:"Message"`
}

type Version struct {
	NombreApi     string `json:"nombre_api"`
	Cliente       string `json:"cliente"`
	Version       string `json:"version"`
	VersionModelo string `json:"version_modelo"`
	FechaStartUp  string `json:"fecha_start_up"`
}

type LogStruct struct {
	IdTipoError  int
	MssgError    string
	IdPersona    int
	CanalDigital string
	ApiKey       string
	IdToken      int
	Token        string
	Endpoint     string
	IpAddress    string
}

type AltaUserResponse struct {
	IdPersona    int    `json:"id_persona"`
	CanalDigital string `json:"canal_digital"`
	Message      string `json:"message"`
}

type GetJWTResponse struct {
	AccessToken string `json:"access_token"`
}

type CheckJWTResponse struct {
	IdPersona   int    `json:"id_persona"`
	TokenStatus string `json:"token_status"`
}

type Credentials struct {
	IdPersona    int
	ApiKey       string
	CanalDigital string
}

type LoginResponse struct {
	Username     string `json:"username"`
	Status       string `json:"status"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Hash2FA      string `json:"hash_2fa"`
}

type UpsertAccessToken struct {
	IdPersona    int
	CanalDigital string
	ApiKey       string
	AccessToken  string
	RefreshToken string
}

type QRBase64UrlResponse struct {
	Url        string `json:"url"`
	CodeBase64 string `json:"code_base_64"`
	ManualSeed string `json:"manual_seed"`
}

// ----------------------------------------------//

type AltaBossResponse struct {
	//IdBosses int    `json:"id_bosses"`
	Nombre  string `json:"nombre"`
	Message string `json:"message"`
}

type AltaAnuncioResponse struct {
	Id    int       `json:"id"`
	Texto string    `json:"texto"`
	Fecha time.Time `json:"fecha"`
	Error string    `json:"error"`
}
