package hysteria

import (
	"log"
	"net/http"
	"plantilla_api/cmd/utils"
	"plantilla_api/internal/version/domains"
	"plantilla_api/internal/version/ports"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HysteriaHandler struct {
	ports.HysteriaService
	*logrus.Logger
}

func NewHysteriaHandler(service ports.SecurityService, loggerInstance *logrus.Logger) *HysteriaHandler {
	if loggerInstance == nil {
		log.Fatal("es necesaria una instancia de logueo para que la app incie")
	}
	return &HysteriaHandler{service, loggerInstance}
}

func (h *HysteriaHandler) AltaBoss(c *gin.Context) {
	var version_api domains.Version

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	version_api.NombreApi = "Hysteria"

	c.JSON(200, version_api)
}
