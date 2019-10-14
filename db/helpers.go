package db

import (
	"fmt"
)

// checkPassword returns true if password already exists in the db
func checkPassword(pass *Password) (bool, error) {
	db, err := decrypt()
	if err != nil {
		return false, err
	}

	query := fmt.Sprintf("select * from passwords where `name`='%s' and `group`='%s'", pass.Name, pass.Group)
	passwords, err := db.doSelect(query)
	if err != nil {
		return false, err
	}

	if len(passwords) > 0 {
		return true, nil
	}

	return false, nil
}
