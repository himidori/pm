package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/himidori/pm/db"
)

func vaultClient() (*api.Client, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return client, nil
}

// ExportPasswords exports all the passwords to the Vault secure storage
func ExportPasswords() error {
	client, err := vaultClient()
	if err != nil {
		return err
	}

	passwords, err := db.SelectAll()
	if err != nil {
		return err
	}

	for _, password := range passwords {
		fmt.Println("Export a password with id", password.Id)
		data := map[string]interface{}{
			password.Name: password,
		}
		path := "cubbyhole/%s/%s"
		vaultPath := fmt.Sprintf(path, password.Group, password.Name)
		if password.Group == "" {
			vaultPath = fmt.Sprintf(path, "null", password.Name)
		}
		_, err := client.Logical().Write(vaultPath, data)
		if err != nil {
			fmt.Println("cannot write a password with id", password.Id, "to the Vault: ", err.Error())
		}
	}
	return nil
}
