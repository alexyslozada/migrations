package main

import (
	"log"

	"flag"

	"fmt"

	"path/filepath"

	"github.com/alexyslozada/migrations/configuration"
	"github.com/alexyslozada/migrations/connection"
	"github.com/alexyslozada/migrations/model"
)

func main() {
	configFile := flag.String("config", "", "Ubicación del archivo de configuración. Debe incluir el nombre del archivo: Ej: /tu/path/configuration.json")
	sqlFiles := flag.String("migration", "", "Ubicación de los archivos de migración")
	flag.Parse()

	if *configFile == "" || *sqlFiles == "" {
		flag.Usage()
		log.Fatalf("Los flag -config y -migration son obligatorios")
	}
	configuration.LoadConfiguration(*configFile)
	model.LoadStorage(configuration.Get())
	m := model.Migration{}

	db := connection.Connection()
	err := m.Setup(db)
	if err != nil {
		log.Fatalf("no se pudo inicializar las migraciones en la base de datos: %v", err)
	}

	fs := ReadFiles(*sqlFiles)
	process(*sqlFiles, fs, db)

	fmt.Println("Proceso realizado correctamente.")
}

func process(src string, fs []string, db *connection.MyDB) {
	for _, v := range fs {
		m := model.Migration{}
		m.FileName = v
		t, err := m.FindByName(db)
		if err != nil {
			log.Fatalf("no se pudo consultar la migración en la base de datos: %v", err)
		}
		if t.ID > 0 {
			continue
		}
		fmt.Printf("Procesando el archivo: %s\n", m.FileName)
		m.Content = string(ReadContent(filepath.Join(src, m.FileName)))
		err = m.Execute(db)
		if err != nil {
			log.Fatalf("no se pudo ejecutar la migración: %v", err)
		}

		err = m.Create(db)
		if err != nil {
			log.Fatalf("no se pudo insertar la migración en la bd: %v", err)
		}
	}
}
