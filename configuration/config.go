package configuration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var (
	once   sync.Once
	config = &Configuration{}
)

// Configuration tiene la información de conexión a la base de datos
type Configuration struct {
	// Engine es el motor de base de datos a usar, ej: postgres, mysql, mssql
	Engine string `json:"engine"`

	// DBServer es la dirección o el dominio donde se encuentra el servidor de base de datos
	// ej: 127.0.0.1, localhost, 192.168.0.33
	DBServer string `json:"db_server"`

	// DBPort es el puerto del servidor de bases de datos
	DBPort uint16 `json:"db_port"`

	// DBName nombre de la base de datos donde se ejecutarán las migraciones
	DBName string `json:"db_name"`

	// DBUser nombre del usuario que tiene los privilegios de ejecutar comandos
	// tipo DDL (CREATE, DROP, ALTER, etc)
	DBUser string `json:"db_user"`

	// DBPassword password del usuario de la bd.
	DBPassword string `json:"db_password"`

	// DBSslmode es usado en las conexiones a postgres, si no sabe cual es el valor
	// a usar, utilice `disable`
	DBSslmode string `json:"db_sslmode"`

	// DBSSLRootCert es usado para ubicar el certificado del servidor postgres cuando sslmode es `required`
	DBSSLRootCert string `json:"db_ssl_root_cert"`
}

// Get devuelve la configuración
func Get() *Configuration {
	return config
}

// LoadConfiguration lee el archivo configuration.json
// y lo carga en un objeto de la estructura Configuration
func LoadConfiguration(src string) {
	once.Do(func() {
		b, err := ioutil.ReadFile(src)
		if err != nil {
			log.Fatalf("no se pudo leer el archivo de configuración: %s", err.Error())
		}

		err = json.Unmarshal(b, config)
		if err != nil {
			log.Fatalf("no se pudo parsear el archivo de configuración: %s", err.Error())
		}

		if config.Engine == "" {
			log.Fatal("no se ha cargado la información de configuración. Utilice el flag: -config=/path/configuration.json")
		}
	})
}
