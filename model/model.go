package model

import (
	"time"

	"github.com/alexyslozada/migrations/connection"
)

type Migration struct {
	ID        int
	FileName  string
	CreatedAt time.Time
	Content   string
}

func (m *Migration) Setup(db *connection.MyDB) error {
	return s.Setup(db)
}

func (m *Migration) Create(db *connection.MyDB) error {
	return s.Create(db, m.FileName)
}

func (m *Migration) FindByName(db *connection.MyDB) (*Migration, error) {
	return s.FindByName(db, m.FileName)
}

func (m *Migration) Execute(db *connection.MyDB) error {
	return s.Execute(db, m.FileName, m.Content)
}
