package migrations

import (
	"fmt"
	"log"
)

func (m *Migration) userRoom() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS user_room(
			id INT PRIMARY KEY AUTO_INCREMENT, 
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			room_id INT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP(),
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP()
		);
	`)
	if err != nil {
		return fmt.Errorf("cannot migrate the user_room table %w", err)
	}

	log.Println("migrated user_room table successfully ✨")

	return nil
}
