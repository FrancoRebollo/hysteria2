package config

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
	once            sync.Once
	globalConfig    *GlobalConfiguration
	globalConfigErr error
)

type (
	// Se Obtiene un configuracion Global, que contiene pequeñas configuraciones de la aplicación
	GlobalConfiguration struct {
		App  *App
		DB   []*DB
		HTTP *HTTP
	}

	//Configuracion de la aplicacion (Version, Nombre, Cliente, Fecha de inicio)
	App struct {
		Name         string
		Environment  string
		Client       string
		Version      string
		FechaStartUp string
	}

	//Configuracion para conexiones a la base de datos
	DB struct {
		Connection     string
		ConnectionType string
		User           string
		Pass           string
		Host           string
		Port           string
		Name           string
	}
	//Configuracion para el servidor HTTP
	HTTP struct {
		Url            string
		Port           string
		AllowedOrigins string
		Environment    string
	}
)

/*
Función que permite cargar las variables de entorno y setear la configuración de la API
Es Privada, solo se puede acceder a ella desde el paquete config y se ejecuta una sola vez
por eso se usa sync.Once para respetar el principio de Singleton
*/

func loadGlobalConfiguration() (*GlobalConfiguration, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	var appName string
	if appName = os.Getenv("APP_NAME"); appName == "" {
		appName = "PLANTILLA API REST"
	}

	var appEnv string
	if appEnv = os.Getenv("APP_ENVIRONMENT"); appEnv == "" {
		appEnv = "DEVELOPMENT"
	}

	var appClient string
	if appClient = os.Getenv("APP_CLIENT"); appClient == "" {
		appClient = "NA10"
	}

	var appVer string
	if appVer = os.Getenv("APP_VERSION"); appVer == "" {
		appVer = "1.0.0"
	}

	startupTime := time.Now()
	fecha := startupTime.Format("02/01/2006 15:04:05")

	app := &App{
		Name:         appName,
		Environment:  appEnv,
		Client:       appClient,
		Version:      appVer,
		FechaStartUp: fecha,
	}

	var dbs []*DB

	var dbConnOracle string
	var dbConnTypeOracle string
	var dbUserOracle string
	var dbPassOracle string
	var dbHostOracle string
	var dbPortOracle string
	var dbNameOracle string

	var dbConnPostgres string
	var dbConnTypePostgres string
	var dbUserPostgres string
	var dbPassPostgres string
	var dbHostPostgres string
	var dbPortPostgres string
	var dbNamePostgres string
	if os.Getenv("DB_DATABASES") == "ORACLE_POSTGRES" {
		if dbConnOracle = os.Getenv("DB_CONNECTION_ORACLE"); dbConnOracle == "" {
			dbConnOracle = "ORACLE"
		}

		if dbConnTypeOracle = os.Getenv("DB_CONNECTION_TYPE_ORACLE"); dbConnTypeOracle == "" {
			dbConnTypeOracle = "POOL"
		}

		if dbUserOracle = os.Getenv("DB_USER_ORACLE"); dbUserOracle == "" {
			dbUserOracle = "INTERFACES_TS"
		}

		if dbPassOracle = os.Getenv("DB_PASS_ORACLE"); dbPassOracle == "" {
			dbPassOracle = "TS#3739423415"
		}

		if dbHostOracle = os.Getenv("DB_HOST_ORACLE"); dbHostOracle == "" {
			dbHostOracle = "192.168.150.28"
		}

		if dbPortOracle = os.Getenv("DB_PORT_ORACLE"); dbPortOracle == "" {
			dbPortOracle = "1521"
		}

		if dbNameOracle = os.Getenv("DB_NAME_ORACLE"); dbNameOracle == "" {
			dbNameOracle = "HOSMAIN"
		}

		dbOracle := &DB{
			Connection:     dbConnOracle,
			ConnectionType: dbConnTypeOracle,
			User:           dbUserOracle,
			Pass:           dbPassOracle,
			Host:           dbHostOracle,
			Port:           dbPortOracle,
			Name:           dbNameOracle,
		}

		dbs = append(dbs, dbOracle)

		if dbConnPostgres = os.Getenv("DB_CONNECTION_POSTGRES"); dbConnPostgres == "" {
			dbConnPostgres = "POSTGRES"
		}

		if dbConnTypePostgres = os.Getenv("DB_CONNECTION_TYPE_POSTGRES"); dbConnTypePostgres == "" {
			dbConnTypePostgres = "POOL"
		}

		if dbUserPostgres = os.Getenv("DB_USER_POSTGRES"); dbUserPostgres == "" {
			dbUserPostgres = "postgres"
		}

		if dbPassPostgres = os.Getenv("DB_PASS_POSTGRES"); dbPassPostgres == "" {
			dbPassPostgres = "admpwd2024"
		}

		if dbHostPostgres = os.Getenv("DB_HOST_POSTGRES"); dbHostPostgres == "" {
			dbHostPostgres = "192.168.150.27"
		}

		if dbPortPostgres = os.Getenv("DB_PORT_POSTGRES"); dbPortPostgres == "" {
			dbPortPostgres = "5432"
		}

		if dbNamePostgres = os.Getenv("DB_NAME_POSTGRES"); dbNamePostgres == "" {
			dbNamePostgres = "HOSMAIN"
		}

		dbPostgres := &DB{
			Connection:     dbConnPostgres,
			ConnectionType: dbConnTypePostgres,
			User:           dbUserPostgres,
			Pass:           dbPassPostgres,
			Host:           dbHostPostgres,
			Port:           dbPortPostgres,
			Name:           dbNamePostgres,
		}

		dbs = append(dbs, dbPostgres)

	} else if os.Getenv("DB_DATABASES") == "ORACLE" {
		if dbConnOracle = os.Getenv("DB_CONNECTION_ORACLE"); dbConnOracle == "" {
			dbConnOracle = "ORACLE"
		}

		if dbConnTypeOracle = os.Getenv("DB_CONNECTION_TYPE_ORACLE"); dbConnTypeOracle == "" {
			dbConnTypeOracle = "POOL"
		}
		if dbUserOracle = os.Getenv("DB_USER_ORACLE"); dbUserOracle == "" {
			dbUserOracle = "INTERFACES_TS"
		}

		if dbPassOracle = os.Getenv("DB_PASS_ORACLE"); dbPassOracle == "" {
			dbPassOracle = "TS#3739423415"
		}

		if dbHostOracle = os.Getenv("DB_HOST_ORACLE"); dbHostOracle == "" {
			dbHostOracle = "192.168.150.28"
		}

		if dbPortOracle = os.Getenv("DB_PORT_ORACLE"); dbPortOracle == "" {
			dbPortOracle = "1521"
		}

		if dbNameOracle = os.Getenv("DB_NAME_ORACLE"); dbNameOracle == "" {
			dbNameOracle = "HOSMAIN"
		}

		dbOracle := &DB{
			Connection:     dbConnOracle,
			ConnectionType: dbConnTypeOracle,
			User:           dbUserOracle,
			Pass:           dbPassOracle,
			Host:           dbHostOracle,
			Port:           dbPortOracle,
			Name:           dbNameOracle,
		}

		dbs = append(dbs, dbOracle)
	} else if os.Getenv("DB_DATABASES") == "POSTGRES" {
		if dbConnPostgres = os.Getenv("DB_CONNECTION_POSTGRES"); dbConnPostgres == "" {
			dbConnPostgres = "POSTGRES"
		}

		if dbConnTypePostgres = os.Getenv("DB_CONNECTION_TYPE_POSTGRES"); dbConnTypePostgres == "" {
			dbConnTypePostgres = "POOL"
		}

		if dbUserPostgres = os.Getenv("DB_USER_POSTGRES"); dbUserPostgres == "" {
			dbUserPostgres = "postgres"
		}

		if dbPassPostgres = os.Getenv("DB_PASS_POSTGRES"); dbPassPostgres == "" {
			dbPassPostgres = "admpwd2024"
		}

		if dbHostPostgres = os.Getenv("DB_HOST_POSTGRES"); dbHostPostgres == "" {
			dbHostPostgres = "192.168.150.27"
		}

		if dbPortPostgres = os.Getenv("DB_PORT_POSTGRES"); dbPortPostgres == "" {
			dbPortPostgres = "5432"
		}

		if dbNamePostgres = os.Getenv("DB_NAME_POSTGRES"); dbNamePostgres == "" {
			dbNamePostgres = "HOSMAIN"
		}

		dbPostgres := &DB{
			Connection:     dbConnPostgres,
			ConnectionType: dbConnTypePostgres,
			User:           dbUserPostgres,
			Pass:           dbPassPostgres,
			Host:           dbHostPostgres,
			Port:           dbPortPostgres,
			Name:           dbNamePostgres,
		}

		dbs = append(dbs, dbPostgres)
	} else {
		return nil, errors.New("no se ha definido el tipo de conexión a la base de datos")
	}

	var httpUrl string
	if httpUrl = os.Getenv("HTTP_URL"); httpUrl == "" {
		httpUrl = "localhost"
	}

	var httpPort string
	if httpPort = os.Getenv("HTTP_PORT"); httpPort == "" {
		httpPort = "13000"
	}

	var httpAllowOrigins string
	if httpAllowOrigins = os.Getenv("HTTP_ALLOWED_ORIGINS"); httpAllowOrigins == "" {
		httpAllowOrigins = "*"
	}

	http := &HTTP{
		Url:            httpUrl,
		Port:           httpPort,
		AllowedOrigins: httpAllowOrigins,
		Environment:    appEnv,
	}

	globalConfig := &GlobalConfiguration{
		App:  app,
		DB:   dbs,
		HTTP: http,
	}

	return globalConfig, nil
}

// GetGlobalConfiguration devuelve la instancia única de la configuración global.
func GetGlobalConfiguration() (*GlobalConfiguration, error) {
	once.Do(func() {
		globalConfig, globalConfigErr = loadGlobalConfiguration()
	})
	return globalConfig, globalConfigErr
}
