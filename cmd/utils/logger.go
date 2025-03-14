package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CustomJSONFormatter struct {
	logrus.JSONFormatter
}

type LogData struct {
	RequestIP string `json:"RequestIP"`
	Endpoint  string `json:"Endpoint"`
}

var loggerInstance *logrus.Logger
var loggerInstanceExec *logrus.Logger
var loggerLock sync.Mutex
var currentHTTPDate string
var currentExecDate string

func (f *CustomJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := map[string]interface{}{
		"level": entry.Level.String(),
		"time":  entry.Time.Format(time.RFC3339),
		"msg":   entry.Message,
	}

	for k, v := range entry.Data {
		data[k] = v
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return append(serialized, '\n'), nil
}

func getNewLoggerInstance(logFileName string) (*logrus.Logger, error) {
	logDir := os.Getenv("LOGS_PATH")

	if logDir == "" {
		return nil, fmt.Errorf("la variable de entorno LOGS_PATH no está definida")
	}

	logFilePath := filepath.Join(logDir, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo de log: %v", err)
	}

	logger := logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&CustomJSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return logger, nil
}

func NewLoggerHTTPInstance() (*logrus.Logger, error) {
	loggerLock.Lock()
	defer loggerLock.Unlock()

	currentDate := time.Now().Format("2006-01-02")
	if loggerInstance == nil || currentDate != currentHTTPDate {
		logFileName := "http_trafic_" + currentDate + ".log"
		newLogger, err := getNewLoggerInstance(logFileName)
		if err != nil {
			return nil, err
		}

		loggerInstance = newLogger
		currentHTTPDate = currentDate
		loggerInstance.Info("Logger HTTP inicializado correctamente")
	}

	return loggerInstance, nil
}

func NewLoggerExecutionInstance() (*logrus.Logger, error) {
	loggerLock.Lock()
	defer loggerLock.Unlock()

	currentDate := time.Now().Format("2006-01-02")
	if loggerInstanceExec == nil || currentDate != currentExecDate {
		logFileName := "execution_" + currentDate + ".log"
		newLogger, err := getNewLoggerInstance(logFileName)
		if err != nil {
			return nil, err
		}

		loggerInstanceExec = newLogger
		currentExecDate = currentDate
		loggerInstanceExec.Info("Logger Execution inicializado correctamente")
	}

	return loggerInstanceExec, nil
}

/*
func LoggerHTTP(c *gin.Context) error {
	logger, err := NewLoggerHTTPInstance()
	if err != nil {
		return err
	}

	data := LogData{
		RequestIP: c.ClientIP(),
		Endpoint:  c.Request.URL.Path,
	}

	logger.WithFields(logrus.Fields{
		"RequestIP": data.RequestIP,
		"Endpoint":  data.Endpoint,
	}).Info("Request received")

	return nil
}
*/

func LoggerHTTP(c *gin.Context) error {
	logger, err := NewLoggerHTTPInstance()
	if err != nil {
		return err
	}

	// Leer el body de la solicitud
	bodyBytes, err := c.GetRawData()
	if err != nil {
		return err
	}

	// Restaurar el body para que otros handlers puedan usarlo
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Parsear el body como JSON
	var bodyJSON map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &bodyJSON); err != nil {
		bodyJSON = map[string]interface{}{"raw": string(bodyBytes)} // Si falla, incluir el body sin procesar
	}

	// Crear los datos del log
	data := LogData{
		RequestIP: c.ClientIP(),
		Endpoint:  c.Request.URL.Path,
	}

	delete(bodyJSON, "password")

	logger.WithFields(logrus.Fields{
		"RequestIP": data.RequestIP,
		"Endpoint":  data.Endpoint,
		"Body":      bodyJSON, // Loguear el body como un objeto JSON
	}).Info("Request received")

	return nil
}

func LoggerMessage(loggerHTTP *logrus.Logger, level, message string) error {
	if loggerHTTP == nil {
		return fmt.Errorf("logger no inicializado")
	}

	switch level {
	case "info":
		loggerHTTP.Info(message)
	case "warn":
		loggerHTTP.Warn(message)
	case "error":
		loggerHTTP.Error(message)
	default:
		loggerHTTP.Info("Default log level: " + message)
	}

	return nil
}

func LoggerExecMessage(loggerExec *logrus.Logger, level, message string) error {
	if loggerExec == nil {
		return fmt.Errorf("logger no inicializado")
	}

	switch level {
	case "info":
		loggerExec.Info(message)
	case "warn":
		loggerExec.Warn(message)
	case "error":
		loggerExec.Error(message)
	default:
		loggerExec.Info("Default log level: " + message)
	}

	return nil
}
