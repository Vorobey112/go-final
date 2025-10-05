package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT DEFAULT "",
    repeat VARCHAR(128) DEFAULT ""
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := false
	if os.IsNotExist(err) {
		install = true
	}

	dbConn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = dbConn.Exec(schema)
		if err != nil {
			dbConn.Close()
			return err
		}
	}
	db = dbConn
	return nil
}

func GetDB() *sql.DB {
	return db
}
