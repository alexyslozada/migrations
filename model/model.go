package model

import (
	"time"

	"github.com/alexyslozada/migrations/connection"
)

// Migration estructura de migraciones
type Migration struct {
	ID        int
	FileName  string
	CreatedAt time.Time
}

// MigrationStore controla el acceso a la tabla de migraciones
type MigrationStore struct {
	Storage Storage
}

// NewStorage devuelve un MigrationStore
func NewStorage(engine string, db *connection.MyDB) *MigrationStore {
	s := newStorage(engine, db)
	return &MigrationStore{s}
}

// Setup crea la tabla de migraciones
func (m *MigrationStore) Setup() error {
	return m.Storage.Setup()
}

// Create registra el nombre del archivo de migración
func (m *MigrationStore) Create(mi *Migration) error {
	return m.Storage.Create(mi.FileName)
}

// FindByName busca un registro en la tabla de migraciones por nombre
func (m *MigrationStore) FindByName(name string) (*Migration, error) {
	return m.Storage.FindByName(name)
}

// Execute ejecuta la migración
func (m *MigrationStore) Execute(content string) error {
	return m.Storage.Execute(content)
}
