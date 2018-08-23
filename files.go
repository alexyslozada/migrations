package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

// ReadFiles lee los archivos terminados en .sql y los procesa en migration
func ReadFiles(src string) []string {
	r := make([]string, 0)
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatalf("no se pudo listar los archivos del directorio: %v", err)
	}

	for _, v := range files {
		if !v.IsDir() {
			if filepath.Ext(v.Name()) == ".sql" {
				r = append(r, v.Name())
			}
		}
	}

	return r
}

// ReadContent lee el contenido del archivo
func ReadContent(filename string) []byte {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("no se pudo leer el contenido del archivo %s: %v", filename, err)
	}

	return f
}
