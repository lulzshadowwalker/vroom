package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lulzshadowwalker/vroom/internal/database/model"
	"github.com/lulzshadowwalker/vroom/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo repo

func (r *AuthRepo) SignIn(email, password string) (*model.User, error) {
	rows, err := r.Db.Query(`
		SELECT u.id, u.username, u.email, u.password, DATE(u.created_at), DATE(u.updated_at), r.id, r.name
		FROM users u 
		LEFT JOIN user_room ru 
		ON u.id = ru.user_id
		LEFT JOIN rooms r
		ON r.id = ru.room_id
		WHERE u.email = ?;
	`, email)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve user from db %w", err)
	}
	defer rows.Close()

	var hashedPwd string
	var checkedPassword bool

	u := model.User{
		Rooms: []model.Room{},
	}

	for rows.Next() {
		var room model.Room
		var mbRoomId sql.NullString
		var mbRoomName sql.NullString

		err := rows.Scan(&u.Id, &u.Username, &u.Email, &hashedPwd, &u.CreatedAt, &u.UpdatedAt, &mbRoomId, &mbRoomName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, utils.NewAppErr("user not found")
			}

			return nil, fmt.Errorf("cannot scan user row %w", err)
		}

		if !checkedPassword {
			err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
			if err != nil {
				return nil, utils.NewAppErr("invalid credentials")
			}
		}

		if mbRoomId.Valid {
			room.Id = mbRoomId.String
			room.Name = mbRoomId.String
			u.Rooms = append(u.Rooms, room)
		}

	}
	if u.Id == 0 {
		return nil, utils.NewAppErr("user not found")
	}

	return &u, nil
}

func (r *AuthRepo) SignUp(username, email, password string) (*model.User, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password %w", err)
	}

	tx, err := r.Db.Begin()
	if err != nil {
		return nil, fmt.Errorf("cannot start a db transaction to insert user into db %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO users(username, email, password) 
		VALUES (?, ?, ?);
	`, username, email, pwdHash)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("cannot insert user into database %w", err)
	}

	u := new(model.User)
	err = tx.QueryRow(`
		SELECT u.id, u.username, u.email, DATE(u.created_at), DATE(u.updated_at)
		FROM users u 
		WHERE u.email = ?;
	`, email).Scan(&u.Id, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("cannot retrieve user from db %w", err)
	}

	tx.Commit()

	return u, nil
}
