package migrations

import (
	"fmt"
	"log"
)

func (m *Migration) rooms() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS rooms(
			id INT PRIMARY KEY AUTO_INCREMENT,
			name NVARCHAR(50) NOT NULL,
			password NVARCHAR(80), -- public/private rooms
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP(),
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP()
		);
	`)
	if err != nil {
		return fmt.Errorf("cannot migrate the rooms table %w", err)
	}

	log.Println("migrated rooms table successfully ✨")

	return nil
}
