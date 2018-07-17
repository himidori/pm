package db

import (
	"os"

	"github.com/Difrex/gpg"
	"github.com/himidori/pm/utils"
)

// encrypting an unencrypted database with a
// GPG public key
func encrypt(path string) error {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	if utils.PathExists(dbPath) {
		err := utils.Rmfile(dbPath)
		if err != nil {
			return err
		}
	}

	keyPath := os.Getenv("HOME") + "/.PM/.key"

	var err error
	if utils.PathExists(keyPath) {
		key, err := utils.ReadFile(keyPath)
		if err != nil {
			return err
		}

		err = gpg.EncryptFile(key, path, dbPath)
	} else {
		err = gpg.EncryptFileRecipientSelf(path, dbPath)
	}

	if err != nil {
		return err
	}

	// deleting the unencrypted database file
	return utils.Rmfile(path)
}

// decrypting a GPG ecrypted database
// and returning a new DB struct
func decrypt() (*DB, error) {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	path := utils.GetPrefix() + ".pm" + utils.GenerateName()

	err := gpg.DecryptFile(dbPath, path)
	if err != nil {
		return nil, err
	}

	return newDB(path)
}
