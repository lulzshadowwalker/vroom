package migrations

import (
	"fmt"
	"log"
)

func (m *Migration) users() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			id CHAR(32) PRIMARY KEY, -- uuid v4 
			username NVARCHAR(32) NOT NULL,
			email NVARCHAR(254) NOT NULL,
			password NVARCHAR(80) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP(),
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP()
		);
	`)
	if err != nil {
		return fmt.Errorf("cannot migrate users table %w", err)
	}

	log.Println("migrated users table successfully ✨")

	return nil
}
