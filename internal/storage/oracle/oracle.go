package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"plantilla_api/cmd/config"
	"plantilla_api/cmd/utils"
	"sync"
	"time"

	_ "github.com/sijms/go-ora/v2"
	"github.com/sirupsen/logrus"
)

type OracleDB struct {
	db         *sql.DB
	config     *config.DB
	loggerExec *logrus.Logger
	mutex      sync.Mutex
	running    bool
}

var instance *OracleDB
var once sync.Once

func GetInstance(c *config.DB, loggerExec *logrus.Logger) (*OracleDB, error) {
	var err error

	once.Do(func() {
		instance = &OracleDB{
			config:     c,
			running:    false,
			loggerExec: loggerExec,
		}
		err = instance.connect()
	})

	if err != nil {
		instance.logMessage("error", err.Error())
		return nil, err
	}
	return instance, nil
}

func (p *OracleDB) connect() error {

	var err error

	dsn := fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
		p.config.User, p.config.Pass, p.config.Host, p.config.Port, p.config.Name)

	p.db, err = sql.Open("oracle", dsn)
	if err != nil {
		utils.LoggerMessage(p.loggerExec, "error", fmt.Sprintf("Error opening Oracle database: %v", err))
		return err
	}

	err = p.db.Ping()
	if err != nil {
		utils.LoggerMessage(p.loggerExec, "error", fmt.Sprintf("Error opening Oracle database: %v", err))
		return err
	}

	if !p.running {
		go p.reconnectOnFailure()
		p.running = true
	}

	return nil
}

func (p *OracleDB) reconnectOnFailure() {
	for {
		time.Sleep(10 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err := p.db.PingContext(ctx)
		if err != nil {

			utils.LoggerExecMessage(p.loggerExec, "error", err.Error())
			p.mutex.Lock()
			defer cancel()

			p.db.Close()
			err := p.connect()

			if err != nil {
				fmt.Println("Error reconectado a oracle...")
				utils.LoggerExecMessage(p.loggerExec, "error", err.Error())
			} else {
				utils.LoggerExecMessage(p.loggerExec, "info", "Reconnected to the Oracle database successfully.")
			}
			p.mutex.Unlock()

		}

		cancel()
	}
}

func (p *OracleDB) logMessage(level string, message string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.loggerExec != nil {
		switch level {
		case "error":
			p.loggerExec.Error(message)
		case "warn":
			p.loggerExec.Warn(message)
		case "info":
			p.loggerExec.Info(message)
		}
	} else {
		fmt.Printf("[%s] %s\n", level, message)
	}
}

func (p *OracleDB) GetDB() *sql.DB {
	return p.db
}
