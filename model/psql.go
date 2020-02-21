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

type Psql struct {
	DB *connection.MyDB
}

// NewPsql devuelve un puntero a Psql
func NewPsql(db *connection.MyDB) *Psql {
	return &Psql{db}
}

func (p *Psql) getDB() *sql.DB {
	return p.DB.DB
}

// Setup crea la tabla de migraciones en la base de datos
func (p *Psql) Setup() error {
	stmt, err := p.getDB().Prepare(setupPsql)
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

// Create inserta el nombre del archivo de migración ejecutado
func (p *Psql) Create(name string) error {
	stmt, err := p.getDB().Prepare(insertPsql)
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

// FindByName busca los nombres de migración
func (p *Psql) FindByName(name string) (*Migration, error) {
	m := &Migration{}
	stmt, err := p.getDB().Prepare(findByNamePsql)
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

// Execute ejecuta la migración encontrada
func (p *Psql) Execute(content string) error {
	_, err := p.getDB().Exec(content)
	if err != nil {
		log.Printf("no se pudo preparar la sentencia para ejecutar la migración %v", err)
		return err
	}

	return nil
}
