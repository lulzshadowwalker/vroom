package repo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lulzshadowwalker/vroom/internal/database/model"
	"github.com/lulzshadowwalker/vroom/pkg/utils"
)

type RoomsRepo repo

func (r *RoomsRepo) Create(name string, ownerId int) (roomId string, err error) {
	id := uuid.New().String()

	_, err = r.Db.Exec(`
		INSERT INTO rooms(id, name, created_by)
		VALUES(?, ?, ?);
	`, id, name, ownerId)
	if err != nil {
		// TODO add logging
		return "", utils.NewAppErr(fmt.Sprintf("cannot create room %q", err)) // TODO
	}

	return id, nil
}
 
func (r *RoomsRepo) GetUserRooms(userId int) ([]model.Room, error) {
  rows, err := r.Db.Query(`
      SELECT id, name FROM rooms r 
      JOIN user_room ur 
      ON ur.user_id = r.id
      WHERE ur.user_id = ?; 
    `, userId);
  if err != nil {
    return nil, fmt.Errorf("cannot retrieve user rooms %w", err);
  }

  rooms := []model.Room{}
  for rows.Next() {
    r := model.Room{} 
    err = rows.Scan(r.Id, r.Name)
    if err != nil {
      return nil, fmt.Errorf("cannot read room row %w", err)
    }

    rooms = append(rooms, r) 
  }

  return rooms, nil
}
