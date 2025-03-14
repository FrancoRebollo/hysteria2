package repository

import (
	"context"
	"fmt"
	"log"
	"plantilla_api/internal/storage/oracle"
	"plantilla_api/internal/storage/postgres"
)

type SecurityRepository struct {
	dbOracl *oracle.OracleDB
	dbPost  *postgres.PostgresDB
}

func NewSecurityRepository(dbOracl *oracle.OracleDB, dbPost *postgres.PostgresDB) *SecurityRepository {
	return &SecurityRepository{
		dbOracl: dbOracl,
		dbPost:  dbPost,
	}
}

func (v SecurityRepository) GetVersion(ctx context.Context) (string, error) {
	var version string
	//Query para obtener el campo "version"
	query := `SELECT version_modelo FROM ts_sec.version_modelo vm
			where vm.fecha_last_update = (select max(fecha_last_update) from ts_sec.version_modelo)`

	// Ejecutamos el query
	rows, err := v.dbPost.GetDB().QueryContext(ctx, query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Iteramos sobre los resultados
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		fmt.Printf("Versión: %s\n", version)
	}

	// Verificamos si ocurrió algún error durante la iteración
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return "", err
	}

	return version, nil
}
