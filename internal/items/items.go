package items

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/google/uuid"
)

func getSaveFileLocation() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(usr.HomeDir, ".vault", "todo.vault")
}

type Item struct {
	Id      uuid.UUID
	Name    string
	History []string
}

func NewItem(name string) Item {
	return Item{Id: uuid.New(), Name: name}
}

func SaveVaultToFile(vault []Item) {
	file, err := os.Create(getSaveFileLocation())

	if err != nil {
		log.Fatal(err)
	}

	vaultBlob, err := json.Marshal(vault)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.WriteString(string(vaultBlob))

	if err != nil {
		log.Fatal(err)
	}
}

func LoadVaultFromFile() []Item {
	vaultFileExists, err := fileExists(getSaveFileLocation())

	if err != nil {
		log.Fatal("Something went wrong when reading the vault file")
	}

	if !vaultFileExists {
		var vault []Item
		SaveVaultToFile(vault)
	}

	data, err := os.ReadFile(getSaveFileLocation())
	if err != nil {
		log.Fatal(err)
	}

	var vault []Item
	json.Unmarshal(data, &vault)
	return vault
}

func fileExists(path string) (bool, error) {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		log.Fatal("Failed to create directories:", err)
	}

	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
