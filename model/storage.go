package model

import (
	"log"

	"github.com/alexyslozada/migrations/connection"
)

const (
	// Postgres string del nombre del motor de base de datos
	Postgres = "postgres"
	// Mysql string del nombre del motor de base de datos
	Mysql = "mysql"
	// Mssql string del nombre del motor de base de datos
	Mssql = "mssql"
)

type Storage interface {
	Setup() error
	Create(name string) error
	FindByName(name string) (*Migration, error)
	Execute(content string) error
}

// newStorage carga el DAO para el motor necesario
// y lo devuelve
func newStorage(engine string, db *connection.MyDB) Storage {
	var s Storage
	switch engine {
	case Postgres:
		s = NewPsql(db)
	case Mysql:
		fallthrough
	case Mssql:
		fallthrough
	default:
		log.Fatalf("el motor %s aún no está configurado", engine)
	}

	return s
}
