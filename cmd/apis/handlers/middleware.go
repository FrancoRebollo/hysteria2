package apis

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"plantilla_api/cmd/apis/handlers/seguridad"
	"plantilla_api/cmd/utils"
	"plantilla_api/internal/version/domains"
	"strings"

	"github.com/gin-gonic/gin"
)

func MiddlewareAuthorization(s seguridad.SecurityHandler) gin.HandlerFunc {
	return func(c *gin.Context) {

		var middlewareError error
		middlewareValidation := domains.Credentials{}
		accessToken := c.GetHeader("Authorization")
		accessBear := strings.TrimPrefix(accessToken, "Bearer ")

		var requestBody struct {
			IdPersona    int    `json:"id_persona"`
			ApiKey       string `json:"api_key"`
			CanalDigital string `json:"canal_digital"`
		}

		defer func() {
			if middlewareError != nil {
				if err := s.LogProcedure(c, requestBody, middlewareError.Error(), accessBear, 0); err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "imposible log errors"})
					c.Abort()
					return
				}
			}
		}()

		if accessBear == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inexistente"})
			c.Abort()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo leer el cuerpo de la solicitud"})
			middlewareError = err
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := json.Unmarshal(body, &requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error al leer el cuerpo de la solicitud"})
			middlewareError = err
			c.Abort()
			return
		}

		claims, err := utils.GetClaimsFromToken(accessBear, "ACCESS")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			middlewareError = err
			c.Abort()
			return
		}

		if idPersona, ok := claims["id_persona"].(float64); ok {
			middlewareValidation.IdPersona = int(idPersona)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no se pueden leer los claims del token"})
			middlewareError = err
			c.Abort()
			return
		}

		if apiKey, ok := claims["api_key"].(string); ok {
			middlewareValidation.ApiKey = apiKey
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no se pueden leer los claims del token"})
			middlewareError = err
			c.Abort()
			return
		}

		if canalDigital, ok := claims["canal_digital"].(string); ok {
			middlewareValidation.CanalDigital = canalDigital
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no se pueden leer los claims del token"})
			middlewareError = err
			c.Abort()
			return
		}

		if err = s.SecurityService.MiddlewareValidationsAPI(c, middlewareValidation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			middlewareError = err
			c.Abort()
			return
		}

		c.Next()
	}
}
