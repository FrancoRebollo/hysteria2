package domains

type RequestUser struct {
	ApiKey    string `json:"api_key"`
	TipoCanal string `json:"tipo_canal"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type RequestLogin struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ApiKey       string `json:"api_key"`
	CanalDigital string `json:"canal_digital"`
}

type RequestLoginDos struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CanalDigital string `json:"canal_digital"`
}

type RequestAltaUser struct {
	IdPersona     int    `json:"id_persona"`
	CanalDigital  string `json:"canal_digital"`
	LoginName     string `json:"login_name"`
	Password      string `json:"password"`
	IdMailPersona int    `json:"id_mail_persona"`
	IdTePersona   int    `json:"id_te_persona"`
}

type RequestGetJWT struct {
	RefreshToken string `json:"refresh_token"`
}

type RequestCheckJWT struct {
	TokenJWT     string `json:"token_jwt"`
	ApiKey       string `json:"api_key"`
	IdPersona    int    `json:"id_persona"`
	CanalDigital string `json:"canal_digital"`
}

type RequestRevokePer struct {
	IdPersonaRevoke int    `json:"id_persona_revoke"`
	Revoke          string `json:"revoke"`
}

type RequestRevokeCanalDigital struct {
	CanalDigitalRevoke string `json:"canal_digital_revoke"`
	Revoke             string `json:"revoke"`
}

type RequestRevokeCanalDigPer struct {
	IdPersonaRevoke    int    `json:"id_persona_revoke"`
	CanalDigitalRevoke string `json:"canal_digital_revoke"`
	Revoke             string `json:"revoke"`
}

type RequestLogin2FA struct {
	Hash2FA      string `json:"hash_2fa"`
	Code         string `json:"code"`
	ApiKey       string `json:"api_key"`
	CanalDigital string `json:"canal_digital"`
	Username     string `json:"username"`
}

type RequestLogin2FADos struct {
	Hash2FA      string `json:"hash_2fa"`
	Code         string `json:"code"`
	CanalDigital string `json:"canal_digital"`
	Username     string `json:"username"`
}

type RequestRecuperacionPassword struct {
	LoginName string `json:"login_name"`
	ApiKey    string `json:"api_key"`
}

type RequestRecuperacionPasswordDos struct {
	LoginName string `json:"login_name"`
}

type ValidarCanalDigital struct {
	IdPersona    string `json:"id_persona"`
	CanalDigital string `json:"canal_digital"`
}

type CrearCanalDigital struct {
	CanalDigital string `json:"canal_digital"`
}

type CambioPassword struct {
	IdPersona      int    `json:"id_persona"`
	CanalDigital   string `json:"canal_digital"`
	ActualPassword string `json:"actual_password"`
	NuevaPassword  string `json:"nueva_password"`
}

type ActivarUser2FA struct {
	Activo string `json:"activo"`
}

type Generate2FAQR struct {
	AccessToken string
	Hash2FA     string
	Username    string
	ApiKey      string
}

type Generate2FAQRDos struct {
	Username string `json:"username"`
}
