# migrations
Permite ejecutar los archivos .sql de migraciones en cualquier base de datos

## Uso
```bash
$ go get -u github.com/alexyslozada/migrations/...
$ cd $GOPATH/src/github.com/alexyslozada/migrations
$ cp configuration.json.example configuration.json
// Edita el archivo configuration.json para que apunte a tu base de datos.

$ go build
$ ./migrations -config=/path/to/your/configuration.json -migration=/path/to/you/directory/contains/sql-files/
```
Los flag -config y -migration son obligatorios. -config indica dónde se encuentra el archivo configuration.json y el flag -migration indica dónde se encuentran los archivos sql de tus migraciones.

## Configuration.json
Este archivo contiene la información necesaria para la conexión a la base de datos.
Copia y pega el archivo `configuration.json.example` en un archivo `configuration.json` y edita su contenido para la correcta conexión a la base de datos.
