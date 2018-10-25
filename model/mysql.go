package model

import (
	"log"

	"database/sql"

	"github.com/alexyslozada/migrations/connection"
)

const (
	setupMsql = `CREATE TABLE IF NOT EXISTS migrations(
		id INT AUTO_INCREMENT NOT NULL,
		file_name VARCHAR(1024) NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT migrations_id_pk PRIMARY KEY (id),
		CONSTRAINT migrations_file_name_uk UNIQUE (file_name)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	insertMsql     = "INSERT INTO migrations (file_name) VALUES (?)"
	findByNameMsql = "SELECT id, file_name, created_at FROM migrations WHERE file_name = ?"
)

type mysql struct{}

func (*mysql) Setup(db *connection.MyDB) error {
	stmt, err := db.DB.Prepare(setupMsql)
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

func (*mysql) Create(db *connection.MyDB, name string) error {
	stmt, err := db.DB.Prepare(insertMsql)
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

func (*mysql) FindByName(db *connection.MyDB, name string) (*Migration, error) {
	m := &Migration{}
	stmt, err := db.DB.Prepare(findByNameMsql)
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

func (*mysql) Execute(db *connection.MyDB, name, query string) error {
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
