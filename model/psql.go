package model

import (
	"log"

	"database/sql"

	"github.com/alexyslozada/migrations/connection"
)

const (
	setupPsql = `CREATE TABLE IF NOT EXISTS migrations(
		id SERIAL NOT NULL,
		file_name VARCHAR(1024) NOT NULL,
		created_at timestamp NOT NULL DEFAULT now(),
		CONSTRAINT migrations_id_pk PRIMARY KEY (id),
		CONSTRAINT migrations_file_name_uk UNIQUE (file_name)
	)`
	insertPsql     = "INSERT INTO migrations (file_name) VALUES ($1)"
	findByNamePsql = "SELECT id, file_name, created_at FROM migrations WHERE file_name = $1"
)

type Psql struct{}

func (*Psql) Setup(db *connection.MyDB) error {
	stmt, err := db.DB.Prepare(setupPsql)
	if err != nil {
		log.Printf("no se pudo preparar la consulta para crear la tabla de migraciones: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("no se pudo crear la tabla de migraciones: %v", err)
		return err
	}

	return nil
}

func (*Psql) Create(db *connection.MyDB, name string) error {
	stmt, err := db.DB.Prepare(insertPsql)
	if err != nil {
		log.Printf("no se pudo preparar la sentencia para insertar la migración: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	if err != nil {
		log.Printf("no se pudo ejecutar la inserción de data en migrations: %v", err)
		return err
	}

	return nil
}

func (*Psql) FindByName(db *connection.MyDB, name string) (*Migration, error) {
	m := &Migration{}
	stmt, err := db.DB.Prepare(findByNamePsql)
	if err != nil {
		log.Printf("no se pudo preparar la sentencia para consultar por nombre la migración: %v", err)
		return m, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(name).Scan(&m.ID, &m.FileName, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return m, nil
	}
	if err != nil {
		log.Printf("no se pudo ejecutar la consulta de data en migrations del archivo %s: %v", name, err)
		return m, err
	}

	return m, nil
}

func (*Psql) Execute(db *connection.MyDB, name, query string) error {
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Printf("no se pudo preparar la sentencia para ejecutar la migración del archivo %s: %v", name, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Printf("no se pudo ejecutar la sentencia de migración del archivo %s: %v", name, err)
		return err
	}

	return nil
}
