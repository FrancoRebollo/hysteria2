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

	//verifico si existe ya el boss con ese nombre en la BD
	query := `SELECT nombre FROM hysteria.bosses WHERE nombre = $1`

	rows, err := tx.QueryContext(ctx, query, fmt.Sprint(AltaBoss.Nombre))

	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		//si no existe, inserto.
		insert := "INSERT INTO hysteria.bosses (nombre) VALUES ($1)"

		_, err = tx.ExecContext(ctx, insert, fmt.Sprint(AltaBoss.Nombre))

		if err != nil {
			return "", err
		}

		message = "Boss creado exitosamente - "
	}

	rows.Close()

	// query = `SELECT nombre FROM hysteria.bosses WHERE nombre = $1`

	// rows, err = tx.QueryContext(ctx, query, fmt.Sprint(AltaBoss.Nombre))

	// if err != nil {
	// 	return "", err
	// }

	//inserto el resto de los datos
	insert := `UPDATE hysteria.bosses
	SET 
		respawn_time = $2,
		interval_respawn_time = $3,
		unidad_interval_respawn_time = $4,
		lunes = $5,
		martes = $6,
		miercoles = $7,
		jueves = $8,
		viernes = $9,
		sabado = $10,
		domingo = $11
	WHERE nombre = $1`

	// `INSERT INTO hysteria.bosses
	// 	(nombre,respawn_time,interval_respawn_time,unidad_interval_respawn_time,
	// 	lunes,martes,miercoles,jueves,viernes,sabado,domingo)
	//  	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)` //$12 <== saco para manejo de id?

	_, err = tx.ExecContext(ctx, insert, //AltaBoss.IdBosses,
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

func (v HysteriaRepository) AltaAnuncio(ctx context.Context, AltaAnuncio domains.RequestAltaAnuncio) (string, error) {
	var message string

	// Iniciar transacción
	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	// Rollback en caso de error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Comenzar la inserción del nuevo anuncio
	insert := `INSERT INTO hysteria.anuncios (texto, fecha) 
               VALUES ($1, $2) RETURNING id`

	// Ejecutar la inserción
	var id int
	err = tx.QueryRowContext(ctx, insert, AltaAnuncio.Texto, AltaAnuncio.Fecha).Scan(&id)
	if err != nil {
		return "", err
	}

	// Mensaje de éxito
	message = fmt.Sprintf("Anuncio creado exitosamente con ID: %d", id)

	// Confirmar la transacción
	if err = tx.Commit(); err != nil {
		return "", err
	}

	return message, nil
}
