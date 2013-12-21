package database

import (
	"fmt"
	"tmsu/common"
 	"tmsu/log"
)

func (db *Database) DBConfigGetConfig() (common.DBConfig, error) {
	sql := `select * from dbconfig`

	rows, err := db.connection.Query(sql)
	if err != nil {
		log.Fatalf("Failed to get information from the database: %v", err)
	}
	defer rows.Close()
	var fingerPrintCommand string
	if !rows.Next() {
		log.Fatalf("No dbconfig set yet")
	}
	rows.Scan(&fingerPrintCommand)
	return common.DBConfig{fingerPrintCommand}, nil
}

func (db *Database) DBConfigSetConfig(config common.DBConfig) () {
	var sql string
	sql = `delete from dbconfig`
	_, err := db.connection.Exec(sql)
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	}

	sql = fmt.Sprintf(`insert into dbconfig values ('%v');`, config.FingerPrintCommand)
	_, err = db.connection.Exec(sql)
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	}
	return
}
