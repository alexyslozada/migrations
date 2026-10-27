package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alexyslozada/migrations/v2/configuration"
	"github.com/alexyslozada/migrations/v2/connection"
	"github.com/alexyslozada/migrations/v2/model"
	"github.com/fatih/color"
)

var (
	colorGray  = color.New(color.FgHiBlack)
	colorGreen = color.New(color.FgGreen)
	colorRed   = color.New(color.FgRed)
)

func ts() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func main() {
	configFile := flag.String("config", "", "Ubicación del archivo de configuración. Debe incluir el nombre del archivo: Ej: /tu/path/configuration.json")
	sqlFiles := flag.String("migration", "", "Ubicación de los archivos de migración")
	flag.Parse()

	if *configFile == "" || *sqlFiles == "" {
		flag.Usage()
		return
	}

	configuration.LoadConfiguration(*configFile)
	cnfg := configuration.Get()

	db := connection.Connection(cnfg)
	ms := model.NewStorage(cnfg.Engine, db)

	err := ms.Setup()
	if err != nil {
		colorRed.Printf("[%s] error al inicializar la tabla de migraciones: %v\n", ts(), err)
		os.Exit(1)
	}

	files := ReadFiles(*sqlFiles)
	processFiles(*sqlFiles, files, ms)

	fmt.Println("Proceso realizado correctamente.")
}

func processFiles(src string, files []string, ms *model.MigrationStore) {
	processed, err := ms.FindNames()
	if err != nil {
		colorRed.Printf("[%s] error al consultar las migraciones registradas: %v\n", ts(), err)
		os.Exit(1)
	}

	for _, v := range files {
		m := model.Migration{}
		m.FileName = v

		if _, ok := processed[m.FileName]; ok {
			colorGray.Printf("[%s] omitido  %s\n", ts(), m.FileName)
			continue
		}

		contents := string(ReadContent(filepath.Join(src, m.FileName)))

		err = ms.Execute(contents)
		if err != nil {
			colorRed.Printf("[%s] fallido  %s: %v\n", ts(), m.FileName, err)
			os.Exit(1)
		}

		err = ms.Create(&m)
		if err != nil {
			colorRed.Printf("[%s] fallido  %s: %v\n", ts(), m.FileName, err)
			os.Exit(1)
		}

		colorGreen.Printf("[%s] migrado  %s\n", ts(), m.FileName)
	}
}
