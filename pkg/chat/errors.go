package chat

import "errors"

var (
	ErrRoomNotExist = errors.New("room does not exist")
	ErrNoActiveRoom = errors.New("no active room")
)
