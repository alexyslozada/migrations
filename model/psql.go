package model

import (
	"context"
	"errors"
	"log"

	"github.com/alexyslozada/migrations/v2/connection"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	DB *pgxpool.Pool
}

// NewPsql devuelve un puntero a Psql
func NewPsql(db *connection.MyDB) *Psql {
	return &Psql{DB: db.DB}
}

// Setup crea la tabla de migraciones en la base de datos
func (p *Psql) Setup() error {
	_, err := p.DB.Exec(context.Background(), setupPsql)
	if err != nil {
		log.Printf("no se pudo crear la tabla de migraciones: %v", err)
		return err
	}

	return nil
}

// Create inserta el nombre del archivo de migración ejecutado
func (p *Psql) Create(name string) error {
	_, err := p.DB.Exec(context.Background(), insertPsql, name)
	if err != nil {
		log.Printf("no se pudo insertar la migración '%s': %v", name, err)
		return err
	}

	return nil
}

// FindByName busca un registro en la tabla de migraciones por nombre
func (p *Psql) FindByName(name string) (*Migration, error) {
	m := &Migration{}
	row := p.DB.QueryRow(context.Background(), findByNamePsql, name)
	err := row.Scan(&m.ID, &m.FileName, &m.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return m, nil
	}
	if err != nil {
		log.Printf("no se pudo consultar la migración '%s': %v", name, err)
		return m, err
	}

	return m, nil
}

// Execute ejecuta el contenido SQL de una migración
func (p *Psql) Execute(content string) error {
	_, err := p.DB.Exec(context.Background(), content)
	if err != nil {
		log.Printf("no se pudo ejecutar la migración: %v", err)
		return err
	}

	return nil
}
