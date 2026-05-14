package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/alexyslozada/migrations/v2/configuration"
	"github.com/alexyslozada/migrations/v2/connection"
	"github.com/alexyslozada/migrations/v2/model"
)

func main() {
	configFile := flag.String("config", "", "Ubicación del archivo de configuración. Debe incluir el nombre del archivo: Ej: /tu/path/configuration.json")
	sqlFiles := flag.String("migration", "", "Ubicación de los archivos de migración")
	verbose := flag.Bool("v", false, "Muestra los logs del proceso de migración")
	flag.BoolVar(verbose, "verbose", false, "Muestra los logs del proceso de migración")
	flag.Parse()

	if *configFile == "" || *sqlFiles == "" {
		flag.Usage()
		return
	}

	out := io.Discard
	if *verbose {
		out = os.Stdout
	}
	logger := log.New(out, "", log.LstdFlags)

	configuration.LoadConfiguration(*configFile)
	cnfg := configuration.Get()

	logger.Printf("conectando a la base de datos (%s)...", cnfg.Engine)
	db := connection.Connection(cnfg)

	ms := model.NewStorage(cnfg.Engine, db)

	logger.Println("inicializando tabla de migraciones...")
	err := ms.Setup()
	if err != nil {
		log.Fatalf("no se pudo inicializar las migraciones en la base de datos: %v", err)
	}

	files := ReadFiles(*sqlFiles)
	logger.Printf("%d archivo(s) encontrado(s) en %s", len(files), *sqlFiles)

	processFiles(*sqlFiles, files, ms, logger)

	fmt.Println("Proceso realizado correctamente.")
}

func processFiles(src string, files []string, ms *model.MigrationStore, logger *log.Logger) {
	for _, v := range files {
		m := model.Migration{}
		m.FileName = v
		t, err := ms.FindByName(m.FileName)
		if err != nil {
			log.Fatalf("no se pudo consultar la migración en la base de datos: %v", err)
		}

		if isProcessed(t.ID) {
			logger.Printf("omitiendo %s (ya aplicada)", m.FileName)
			continue
		}

		logger.Printf("aplicando %s...", m.FileName)
		contents := string(ReadContent(filepath.Join(src, m.FileName)))

		err = ms.Execute(contents)
		if err != nil {
			log.Fatalf("no se pudo ejecutar la migración: %v", err)
		}

		err = ms.Create(&m)
		if err != nil {
			log.Fatalf("no se pudo insertar la migración en la bd: %v", err)
		}

		logger.Printf("migración %s aplicada correctamente", m.FileName)
	}
}

func isProcessed(ID int) bool {
	return ID > 0
}
