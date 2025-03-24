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

func NewHysteriaHandler(service ports.HysteriaService, loggerInstance *logrus.Logger) *HysteriaHandler {
	if loggerInstance == nil {
		log.Fatal("es necesaria una instancia de logueo para que la app incie")
	}
	return &HysteriaHandler{service, loggerInstance}
}

func (h *HysteriaHandler) AltaBoss(c *gin.Context) {
	var request domains.RequestAltaBoss

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contenedor, err := h.AltaBossAPI(c, request) //<= invoco servicio

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, contenedor)
}

func (h *HysteriaHandler) AltaAnuncio(c *gin.Context) {
	var request domains.RequestAltaAnuncio

	if err := utils.LoggerHTTP(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logger no disponible"})
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contenedor, err := h.AltaAnuncioAPI(c, request) //<= invoco servicio

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, contenedor)
}
