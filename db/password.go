package db

import (
	"fmt"

	pwordgen "github.com/cmiceli/password-generator-go"
)

// Password struct
// holds every database field
type Password struct {
	Id       int
	Name     string
	Resource string
	Password string
	Username string
	Comment  string
	Group    string
}

// method used for inserting a new password
// into the database
func AddPassword(pass *Password) error {
	ok, err := checkPassword(pass)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("Password %s in group %s already present", pass.Name, pass.Group)
	}

	db, err := decrypt()
	if err != nil {
		return err
	}

	query := `
insert into passwords(name, resource, password, username, comment, 'group')
values (?, ?, ?, ?, ?, ?)`
	return db.doQuery(query, pass.Name, pass.Resource, pass.Password,
		pass.Username, pass.Comment, pass.Group)
}

// method used for removing a password
// from the database
func RemovePassword(id int) error {
	db, err := decrypt()
	if err != nil {
		return err
	}

	return db.doQuery("delete from passwords where id=?", id)
}

// method used for fetching
// all passwords stored in the DB
func SelectAll() ([]*Password, error) {
	db, err := decrypt()
	if err != nil {
		return nil, err
	}

	return db.doSelect("select id, name, username, resource, password" +
		", comment, `group` from passwords order by `group`")
}

// method used for selecting passwords
// when the -n flag is provided
func SelectByName(name string) ([]*Password, error) {
	db, err := decrypt()
	if err != nil {
		return nil, err
	}

	if name == "all" {
		return SelectAll()
	}

	query := "select id, name, username, resource, password" +
		", comment, `group` from passwords where name=?"
	return db.doSelect(query, name)
}

// method used for selecting passwords
// when the -g flag is provided
func SelectByGroup(name string) ([]*Password, error) {
	db, err := decrypt()
	if err != nil {
		return nil, err
	}

	query := "select id, name, username, resource, password" +
		", comment, `group` from passwords where `group`=?"
	return db.doSelect(query, name)
}

// method used for selecting passwords
// when both -n and -g flags are provided
func SelectByGroupAndName(name string, group string) ([]*Password, error) {
	db, err := decrypt()
	if err != nil {
		return nil, err
	}

	query := "select id, name, username, resource, password" +
		", comment, `group` from passwords where name=?" +
		" and `group`=?"
	return db.doSelect(query, name, group)
}

// method used for generating a password
// of given length
func GeneratePassword(length int) (string, error) {
	return pwordgen.NewPassword(length), nil
}
