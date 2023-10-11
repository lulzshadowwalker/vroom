package repo

import "database/sql"

type repo struct {
	Db *sql.DB
}
