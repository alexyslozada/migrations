package connection

import (
	"context"
	"fmt"
	"log"

	"github.com/alexyslozada/migrations/configuration"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// Postgres string del nombre del motor de base de datos
	Postgres = "postgres"
	// Mysql string del nombre del motor de base de datos
	Mysql = "mysql"
	// Mssql string del nombre del motor de base de datos
	Mssql = "mssql"
)

// MyDB estructura que tiene un pool de conexiones pgxpool
type MyDB struct {
	DB *pgxpool.Pool
}

// Connection se conecta a la base de datos y devuelve el pool de conexiones
func Connection(config *configuration.Configuration) *MyDB {
	connStr := connectionString(config)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Error al conectarse a la BD: %v", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		log.Fatalf("Error al hacer ping a la BD: %v", err)
	}

	return &MyDB{DB: pool}
}

// connectionString devuelve la cadena de conexión según el motor configurado
func connectionString(config *configuration.Configuration) string {
	switch config.Engine {
	case Postgres:
		dsn := fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
			config.DBUser,
			config.DBPassword,
			config.DBServer,
			config.DBPort,
			config.DBName,
			config.DBSslmode,
		)
		if config.DBSslmode == "require" {
			dsn = fmt.Sprintf("%s sslrootcert=%s", dsn, config.DBSSLRootCert)
		}
		return dsn
	case Mysql:
		fallthrough
	case Mssql:
		fallthrough
	default:
		log.Fatalf("El motor de base de datos %s no está configurado aún.", config.Engine)
	}

	return ""
}
