package notalib

import (
	"database/sql"
	"fmt"
	"log"
)

func Migrate(db *sql.DB) {
	versions := []string{
		`CREATE TABLE IF NOT EXISTS reminds (
        id INTEGER PRIMARY KEY ASC,
        tag TEXT NOT NULL,
        title TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_at DATETIME NULL DEFAULT NULL,
        priority INTEGER NOT NULL DEFAULT 0,
        scheduled_at DATETIME NULL,
        period TEXT DEFAULT NULL,
        finished_at DATETIME DEFAULT NULL);`,
		`ALTER TABLE reminds ADD COLUMN secret INTEGER DEFAULT 0;`,
		`update reminds set secret = 1, priority = 0 where priority = -1;`,
	}

	var userVersion int
	err := db.QueryRow("PRAGMA user_version;").Scan(&userVersion)
	if err != nil {
		log.Fatal(err)
		panic("failed to load db version")
	}
	for i := userVersion; i < len(versions); i++ {
		_, err = db.Exec(versions[i])
		if err != nil {
			log.Fatal(err)
			panic("db migration failed")
		}
		_, err = db.Exec(fmt.Sprintf(`PRAGMA user_version = %d`, i+1))
		if err != nil {
			log.Fatal(err)
			panic("db migration failed")
		}
	}
}
