package main

import (
	"fmt"
	"log"
	"os"
	"plantilla_api/cmd/apis/handlers/hysteria"
	"plantilla_api/cmd/apis/handlers/seguridad"
	"plantilla_api/cmd/apis/router"
	"plantilla_api/cmd/config"
	"plantilla_api/cmd/utils"
	"plantilla_api/internal/storage/oracle"
	"plantilla_api/internal/storage/postgres"
	"plantilla_api/internal/version/repository"
	"plantilla_api/internal/version/services"
)

func main() {

	var dbPostgres *postgres.PostgresDB

	var dbOracle *oracle.OracleDB

	config, err := config.GetGlobalConfiguration()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	loggerExec, err := utils.NewLoggerExecutionInstance()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, v := range config.DB {
		if v.Connection == "ORACLE_POSTGRES" {
			dbOracle, err = oracle.GetInstance(v, loggerExec)
			if err != nil {
				fmt.Println("Error conectando a la base de datos de Oracle:", err)
				utils.LoggerExecMessage(loggerExec, "error", err.Error())
				os.Exit(1)
			}
			dbPostgres, err = postgres.GetInstance(v, loggerExec)
			if err != nil {
				fmt.Println("Error conectando a la base de datos de Postgres:", err)
				utils.LoggerExecMessage(loggerExec, "error", err.Error())
				os.Exit(1)
			}

		} else if v.Connection == "POSTGRES" {
			dbPostgres, err = postgres.GetInstance(v, loggerExec)
			if err != nil {
				fmt.Println("Error conectando a la base de datos de Postgres:", err)
				utils.LoggerExecMessage(loggerExec, "error", err.Error())
				os.Exit(1)

			}

		} else if v.Connection == "ORACLE" {
			fmt.Println("Postgres")
			dbOracle, err = oracle.GetInstance(v, loggerExec)
			if err != nil {
				fmt.Println("Error conectando a la base de datos de Oracle:", err)
				utils.LoggerExecMessage(loggerExec, "error", err.Error())
				os.Exit(1)

			}

		}
	}

	if dbPostgres != nil {
		if err != nil {
			fmt.Println("Error conectando a la base de datos de Postgres:", err)
			utils.LoggerExecMessage(loggerExec, "error", err.Error())
			os.Exit(1)
		} else {
			fmt.Println("Conexión a Postgres exitosa")
		}
	}

	if dbOracle != nil {
		if err != nil {
			fmt.Println("Error conectando a la base de datos de Postgres:", err)
			utils.LoggerExecMessage(loggerExec, "error", err.Error())
			os.Exit(1)
		} else {

			fmt.Println("Conexión a Postgres exitosa")
		}
	}

	loggerHTTPInstance, err := utils.NewLoggerHTTPInstance()
	if err != nil {
		log.Fatalf("error creando instancia de logueo")
	}

	hysteriaRepository := repository.NewSecurityRepository(dbOracle, dbPostgres)
	hysteriaService := services.NewSecurityService(hysteriaRepository, *config.App)
	hysteriaHandler := hysteria.NewHysteriaHandler(hysteriaService, loggerHTTPInstance)

	securityRepository := repository.NewSecurityRepository(dbOracle, dbPostgres)
	securityService := services.NewSecurityService(securityRepository, *config.App)
	securityHandler := seguridad.NewSecurityHandler(securityService, loggerHTTPInstance)

	routes, err := router.NewRouter(config.HTTP, *securityHandler, *hysteriaHandler)
	if err != nil {
		if err != utils.LoggerMessage(loggerHTTPInstance, "error", "error iniciando el router") {
			utils.LoggerExecMessage(loggerExec, "error", err.Error())
			os.Exit(1)
		}
	}
	utils.LoggerExecMessage(loggerExec, "info", "Router inicializado correctamente")

	address := fmt.Sprintf("%s:%s", config.HTTP.Url, config.HTTP.Port)
	err = routes.Listen(address)

	if err == nil {
		utils.LoggerExecMessage(loggerExec, "info", "Servidor incializado correctamente")
	}
	if err != nil {
		utils.LoggerExecMessage(loggerExec, "error", "error inicializando el servidor")
		os.Exit(1)
	}

}
