package model

import (
	"log"

	"sync"

	"github.com/alexyslozada/migrations/configuration"
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

var (
	s    Storage
	once sync.Once
)

type Storage interface {
	Setup(db *connection.MyDB) error
	Create(db *connection.MyDB, name string) error
	FindByName(db *connection.MyDB, name string) (*Migration, error)
	Execute(db *connection.MyDB, name, query string) error
}

// LoadStorage carga el DAO para el motor necesario
func LoadStorage(config *configuration.Configuration) {
	once.Do(func() {
		switch config.Engine {
		case Postgres:
			s = &Psql{}
		case Mysql:
			s = &mysql{}
		case Mssql:
			fallthrough
		default:
			log.Fatalf("el motor %s aún no está configurado", config.Engine)
		}
	})
}
