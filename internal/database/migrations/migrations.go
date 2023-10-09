package migrations

import (
	"database/sql"
	"log"

	"github.com/lulzshadowwalker/vroom/internal/database"
)

type Migration struct {
	db *sql.DB
}

var m *Migration

func init() {
	m = &Migration{
		db: database.Db,
	}
}

func Migrate() error {
	log.Println("🥑 running migrations")
	type migrator func() error

	arr := []migrator{
		m.users,
		m.rooms,
		m.userRoom,
		m.messages,
	}

	for _, m := range arr {
		err := m()
		if err != nil {
			return err
		}
	}

	return nil
}
