package database

// The number of tags in the database.
func (db *Database) DBConfigFingerprintExternalProgram() (string, error) {
	sql := `select fingerprint_command from dbconfig`

	rows, err := db.connection.Query(sql)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var fingerPrintCommand string
	rows.Scan(&fingerPrintCommand)
	return fingerPrintCommand, nil
}
