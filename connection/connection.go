package connection

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/alexyslozada/migrations/configuration"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	// Postgres string del nombre del motor de base de datos
	Postgres = "postgres"
	// Mysql string del nombre del motor de base de datos
	Mysql = "mysql"
	// Mssql string del nombre del motor de base de datos
	Mssql = "mssql"
)

// MyDB estructura que tiene un objeto tipo sql.MyDB para proteger que no se cierre la conexión
type MyDB struct {
	DB *sql.DB
}

// Connection se conecta a la base de datos y devuelve el pool de conexiones a la base de datos
// Se recibe un objeto tipo puntero de configuration.Configuration
func Connection(config *configuration.Configuration) *MyDB {
	conn, err := sql.Open(config.Engine, connectionString(config))
	if err != nil {
		log.Fatalf("Error al conectarse a la BD: %v", err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatalf("Error al conectarse a la BD: %v", err)
	}

	db := &MyDB{}
	db.DB = conn

	return db
}

// connectionString devuelve la cadena de conexión del motor al que se va a conectar
func connectionString(config *configuration.Configuration) string {
	dns := ""
	switch config.Engine {
	case Postgres:
		dns = fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
			config.DBUser,
			config.DBPassword,
			config.DBServer,
			config.DBPort,
			config.DBName,
			config.DBSslmode,
		)
		if config.DBSslmode == "require" {
			dns = fmt.Sprintf("%s sslrootcert=%s", dns, config.DBSSLRootCert)
		}
	case Mysql:
		fallthrough
	case Mssql:
		fallthrough
	default:
		log.Fatalf("El motor de base de datos %s no está configurado aún.", config.Engine)
	}

	return dns
}
