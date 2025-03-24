package router

import (
	"plantilla_api/cmd/config"
	"strings"

	apis "plantilla_api/cmd/apis/handlers"
	"plantilla_api/cmd/apis/handlers/hysteria"
	"plantilla_api/cmd/apis/handlers/seguridad"
	"plantilla_api/cmd/utils/constants"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func NewRouter(config *config.HTTP, securityHandler seguridad.SecurityHandler, hysteriaHandler hysteria.HysteriaHandler) (*Router, error) {
	//Debug mode
	if config.Environment == constants.PRODUCCION {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.DefaultConfig()
	allowedOrigins := config.AllowedOrigins
	originsList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	//Server
	router := gin.New()

	//Middlewares
	router.Use(gin.Recovery(), cors.New(ginConfig))

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Rutas
	api := router.Group("/api")
	{
		//Version
		version := api.Group("/version")
		{
			version.GET("/", securityHandler.GetVersion)
		}

		seguridad := api.Group("/usuarios")
		{

			seguridad.POST("/login", securityHandler.Login)
			seguridad.POST("/access_token", securityHandler.GetJWT)
			seguridad.GET("/validacion_token", securityHandler.ValidateJWT)

			seguridad.POST("/recuperar_contraseñas", securityHandler.RecuperacionPassword)
			seguridad.POST("/validacion_canales_digitales", securityHandler.ValidarCanalDigital)
			seguridad.POST("/cambiar_contraseñas", apis.MiddlewareAuthorization(securityHandler), securityHandler.CambioPassword)

			// DOBLE FACTOR AUTENTICACION ENDPOINTS
			seguridad.POST("/2fa", apis.MiddlewareAuthorization(securityHandler), securityHandler.ActivarUser2FA)
			seguridad.POST("/2fa/qr", securityHandler.Generate2FAQR)
			seguridad.POST("/2fa/login", securityHandler.Login2FA)
		}

		config := api.Group("/config")
		{
			config.POST("/canales_digitales", securityHandler.AltaUser)
			config.POST("/accesos/personas", apis.MiddlewareAuthorization(securityHandler), securityHandler.RevokePersona)
			config.POST("/accesos/canales_digitales", apis.MiddlewareAuthorization(securityHandler), securityHandler.RevokeCanalDigital)
			config.POST("/accesos/personas_canales_digitales", apis.MiddlewareAuthorization(securityHandler), securityHandler.RevokeCanalDigitalPersona)
			config.POST("/creacion_canales_digitales", securityHandler.CrearCanalDigital)

		}

		hysteria := api.Group("/hysteria")
		{
			hysteria.POST("/altaBoss", hysteriaHandler.AltaBoss)
			hysteria.POST("/altaAnuncio", hysteriaHandler.AltaAnuncio)
		}

	}

	return &Router{
		router,
	}, nil
}

func (r *Router) Listen(addr string) error {
	return r.Run(addr)
}
