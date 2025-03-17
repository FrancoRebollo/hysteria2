package repository

import (
	"context"
	"fmt"
	"plantilla_api/internal/version/domains"
)

// ----------------------------------------------- //

func (v HysteriaRepository) AltaBoss(ctx context.Context, AltaBoss domains.RequestAltaBoss) (string, error) {
	var message string

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	//verifico si existe ya el boss con ese ID en la BD
	query := `SELECT id_bosses FROM hysteria.bosses WHERE id_bosses = $1`

	rows, err := tx.QueryContext(ctx, query, fmt.Sprint(AltaBoss.IdBosses))

	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		//si no existe, inserto.
		insert := "INSERT INTO hysteria.bosses (id_bosses) VALUES ($1)"

		_, err = tx.ExecContext(ctx, insert, fmt.Sprint(AltaBoss.IdBosses))

		if err != nil {
			return "", err
		}

		message = "Boss creado exitosamente - "
	}

	rows.Close()

	query = `SELECT id_bosses FROM hysteria.bosses WHERE id_bosses = $1`

	rows, err = tx.QueryContext(ctx, query, fmt.Sprint(AltaBoss.IdBosses))

	if err != nil {
		return "", err
	}
	//inserto el resto de los datos
	insert := `INSERT INTO hysteria.bosses
		(nombre,respawn_time,interval_respawn_time,unidad_interval_respawn_time,
		lunes,martes,miercoles,jueves,viernes,sabado,domingo)
	 	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	_, err = tx.ExecContext(ctx, insert, AltaBoss.IdBosses,
		AltaBoss.Nombre,
		AltaBoss.RespawnTime,
		AltaBoss.IntervalRespawnTime,
		AltaBoss.UnidadIntervalRespawnTime,
		AltaBoss.Lunes,
		AltaBoss.Martes,
		AltaBoss.Miercoles,
		AltaBoss.Jueves,
		AltaBoss.Viernes,
		AltaBoss.Sabado,
		AltaBoss.Domingo)

	if err != nil {
		return "", err
	}

	//commit
	if err = tx.Commit(); err != nil {
		return "", err
	}

	return message, nil
}
