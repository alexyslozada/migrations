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
	cnfg := configuration.Get()

	db := connection.Connection(cnfg)
	ms := model.NewStorage(cnfg.Engine, db)

	err := ms.Setup()
	if err != nil {
		log.Fatalf("no se pudo inicializar las migraciones en la base de datos: %v", err)
	}

	fs := ReadFiles(*sqlFiles)
	process(*sqlFiles, fs, ms)

	fmt.Println("Proceso realizado correctamente.")
}

func process(src string, fs []string, ms *model.MigrationStore) {
	for _, v := range fs {
		m := model.Migration{}
		m.FileName = v
		t, err := ms.FindByName(m.FileName)
		if err != nil {
			log.Fatalf("no se pudo consultar la migración en la base de datos: %v", err)
		}

		if t.ID > 0 {
			continue
		}

		fmt.Printf("Procesando el archivo: %s\n", m.FileName)
		content := string(ReadContent(filepath.Join(src, m.FileName)))
		err = ms.Execute(content)
		if err != nil {
			log.Fatalf("no se pudo ejecutar la migración: %v", err)
		}

		err = ms.Create(&m)
		if err != nil {
			log.Fatalf("no se pudo insertar la migración en la bd: %v", err)
		}
	}
}
