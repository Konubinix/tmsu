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
		return common.DBConfig{""}, nil
	}
	defer rows.Close()
	var fingerPrintCommand string
	if !rows.Next() {
		fingerPrintCommand = ""
	}
	rows.Scan(&fingerPrintCommand)
	return common.DBConfig{fingerPrintCommand}, nil
}

func (db *Database) DBConfigSetConfig(config common.DBConfig) () {
	var sql string
	sql = `delete from dbconfig`
	_, err := db.connection.Exec(sql)
	if err != nil {
		log.Infof(2, "Something went wrong when removing the previous info: %v", err)
	}

	sql = fmt.Sprintf(`insert into dbconfig values ('%v');`, config.FingerPrintCommand)
	_, err = db.connection.Exec(sql)
	if err != nil {
		log.Fatalf("Something went wrong when inserting the dbconfig info: %v", err)
	}
	return
}
