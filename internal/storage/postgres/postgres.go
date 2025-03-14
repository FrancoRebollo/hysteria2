package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"plantilla_api/cmd/config"
	"plantilla_api/cmd/utils"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type PostgresDB struct {
	db         *sql.DB
	config     *config.DB
	loggerExec *logrus.Logger
	mutex      sync.Mutex
	running    bool
}

var instance *PostgresDB
var once sync.Once

func GetInstance(c *config.DB, loggerExec *logrus.Logger) (*PostgresDB, error) {
	var err error

	once.Do(func() {
		instance = &PostgresDB{
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

func (p *PostgresDB) connect() error {

	var err error

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		p.config.User, p.config.Pass, p.config.Host, p.config.Port, p.config.Name)

	p.db, err = sql.Open("postgres", dsn)
	if err != nil {
		p.logMessage("error", fmt.Sprintf("Error opening PostgreSQL database: %v", err))
		return err
	}

	err = p.db.Ping()
	if err != nil {
		p.logMessage("error", fmt.Sprintf("Error connecting to PostgreSQL database: %v", err))
		return err
	}

	if !p.running {
		go p.reconnectOnFailure()
		p.running = true
	}

	return nil
}

func (p *PostgresDB) reconnectOnFailure() {
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
				fmt.Println("Error reconectado a Postgres...")
				utils.LoggerExecMessage(p.loggerExec, "error", err.Error())
			} else {
				utils.LoggerExecMessage(p.loggerExec, "info", "Reconnected to the Oracle database successfully.")
			}
			p.mutex.Unlock()

		}

		cancel()
	}

}

func (p *PostgresDB) logMessage(level string, message string) {
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

func (p *PostgresDB) GetDB() *sql.DB {
	return p.db
}
