package seguridad

import (
	"fmt"
	"plantilla_api/internal/version/domains"
)

func ValidateAltaUser(altaUser domains.RequestAltaUser) error {
	if altaUser.IdPersona == 0 {
		return fmt.Errorf("es necesaria un id_persona para el alta")
	}

	if altaUser.CanalDigital == "" {
		return fmt.Errorf("se debe proveer un canal digital para el alta")
	}

	return nil
}

func ValidateGetJWT(requestJWT domains.RequestGetJWT) error {

	return nil
}

func ValidateCheckJWT(requestJWT domains.RequestCheckJWT) error {
	return nil
}
