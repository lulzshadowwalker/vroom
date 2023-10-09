package migrations

import (
	"fmt"
	"log"
)

func (m *Migration) messages() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS messages(
			id INT PRIMARY KEY AUTO_INCREMENT, 
			sender_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			room_id INT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP(),
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
    		CHECK (NULLIF(content, '') IS NOT NULL)
		);
	`)
	if err != nil {
		return fmt.Errorf("cannot migrate the messages table %w", err)
	}

	log.Println("migrated messages table successfully ✨")

	return nil
}
